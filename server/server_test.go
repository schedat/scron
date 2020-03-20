package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchJobs(t *testing.T) {
	wd, _ := os.Getwd()
	jobsPath := wd + "/fixtures"

	config := SchedulerConfig{ConfigPath: jobsPath}
	server := &SchedulerServer{Config: config}

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

func getJobsFromResponse(resp *http.Response) (jobs []Job) {
	json.NewDecoder(resp.Body).Decode(&jobs)
	return jobs
}
