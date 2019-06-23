package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

func startMonitoringServer(config *Config, pdb *PlayerDatabase, startTime time.Time) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("tmpl/monitoring.html")

		buf := bytes.NewBuffer(nil)
		f, _ := os.Open(logFilePath)
		io.Copy(buf, f)
		f.Close()

		uptime := time.Since(startTime).Round(time.Second)
		data := struct {
			Config    *Config
			PDB       *PlayerDatabase
			Uptime    time.Duration
			RecentLog string
		}{
			config,
			pdb,
			uptime,
			string(buf.Bytes()),
		}

		t.Execute(w, data)
	})

	address := fmt.Sprintf("0.0.0.0:%d", config.MonitoringPort)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Println(err)
	}
}
