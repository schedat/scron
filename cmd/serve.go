package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/schedat/scron/server"

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
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get current directory: %v", err)
			return
		}

		server, err := server.NewScheduler(
			server.SchedulerConfig{ConfigPath: wd + "/config"},
		)

		if err != nil {
			log.Fatalf("Cannot create Scheduler %v", err)
			return
		}

		if err := http.ListenAndServe(":5000", server); err != nil {
			log.Fatalf("could not listen on port 5000 %v", err)
		}
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
