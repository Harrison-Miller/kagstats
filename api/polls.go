package main

import (
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func PollRoutes(public *mux.Router, protected *mux.Router) {
	public.HandleFunc("/poll", GetPoll).Methods("GET")
	protected.HandleFunc("/poll", AnswerPoll).Methods("POST")
}

func GetPoll(w http.ResponseWriter, r *http.Request) {
	// get current poll
	var poll models.Poll

	now := time.Now()
	err := db.Get(&poll, "SELECT ID,name,description FROM polls WHERE startAt<? AND endAt>?", now.Unix(), now.Unix())
	if err != nil {
		http.Error(w, "No current poll", http.StatusNotFound)
		return
	}

	// get questions for poll
	err = db.Select(&poll.Questions, "SELECT questionID,question,options,required FROM poll_questions WHERE pollID=?", poll.ID)
	if err != nil {
		http.Error(w, "Could not get questions for poll", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &poll)
}

func AnswerPoll(w http.ResponseWriter, r *http.Request) {

}