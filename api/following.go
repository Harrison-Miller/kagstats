package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func followedError(w http.ResponseWriter, err error) {
	log.Printf("Error getting followed: %v\n", err)
	http.Error(w, "Error getting followed", http.StatusInternalServerError)
}

type FollowingResp struct {
	FollowingCount int64 `json:"followingCount"`
}

func FollowingRoutes(public *mux.Router, protected *mux.Router) {
	public.HandleFunc("/players/{id:[0-9]+}/followers", GetFollowersCount).Methods("GET")

	protected.HandleFunc("/players/{id:[0-9]+}/follow", FollowPlayer).Methods("POST")
	protected.HandleFunc("/players/{id:[0-9]+}/follow", IsFollowingPlayer).Methods("GET")
	protected.HandleFunc("/players/{id:[0-9]+}/unfollow", UnfollowPlayer).Methods("POST")
	protected.HandleFunc("/following/stats", GetFollowingStats).Methods("GET")
}

func GetFollowersCount(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	var count int64
	err = db.Get(&count, "SELECT COUNT(*) FROM followers WHERE followedID=?", playerID)
	if err != nil {
		followedError(w, err)
		return
	}

	resp := FollowingResp{
		FollowingCount: count,
	}

	JSONResponse(w, &resp)
}

func FollowPlayer(w http.ResponseWriter, r *http.Request) {
	followedID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	_, err = db.Exec("INSERT INTO followers (playerID, followedID) VALUES (?, ?)", claims.PlayerID, followedID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func IsFollowingPlayer(w http.ResponseWriter, r *http.Request) {
	followedID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, err := GetClaims(r)
	if err != nil {
		followedError(w, err)
		return
	}

	var count int64
	err = db.Get(&count, "SELECT COUNT(*) FROM followers WHERE playerID=? AND followedID=?", claims.PlayerID, followedID)
	if err != nil {
		followedError(w, err)
		return
	}
	resp := FollowingResp{
		FollowingCount: count,
	}

	JSONResponse(w, &resp)
}

func UnfollowPlayer(w http.ResponseWriter, r *http.Request) {
	followedID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	claims, _ := GetClaims(r)

	_, err = db.Exec("DELETE FROM followers WHERE playerID=? AND followedID=?", claims.PlayerID, followedID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func GetFollowingStats(w http.ResponseWriter, r *http.Request) {
	stats := []BasicStats{}

	claims, _ := GetClaims(r)

	err := db.Select(&stats, `SELECT basic_stats.*, p.ID "player.ID", p.username "player.username",
		p.charactername "player.charactername", p.clantag "player.clantag", p.clanID "player.clanID", p.joinedClan "player.joinedClan", p.oldgold "player.oldgold",
		p.registered "player.registered", p.role "player.role", p.avatar "player.avatar", p.tier "player.tier",
		p.gold "player.gold", p.silver "player.silver", p.bronze "player.bronze", p.participation "player.participation",
		p.github "player.github", p.community "player.community", p.mapmaker "player.mapmaker", p.moderation "player.moderation",
		p.leaderboardBan "player.leaderboardBan", p.statsBan "player.statsBan", c.name "player.clan_info.name" 
		FROM (SELECT * FROM followers WHERE playerID=?) as followers 
		INNER JOIN basic_stats ON basic_stats.playerID=followers.followedID 
		INNER JOIN players AS p ON basic_stats.playerID=p.ID 
		LEFT JOIN clan_info as c ON p.clanID=c.ID`, claims.PlayerID)
	if err != nil {
		followedError(w, err)
		return
	}

	JSONResponse(w, &stats)
}