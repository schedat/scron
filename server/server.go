package server

import (
	"bufio"
	"encoding/json"
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
		Config    SchedulerConfig
		jobs      []job
		schedules []Schedule
	}

	// Job provides public information about a job
	Job struct {
		ID      string
		Name    string
		Enabled bool
	}

	// Schedule associates trigger to a job
	Schedule struct {
		Job         string
		Description string
		Trigger     string // Cron expression
	}
)

// NewScheduler creates and initializes new instance of SchedulerServer
func NewScheduler(config SchedulerConfig) (*SchedulerServer, error) {
	var path = config.ConfigPath + "/config.yml"

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cfg, err := Parse(bufio.NewReader(file))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	server := SchedulerServer{Config: config, jobs: cfg.Jobs}
	return &server, nil
}

func (p *SchedulerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	router.HandleFunc("/jobs", http.HandlerFunc(p.handleJobs))
	router.HandleFunc("/schedules", http.HandlerFunc(p.handleSchedules))

	router.ServeHTTP(w, r)
}

func (p *SchedulerServer) handleSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p.schedules)
	} else if r.Method == http.MethodPost {
		var payload Schedule
		json.NewDecoder(r.Body).Decode(&payload)
		p.schedules = append(p.schedules, payload)

		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *SchedulerServer) handleJobs(w http.ResponseWriter, r *http.Request) {
	var payload []Job
	for _, job := range p.jobs {
		payload = append(payload, Job{
			ID:      job.ID,
			Name:    job.Name,
			Enabled: job.Enabled,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
