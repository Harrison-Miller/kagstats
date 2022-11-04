package main

import (
	"encoding/json"
	"fmt"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func PollRoutes(public *mux.Router, protected *mux.Router) {
	public.HandleFunc("/poll", GetPoll).Methods("GET")
	protected.HandleFunc("/poll", AnswerPoll).Methods("POST")
	protected.HandleFunc("/poll/completed", PollCompleted).Methods("GET")

	pollViewer := protected.NewRoute().Subrouter()
	pollViewer.Use(WithPermissions([]string{"poll_viewer"}))
	pollViewer.HandleFunc("/poll/{id:[0-9]+}", GetPollByID).Methods("GET")
	pollViewer.HandleFunc("/poll/{id:[0-9]+}/download", DownloadResponses).Methods("GET")
	pollViewer.HandleFunc("/polls", ListPolls).Methods("GET")
}

type PollCompletedResp struct {
	Completed bool `json:"completed"`
}

func PollCompleted(w http.ResponseWriter, r *http.Request) {
	// get current poll
	var poll models.Poll

	now := time.Now()
	err := db.Get(&poll, "SELECT ID,name,description FROM polls WHERE startAt<? AND endAt>? LIMIT 1", now.Unix(), now.Unix())
	if err != nil {
		http.Error(w, "No current poll", http.StatusNotFound)
		return
	}

	// has this player already answered the poll
	claims, _ := GetClaims(r)
	var answers []models.PollAnswer

	err = db.Select(&answers, "SELECT questionID,answer FROM poll_answers WHERE pollID=? AND playerID=?", poll.ID, claims.PlayerID)
	if err != nil {
		log.Printf("could not get player answers: %s\n", err)
		http.Error(w, "could not get player answers", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, PollCompletedResp{
		Completed: len(answers) > 0,
	})
}

func ListPolls(w http.ResponseWriter, r *http.Request) {
	// get polls
	polls := []models.Poll{}
	err := db.Select(&polls, "SELECT ID,name,description,startAt,endAt FROM polls ORDER BY startAt ASC")
	if err != nil {
		http.Error(w, "failed to get polls", http.StatusInternalServerError)
		return
	}

	for i, poll := range polls {
		// get questions for poll
		err = db.Select(&poll.Questions, "SELECT questionID,question,options,required FROM poll_questions WHERE pollID=?", poll.ID)
		if err != nil {
			http.Error(w, "Could not get questions for poll", http.StatusInternalServerError)
			return
		}

		// get # of responses
		if len(poll.Questions) > 0 {
			count := 0
			err := db.Get(&count, "SELECT COUNT(*) FROM poll_answers WHERE questionID=?", poll.Questions[0].QuestionID)
			if err != nil {
				http.Error(w, "Could not get number of responses", http.StatusInternalServerError)
				return
			}
			poll.Responses = count
		}

		polls[i] = poll
	}

	JSONResponse(w, &polls)
}

func getFileName(title string) string {
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " ", "_")
	title = strings.Trim(title, ",.`!?:;\"/\\")
	return title + ".csv"
}

type answerSet struct {
	stats BasicStats `json:"stats" db:"stats"`
	answers []models.PollAnswer
}

func generateCSV(poll models.Poll, answers map[int64]*answerSet) (io.Reader, int) {
	b := strings.Builder{}
	// column names
	b.WriteString("playerID,")
	// questions
	for _, question := range poll.Questions {
		b.WriteString("\"" + question.Question + "\",")
	}
	// player stats
	b.WriteString("ArcherKills,ArcherDeaths,BuilderKills,BuilderDeaths,KnightKills,KnightDeaths,TotalKills,TotalDeaths\n")

	// write answer sets
	for _, set := range answers {
		// playerID
		b.WriteString(fmt.Sprint(set.stats.PlayerID,","))

		// write answers
		for _, answer := range set.answers {
			b.WriteString("\"" + answer.Answer + "\",")
		}
		// write stats
		s := set.stats
		b.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d\n", s.ArcherKills,s.ArcherDeaths,s.BuilderKills,s.BuilderDeaths,s.KnightKills,s.KnightDeaths,s.TotalKills,s.TotalDeaths))
	}

	return strings.NewReader(b.String()), b.Len()
}

func DownloadResponses(w http.ResponseWriter, r *http.Request) {
	pollID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	// get poll
	var poll models.Poll

	err = db.Get(&poll, "SELECT ID,name,description,startAt,endAt FROM polls WHERE ID=?", pollID)
	if err != nil {
		http.Error(w, "failed to get poll", http.StatusNotFound)
		return
	}

	// get questions for poll
	err = db.Select(&poll.Questions, "SELECT questionID,question,options,required FROM poll_questions WHERE pollID=? ORDER BY questionID ASC", poll.ID)
	if err != nil {
		log.Printf("error getting questions: %s\n", err)
		http.Error(w, "Could not get questions for poll", http.StatusInternalServerError)
		return
	}

	answers := []models.PollAnswer{}
	// get # of responses
	if len(poll.Questions) > 0 {
		err := db.Select(&answers, "SELECT questionID,answer,playerID FROM poll_answers WHERE pollID=? ORDER BY playerID ASC, questionID ASC", poll.ID)
		if err != nil {
			log.Printf("error getting answers: %s\n", err)
			http.Error(w, "Could not get answers", http.StatusInternalServerError)
			return
		}
	}

	answersWithStats := map[int64]*answerSet{}
	for _, answer := range answers {
		if set, ok := answersWithStats[answer.PlayerID]; ok {
			set.answers = append(set.answers, answer)
		} else {
			set := &answerSet{}
			set.answers = append(set.answers, answer)
			// get stats
			err = db.Get(&set.stats, basicQuery+"WHERE basic_stats.playerID=?", answer.PlayerID)
			if err != nil {
				playerNotFoundError(w, err)
				return
			}
			answersWithStats[answer.PlayerID] = set
		}
	}

	reader, size := generateCSV(poll, answersWithStats)

	//Send the headers before sending the file
	fileName := getFileName(poll.Name)
	log.Println("fileName: ", fileName)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", fmt.Sprint(size))

	//Send the file
	io.Copy(w, reader)
}

func GetPollByID(w http.ResponseWriter, r *http.Request) {
	pollID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	// get poll
	var poll models.Poll

	err = db.Get(&poll, "SELECT ID,name,description,startAt,endAt FROM polls WHERE ID=?", pollID)
	if err != nil {
		http.Error(w, "failed to get poll", http.StatusNotFound)
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

func GetPoll(w http.ResponseWriter, r *http.Request) {
	// get current poll
	var poll models.Poll

	now := time.Now()
	err := db.Get(&poll, "SELECT ID,name,description,startAt,endAt FROM polls WHERE startAt<? AND endAt>? LIMIT 1", now.Unix(), now.Unix())
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
