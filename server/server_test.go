package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobs(t *testing.T) {
	server := newFakeScheduler(t)

	t.Run("returns all jobs", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := make([]Job, 0)
		json.NewDecoder(response.Body).Decode(&got)

		want := []Job{
			Job{ID: "backup-database", Name: "Backup User Database", Enabled: false},
			Job{ID: "renew-letsencrypt", Name: "Renew LetsEncrypt certificates", Enabled: true},
		}

		assert.Equal(t, got, want)
	})
}

func TestSchedules(t *testing.T) {
	t.Run("Schedule a job", func(t *testing.T) {
		server := newFakeScheduler(t)

		var jsonStr = []byte(`
		{
			"job":"backup-database",
			"description": "Backup every minute",
			"trigger": "*/1 * * * *"
		}
		`)
		request, _ := http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(jsonStr))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusAccepted)
	})

	t.Run("Get all schedules", func(t *testing.T) {
		server := newFakeScheduler(t)
		addScheduleTo(server)

		request, _ := http.NewRequest(http.MethodGet, "/schedules", nil)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := make([]Schedule, 0)
		json.NewDecoder(response.Body).Decode(&got)

		want := []Schedule{Schedule{
			Job:         "backup-database",
			Description: "Backup every minute",
			Trigger:     "*/1 * * * *",
		}}

		assert.Equal(t, got, want)
	})
}

func newFakeScheduler(t *testing.T) *SchedulerServer {
	t.Helper()
	wd, _ := os.Getwd()
	jobsPath := wd + "/fixtures"

	config := SchedulerConfig{ConfigPath: jobsPath}
	scheduler, err := NewScheduler(config)

	if err != nil {
		t.Error("Unable to create Scheduler")
	}

	return scheduler
}

func getJobsFromResponse(resp *http.Response) (jobs []Job) {
	json.NewDecoder(resp.Body).Decode(&jobs)
	return jobs
}

func addScheduleTo(server *SchedulerServer) {
	var jsonStr = []byte(`
		{
			"job":"backup-database",
			"description": "Backup every minute",
			"trigger": "*/1 * * * *"
		}
		`)
	request, _ := http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
}
