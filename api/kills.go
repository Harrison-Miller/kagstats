package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)


func getKills(w http.ResponseWriter, r *http.Request) {
	var kills []models.Kill
	var start int64
	var limit int64 = 100

	if v := r.URL.Query().Get("start"); v != "" {
		s, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not parse start: %v", err), http.StatusBadRequest)
		}

		if s < 0 {
			http.Error(w, "start < 0 is not valid", http.StatusBadRequest)
			return
		}
		start = int64(s)
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not parse limit: %v", err), http.StatusBadRequest)
		}
		limit = Min(int64(l), limit)
	}

	err := db.Select(&kills, "SELECT * FROM kills ORDER BY ID DESC LIMIT ?,?", start, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	next := int(start)+len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct {
		Limit int `json:"limit"`
		Start int `json:"start"`
		Size int `json:"size"`
		Next int `json:"next"`
		Kills []models.Kill `json:"kills"`
	}{
		Limit: int(limit),
		Start: int(start),
		Size: len(kills),
		Next: next,
		Kills: kills,
	})
}

func getKill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	killID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var kill models.Kill
	err = db.Get(&kill, "SELECT * FROM kills WHERE ID=?", int64(killID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Kill not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &kill)
}
