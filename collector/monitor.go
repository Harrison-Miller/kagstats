package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

func startMonitoringServer(config *Config, pdb *PlayerDatabase, startTime time.Time) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("tmpl/monitoring.html")

		uptime := time.Since(startTime).Round(time.Second)
		data := struct {
			Config *Config
			PDB    *PlayerDatabase
			Uptime time.Duration
		}{
			config,
			pdb,
			uptime,
		}

		t.Execute(w, data)
	})

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err)
	}
}
