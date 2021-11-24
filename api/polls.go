package main

import (
	"encoding/json"
	"fmt"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

func PollRoutes(public *mux.Router, protected *mux.Router) {
	public.HandleFunc("/poll", GetPoll).Methods("GET")
	protected.HandleFunc("/poll", AnswerPoll).Methods("POST")
}

type PollResponse struct {
	Completed bool `json:"completed"`
	Poll models.Poll `json:"poll"`
}

func GetPoll(w http.ResponseWriter, r *http.Request) {
	// get current poll
	var poll models.Poll

	now := time.Now()
	err := db.Get(&poll, "SELECT ID,name,description FROM polls WHERE startAt<? AND endAt>? LIMIT 1", now.Unix(), now.Unix())
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

	// has this player already answered the poll
	claims, _ := GetClaims(r)
	var answers []models.PollAnswer

	if claims != nil {
		err = db.Select(&answers, "SELECT questionID,answer FROM poll_answers WHERE pollID=? AND playerID=?", poll.ID, claims.PlayerID)
		if err != nil {
			log.Printf("could not get player answers: %s\n", err)
			http.Error(w, "could not get player answers", http.StatusInternalServerError)
			return
		}
	}


	JSONResponse(w, PollResponse{
		Completed: len(answers) > 0,
		Poll: poll,
	})
}

func validate(answers []models.PollAnswer, poll models.Poll) error {
	for _, question := range poll.Questions {
		hasAnswer := false
		for _, answer := range answers {
			if answer.QuestionID == question.QuestionID {
				hasAnswer = true
				// required
				if question.Required && answer.Answer == "" {
					return fmt.Errorf("%d is a required question", answer.QuestionID)
				}

				// options
				parts := strings.Split(question.Options, ";")
				if len(parts) > 1 {
					validOption := false
					for _, part := range parts {
						if strings.Contains(part, "Other") {
							validOption = true
							break
						} else if part == answer.Answer {
							validOption = true
							break
						}
					}

					if !validOption {
						return fmt.Errorf("must pick a valid option for: %d", answer.QuestionID)
					}
				}

			}
		}

		if !hasAnswer {
			return fmt.Errorf("missing answer for question: %d", question.QuestionID)
		}
	}

	return nil
}

func AnswerPoll(w http.ResponseWriter, r *http.Request) {
	var req []models.PollAnswer
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("could not parse request: %s\n", err)
		http.Error(w, "could not parse request", http.StatusBadRequest)
		return
	}

	// get current poll
	var poll models.Poll

	now := time.Now()
	err = db.Get(&poll, "SELECT ID,name,description FROM polls WHERE startAt<? AND endAt>? LIMIT 1", now.Unix(), now.Unix())
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

	// validate
	err = validate(req, poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save poll answers
	claims, _ := GetClaims(r)

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	defer tx.Rollback()

	for _, answer := range req {
		_, err := tx.Exec("INSERT INTO poll_answers (pollID,playerID,questionID,answer) VALUES(?,?,?,?)",
			poll.ID,claims.PlayerID,answer.QuestionID,answer.Answer)
		if err != nil {
			log.Printf("failed to submit poll answers: %s\n", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("internal error\n")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}