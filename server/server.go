package server

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
)

// SchedulerServer definition
type (
	// SchedulerConfig contains config for SchedulerServer
	SchedulerConfig struct {
		ConfigPath string
	}

	// SchedulerServer schedules and executes jobs
	SchedulerServer struct {
		Config    SchedulerConfig
		jobs      []job
		schedules []schedule
		cron      *cron.Cron
	}

	schedule struct {
		id          cron.EntryID
		job         string
		description string
		trigger     string // Cron expression
	}
)

// Representation
type (
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
	// ScheduledJob specifies Schedule with next execution time
	ScheduledJob struct {
		ID            int
		Job           string
		Description   string
		Trigger       string // Cron expression
		NextExecution time.Time
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

	crn := cron.New()

	server := SchedulerServer{
		Config: config,
		jobs:   cfg.Jobs,
		cron:   crn,
	}
	crn.Start()
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
		var scheds []ScheduledJob

		for _, s := range p.schedules {
			scheds = append(scheds, ScheduledJob{
				ID:            int(s.id),
				Job:           s.job,
				Description:   s.description,
				Trigger:       s.trigger,
				NextExecution: p.cron.Entry(s.id).Next,
			})
		}

		json.NewEncoder(w).Encode(scheds)
	} else if r.Method == http.MethodPost {
		var payload Schedule
		json.NewDecoder(r.Body).Decode(&payload)

		if job := p.findJobByID(payload.Job); job != nil {
			id, err := p.cron.AddFunc(payload.Trigger, func() {
				cmd := exec.Command(job.Program, job.Arguments)
				out, err := cmd.Output()

				if err != nil {
					println(err.Error())
					return
				}

				print(string(out))
			})
			if err == nil {
				sched := schedule{
					id:          id,
					job:         payload.Job,
					description: payload.Description,
					trigger:     payload.Trigger,
				}

				p.schedules = append(p.schedules, sched)

				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

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

func (p *SchedulerServer) findJobByID(job string) *job {
	for _, j := range p.jobs {
		if j.ID == job {
			return &j
		}
	}

	return nil
}
