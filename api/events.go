package main

/*
func getPlayerEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var events []models.Event
	err = db.Select(&events, "SELECT * FROM events WHERE playerID=? LIMIT 100", playerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting events: %s", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &events)
}

func getServerEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var events []models.Event
	err = db.Select(&events, "SELECT * FROM events WHERE serverID=? LIMIT 100", serverID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting events: %s", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &events)
}
*/
