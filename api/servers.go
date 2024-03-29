package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

// GetServers godoc
// @Tags Servers
// @Summary returns a list of all servers
// @Produce json
// @Success 200 {object} []models.Server
// @Router /servers [get]
func GetServers(w http.ResponseWriter, r *http.Request) {
	var servers []models.Server
	err := db.Select(&servers, "SELECT * FROM servers")
	if err != nil {
		log.Printf("Could not get servers: %v\n", err)
		http.Error(w, "Could not get servers", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &servers)
}

// GetServer godoc
// @Tags Servers
// @Summary returns information about a server given an id
// @Produce json
// @Param id path int true "Server ID"
// @Success 200 {object} models.Server
// @Router /server/{id} [get]
func GetServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var server models.Server
	err = db.Get(&server, "SELECT * FROM servers WHERE ID=?", int64(serverID))
	if err != nil {
		log.Printf("Could not get server: %v\n", err)
		http.Error(w, "Could not get server", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &server)
}
