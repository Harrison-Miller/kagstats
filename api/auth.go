package main

import (
	"encoding/json"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var AUTH_SECRET = "password"

type LoginReq struct {
	Username string `json:"username"`
	Token string `json:"token"`
}

type PlayerClaims struct {
	PlayerID int `json:"playerID"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	jwt.StandardClaims
}

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/validate", Validate).Methods("GET")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login LoginReq
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		log.Printf("could not parse login request: %s\n", err)
		http.Error(w, "could not parse login request", http.StatusBadRequest)
		return
	}

	// get player by id
	var player models.Player
	err = db.Get(&player, "SELECT * FROM players WHERE username=?", login.Username)
	if err != nil {
		log.Printf("player not in database: %s\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// validate token against api.kag2d.com
	err = utils.ValidateToken(login.Username, login.Token)
	if err != nil {
		log.Printf("player %s failed to login: %s\n", login.Username, err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	expiration := time.Now().Add(24 * time.Hour * 365)
	claims := PlayerClaims{
		PlayerID: int(player.ID),
		Username: player.Username,
		Avatar: player.Avatar,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(AUTH_SECRET))
	if err != nil {
		log.Printf("error signing jwt: %s\n", err)
		http.Error(w, "failed to sign jwt", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       "KAGSTATS_TOKEN",
		Value:      signed,
		Expires:    expiration,
		//Secure:     true,
		//HttpOnly:   true,
	})

	JSONResponse(w, struct{
		Token string `json:"token"'`
	}{
		Token: signed,
	})
}

func Validate(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("KAGSTATS_TOKEN")
	if err != nil {
		log.Printf("no cookie: %s\n", err)
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	var claims PlayerClaims
	_, err = jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(AUTH_SECRET), nil
	})
	if err != nil {
		log.Printf("bad jwt: %s\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
}
