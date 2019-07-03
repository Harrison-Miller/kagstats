package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
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
