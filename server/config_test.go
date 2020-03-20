package server

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	spec := `
host: app-server-1
jobs:
- id: backup-database
  name: Backup User Database
  program: /usr/bin/backup-database
  arguments: 
  enabled: false
- id: renew-letsencrypt
  name: Renew LetsEncrypt certificates
  program: /bin/certbot
  arguments: renew
  enabled: true
`
	config, err := Parse(strings.NewReader(spec))
	if err != nil {
		fmt.Println(err)
		t.Error("Unable to read config")
	}

	assert.Equal(t, config.Host, "app-server-1")
	assert.Len(t, config.Jobs, 2)

	job1 := config.Jobs[0]
	assert.Equal(t, job1.ID, "backup-database")
	assert.Equal(t, job1.Name, "Backup User Database")
	assert.Equal(t, job1.Program, "/usr/bin/backup-database")
	assert.Equal(t, job1.Arguments, "")
	assert.False(t, job1.Enabled)
}
