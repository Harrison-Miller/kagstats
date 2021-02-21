package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"fmt"
)

const playersAs = ` p.ID "player.ID", p.username "player.username",
p.charactername "player.charactername", p.clantag "player.clantag", p.oldgold "player.oldgold",
p.registered "player.registered", p.role "player.role", p.avatar "player.avatar", p.tier "player.tier",
p.gold "player.gold", p.silver "player.silver", p.bronze "player.bronze", p.participation "player.participation",
p.github "player.github", p.community "player.community", p.mapmaker "player.mapmaker", p.moderation "player.moderation",
p.leaderboardBan "player.leaderboardBan", p.statsBan "player.statsBan" `

func GetIntURLArg(name string, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars[name])
	return value, err
}

func GetURLParam(name string, defaultValue int, r *http.Request) (int, error) {
	if v := r.URL.Query().Get(name); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, err
		}
		return l, nil
	}
	return defaultValue, nil
}

func JSONResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(i)
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func GetClaims(r *http.Request) (*PlayerClaims, error) {
	claims, ok := r.Context().Value("claims").(*PlayerClaims)
	if !ok {
		return nil, fmt.Errorf("no claims")
	}

	return claims, nil
}
