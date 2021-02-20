package main

import (
	"encoding/json"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type RegisterClanReq struct {
	Name string `json:"name"`
}

func clanError(w http.ResponseWriter, err error) {
	log.Printf("Error getting clan: %v\n", err)
	http.Error(w, "Error getting clan", http.StatusInternalServerError)
}

func ClanRoutes(public *mux.Router, protected *mux.Router) {
	protected.HandleFunc("/clans/register", RegisterClan).Methods("POST")
	public.HandleFunc("/clans/{id:[0-9]+}", GetClan).Methods("GET")
	protected.HandleFunc("/clans/{id:[0-9]+}/disband", DisbandClan).Methods("POST")
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
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO clan_info (name, lowerName, createdAt, leaderID) VALUES (?, ?, ?, ?)", req.Name, strings.ToLower(req.Name), time.Now().Unix(), claims.PlayerID)
	if err != nil {
		log.Printf("failed to creat clan: %s\n", err)
		http.Error(w, "failed to create clan", http.StatusBadRequest)
		return
	}

	_, err = tx.Exec("UPDATE players SET clanID=(SELECT ID FROM clan_info WHERE name=?) WHERE ID=?", req.Name, claims.PlayerID)
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
	err = db.Get(&clan, "SELECT * FROM clan_info WHERE ID=? AND banned=false", clanID)
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
		log.Printf("can not delete clan you don't own\n")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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