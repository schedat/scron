package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type (
	// SchedulerConfig contains config for SchedulerServer
	SchedulerConfig struct {
		ConfigPath string
	}

	// SchedulerServer schedules and executes jobs
	SchedulerServer struct {
		Config SchedulerConfig
	}

	// Job provides public information about a job
	Job struct {
		ID      string
		Name    string
		Enabled bool
	}
)

func (p *SchedulerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = p.Config.ConfigPath + "/config.yml"

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to read config file")
		return
	}

	config, err := Parse(bufio.NewReader(file))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to parse config")
		return
	}

	var payload []Job
	for _, job := range config.Jobs {
		payload = append(payload, Job{
			ID:      job.ID,
			Name:    job.Name,
			Enabled: job.Enabled,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
