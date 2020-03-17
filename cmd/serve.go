/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os/exec"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

type entry struct {
	schedule string
	job      string
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch scron server",
	Long: `scron server is a process that frequently checks every configured schedules
to find jobs to run. It also monitors and reports statuses of these jobs.`,
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool)

		entries := []entry{entry{schedule: "*/1 * * * *", job: "echo"}}

		c := cron.New()
		for _, ent := range entries {
			c.AddFunc(ent.schedule, func() {
				cmd := exec.Command(ent.job, "Hello")
				out, err := cmd.Output()

				if err != nil {
					println(err.Error())
					return
				}

				print(string(out))
			})
		}
		c.Start()
		<-done
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
