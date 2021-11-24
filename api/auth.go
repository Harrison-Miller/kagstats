package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var AUTH_SECRET = "password"

type LoginReq struct {
	Username string `json:"username"`
	Token string `json:"token"`
}

type PlayerClaims struct {
	PlayerID int64 `json:"playerID"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	ClanID *int64 `json:"clanID"`
	BannedFromMakingClans bool `json:"bannedFromMakingClans"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

func AuthRoutes(r *mux.Router, protected *mux.Router) {
	r.HandleFunc("/login", Login).Methods("POST")
	protected.HandleFunc("/validate", Validate).Methods("GET")
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
	err = db.Get(&player, "SELECT * FROM players WHERE LOWER(username)=LOWER(?)", login.Username)
	if err != nil {
		log.Printf("player not in database: %s\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get the login.Token from whatever was put in
	tokenRegex, err := regexp.Compile(fmt.Sprintf(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}%s[0-9a-f]{128}?`, strings.ToLower(login.Username)))
	if err != nil {
		log.Printf("failed to compile regex: %s\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	matchBytes := tokenRegex.Find([]byte(login.Token))
	if matchBytes == nil {
		log.Printf("no match in regex: %s\n", login.Token)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	match := strings.TrimSpace(string(matchBytes))

	// validate token against api.kag2d.com
	err = utils.ValidateToken(login.Username, match)
	if err != nil {
		log.Printf("player %s failed to login: %s\n", login.Username, err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get permissions
	permissions := []string{}
	err = db.Select(&permissions, "SELECT permission FROM permissions WHERE playerID=?", player.ID)
	if err != nil {
		log.Printf("error retrieving permissions: %s\n", err)
		http.Error(w, "error retrieving permissions", http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(24 * time.Hour * 365)
	claims := PlayerClaims{
		PlayerID: player.ID,
		Username: player.Username,
		Avatar: player.Avatar,
		BannedFromMakingClans: player.BannedFromMakingClans,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	if player.ClanID != nil {
		claims.ClanID = player.ClanID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
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
		Path: "/",
		Secure: prod,
		SameSite: http.SameSiteStrictMode,
	})

	JSONResponse(w, struct{
		Token string `json:"token"'`
	}{
		Token: signed,
	})
}


func WithPermissions(permissionsRequired []string) func (next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := GetClaims(r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			for _, required := range permissionsRequired {
				found := false
				for _, permission := range claims.Permissions {
					if required == permission {
						found = true
						break
					}
				}

				if !found {
					http.Error(w, "unaauthorized", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("KAGSTATS_TOKEN")
		if err != nil {
			log.Printf("no cookie: %s\n", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
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

		ctx := context.WithValue(r.Context(), "claims", &claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func Validate(w http.ResponseWriter, r *http.Request) {
	// refresh the token as well as validate
	oldClaims, err := GetClaims(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get player by id
	var player models.Player
	err = db.Get(&player, "SELECT * FROM players WHERE ID=?", oldClaims.PlayerID)
	if err != nil {
		log.Printf("player not in database: %s\n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get permissions
	permissions := []string{}
	err = db.Select(&permissions, "SELECT permission FROM permissions WHERE playerID=?", player.ID)
	if err != nil {
		log.Printf("error retrieving permissions: %s\n", err)
		http.Error(w, "error retrieving permissions", http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(24 * time.Hour * 365)
	claims := PlayerClaims{
		PlayerID: player.ID,
		Username: player.Username,
		Avatar: player.Avatar,
		BannedFromMakingClans: player.BannedFromMakingClans,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	if player.ClanID != nil {
		claims.ClanID = player.ClanID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
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
		Path: "/",
		Secure: prod,
		SameSite: http.SameSiteStrictMode,
	})

	JSONResponse(w, struct{
		Token string `json:"token"'`
	}{
		Token: signed,
	})

	return
}
