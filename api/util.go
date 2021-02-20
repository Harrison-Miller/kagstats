package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
