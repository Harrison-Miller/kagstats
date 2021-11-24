package main

import (
	"encoding/json"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type RegisterClanReq struct {
	Name string `json:"name"`
}

type InviteMemberReq struct {
	Username string `json:"username"`
}

const clansAs = ` c.ID "clan.ID", c.name "clan.name", c.createdAt "clan.createdAt", c.leaderID "clan.leaderID" `
const leadersAs = ` l.ID "leader.ID", l.username "leader.username",
l.charactername "leader.charactername", l.clantag "leader.clantag", l.clanID "leader.clanID", c.name "leader.clan_info.name", l.joinedClan "leader.joinedClan", l.oldgold "leader.oldgold",
l.registered "leader.registered", l.role "leader.role", l.avatar "leader.avatar", l.tier "leader.tier",
l.gold "leader.gold", l.silver "leader.silver", l.bronze "leader.bronze", l.participation "leader.participation",
l.github "leader.github", l.community "leader.community", l.mapmaker "leader.mapmaker", l.moderation "leader.moderation",
l.leaderboardBan "leader.leaderboardBan", l.statsBan "leader.statsBan" `

func clanError(w http.ResponseWriter, err error) {
	log.Printf("Error getting clan: %v\n", err)
	http.Error(w, "Error getting clan", http.StatusInternalServerError)
}

func notLeaderError(w http.ResponseWriter) {
	log.Printf("Error not the clan leader\n")
	http.Error(w, "Error not the clan leader", http.StatusUnauthorized)
}

func ClanRoutes(public *mux.Router, protected *mux.Router) {
	protected.HandleFunc("/clans/register", RegisterClan).Methods("POST")
	public.HandleFunc("/clans/{id:[0-9]+}", GetClan).Methods("GET")
	protected.HandleFunc("/clans/{id:[0-9]+}/disband", DisbandClan).Methods("POST")
	protected.HandleFunc("/clans/{id:[0-9]+}/invite", InviteMember).Methods("POST")
	protected.HandleFunc("/clans/{id:[0-9]+}/invites", GetClanInvites).Methods("GET")
	protected.HandleFunc("/clans/{id:[0-9]+}/invites/cancel/{playerID:[0-9]+}", CancelInvite).Methods("POST")
	protected.HandleFunc("/clans/invites", GetMyInvites).Methods("GET")
	protected.HandleFunc("/clans/{id:[0-9]+}/decline", DeclineInvite).Methods("POST")
	protected.HandleFunc("/clans/{id:[0-9]+}/accept", AcceptInvite).Methods("POST")
	public.HandleFunc("/clans/{id:[0-9]+}/members", GetMembers).Methods("GET")
	protected.HandleFunc("/clans/{id:[0-9]+}/kick/{playerID:[0-9]+}", KickMember).Methods("POSt")
	protected.HandleFunc("/clans/leave", LeaveClan).Methods("POST")
	public.HandleFunc("/clans", GetClans).Methods("GET")
}

func RegisterClan(w http.ResponseWriter, r *http.Request) {
	var req RegisterClanReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("could not parse request: %s\n", err)
		http.Error(w, "could not parse request", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO clan_info (name, lowerName, createdAt, leaderID) VALUES (?, ?, ?, ?)", req.Name, strings.ToLower(req.Name), utils.NowAsUnixMilliseconds(), claims.PlayerID)
	if err != nil {
		log.Printf("failed to create clan: %s\n", err)
		http.Error(w, "failed to create clan", http.StatusBadRequest)
		return
	}

	_, err = tx.Exec("UPDATE players SET clanID=(SELECT ID FROM clan_info WHERE name=?), joinedClan=? WHERE ID=?", req.Name, utils.NowAsUnixMilliseconds(), claims.PlayerID)
	if err != nil {
		log.Printf("failed to set clanID:  %s\n", err)
		http.Error(w, "failed to set clan id", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("internal error\n")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// who cares about errors here we already added it to the db successfully
	var clan models.ClanInfo
	db.Get(&clan, "SELECT * FROM clan_info WHERE name=?", req.Name)
	JSONResponse(w, &clan)
}

func GetClan(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	var clan models.ClanInfo
	err = db.Get(&clan, "SELECT c.*, " + leadersAs + " FROM clan_info as c INNER JOIN players as l ON leaderID=l.ID WHERE c.ID=? AND c.banned=false", clanID)
	if err != nil {
		clanError(w, err)
		return
	}

	JSONResponse(w, &clan)
}

func DisbandClan(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	var clan models.ClanInfo
	err = db.Get(&clan, "SELECT * FROM clan_info WHERE ID=?", clanID)
	if err != nil {
		clanError(w, err)
		return
	}

	if clan.LeaderID != claims.PlayerID {
		notLeaderError(w)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM clan_info WHERE ID=?", clanID)
	if err != nil {
		log.Printf("failed to delete clan: %s\n", err)
		http.Error(w, "failed to delete clan", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE players SET clanID=NULL WHERE clanID=?", clanID)
	if err != nil {
		log.Printf("failed to update players clans: %s\n", err)
		http.Error(w, "failed to update players clans", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("internal error\n")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func InviteMember(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	var req InviteMemberReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("could not parse request: %s\n", err)
		http.Error(w, "could not parse request", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	var clan models.ClanInfo
	err = db.Get(&clan, "SELECT * FROM clan_info WHERE ID=?", clanID)
	if err != nil {
		clanError(w, err)
		return
	}

	if clan.LeaderID != claims.PlayerID {
		notLeaderError(w)
		return
	}

	var player models.Player
	err = db.Get(&player, "SELECT * FROM players WHERE username=?", req.Username)
	if err != nil {
		log.Printf("could not find player\n")
		http.Error(w,"could not find player", http.StatusBadRequest)
		return
	}

	if player.ClanID != nil {
		log.Printf("player is already in a clan\n")
		http.Error(w, "player is already in a clan", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO clan_invites (clanID, playerID, sentAt) VALUES (?, ?, ?)", clanID, player.ID, utils.NowAsUnixMilliseconds())
	if err != nil {
		log.Printf("could not make invite: %s\n", err)
		http.Error(w, "could not make invite", http.StatusInternalServerError)
		return
	}
}

func GetClanInvites(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	var clan models.ClanInfo
	err = db.Get(&clan, "SELECT * FROM clan_info WHERE ID=?", clanID)
	if err != nil {
		clanError(w, err)
		return
	}

	if clan.LeaderID != claims.PlayerID {
		notLeaderError(w)
		return
	}

	invites := []models.ClanInvite{}
	err = db.Select(&invites, `SELECT clan_invites.*, ` + playersAs + ` FROM clan_invites INNER JOIN players as p ON playerID=p.ID WHERE clan_invites.clanID=?`, clanID)
	if err != nil {
		log.Printf("Error getting clan invites: %s\n", err)
		http.Error(w, "could not get clan invites", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &invites)
}

func CancelInvite(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	playerID, err := GetIntURLArg("playerID", r)
	if err != nil {
		http.Error(w, "Could not get playerID", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	var clan models.ClanInfo
	err = db.Get(&clan, "SELECT * FROM clan_info WHERE ID=?", clanID)
	if err != nil {
		clanError(w, err)
		return
	}

	if clan.LeaderID != claims.PlayerID {
		notLeaderError(w)
		return
	}

	_, err = db.Exec("DELETE FROM clan_invites WHERE clanID=? AND playerID=?", clanID, playerID)
	if err != nil {
		log.Printf("could not cancel invite:%s\n", err)
		http.Error(w, "could not cancel invite", http.StatusInternalServerError)
		return
	}
}

func GetMyInvites(w http.ResponseWriter, r *http.Request) {
	claims, _ := GetClaims(r)

	invites := []models.ClanInvite{}
	err := db.Select(&invites, `SELECT clan_invites.*, ` + clansAs + "," + leadersAs + ` FROM clan_invites INNER JOIN clan_info as c ON clanID=c.ID INNER JOIN players as l ON leaderID=l.ID WHERE clan_invites.playerID=?`, claims.PlayerID)
	if err != nil {
		log.Printf("Error getting clan invites: %s\n", err)
		http.Error(w, "could not get clan invites", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &invites)
}

func DeclineInvite(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	_, err = db.Exec("DELETE FROM clan_invites WHERE clanID=? AND playerID=?", clanID, claims.PlayerID)
	if err != nil {
		log.Printf("error declining clan invite: %s\n", err)
		http.Error(w, "error declining clan invite", http.StatusInternalServerError)
		return
	}
}

func AcceptInvite(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	// check that the invite still exists
	var invite models.ClanInvite
	err = db.Get(&invite, "SELECT * FROM clan_invites WHERE clanID=? AND playerID=?", clanID, claims.PlayerID)
	if err != nil {
		log.Printf("can't accept non existant invite")
		http.Error(w, "can't accept non exists invite", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = db.Exec("UPDATE players SET clanID=?, joinedClan=? WHERE ID=?", clanID, utils.NowAsUnixMilliseconds(), claims.PlayerID)
	if err != nil {
		log.Printf("error setting clan: %s\n", err)
		http.Error(w, "error setting clan", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM clan_invites WHERE playerID=?", claims.PlayerID)
	if err != nil {
		log.Printf("error deleting invites: %s\n", err)
		http.Error(w, "error deleting invites", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("internal error\n")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func GetMembers(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	members := []BasicStats{}
	err = db.Select(&members, basicQuery + "WHERE clanID=?", clanID)

	JSONResponse(w, &members)
}

func KickMember(w http.ResponseWriter, r *http.Request) {
	clanID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	playerID, err := GetIntURLArg("playerID", r)
	if err != nil {
		http.Error(w, "Could not get playerID", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	var clan models.ClanInfo
	err = db.Get(&clan, "SELECT * FROM clan_info WHERE ID=?", clanID)
	if err != nil {
		clanError(w, err)
		return
	}

	if clan.LeaderID != claims.PlayerID {
		notLeaderError(w)
		return
	}

	if int64(playerID) == claims.PlayerID {
		log.Printf("can't kick yourself\n")
		http.Error(w, "can't kick self", http.StatusBadRequest)
	}

	_, err = db.Exec("UPDATE players SET clanID=NULL WHERE ID=?", playerID)
	if err != nil {
		log.Printf("error kicking player: %s\n", err)
		http.Error(w, "error kicking player", http.StatusInternalServerError)
		return
	}
}

func LeaveClan(w http.ResponseWriter, r *http.Request) {
	claims, _ := GetClaims(r)

	_, err := db.Exec("UPDATE players SET clanID=NULL WHERE ID=?", claims.PlayerID)
	if err != nil {
		log.Printf("error leaving clan: %s\n", err)
		http.Error(w, "error leaving clan", http.StatusInternalServerError)
		return
	}
}

func GetClans(w http.ResponseWriter, r *http.Request) {
	clans := []models.ClanInfo{}
	err := db.Select(&clans, "SELECT (SELECT COUNT(id) FROM players WHERE clanID=c.ID) as membersCount, c.*, " + leadersAs + " FROM clan_info as c JOIN players AS l ON c.leaderID=l.ID")
	if err != nil {
		clanError(w, err)
		return
	}

	JSONResponse(w, &clans)
}