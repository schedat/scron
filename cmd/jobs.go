package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/schedat/scron/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// jobsCmd represents the jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault("endpoint", "http://localhost:5000")
		endpoint := viper.GetString("endpoint")

		resp, err := http.Get(endpoint + "/jobs")
		if err != nil {
			log.Fatalf("Cannot request the scheduler: %v", err)
			return
		}

		payload := make([]server.Job, 0)
		json.NewDecoder(resp.Body).Decode(&payload)

		data := [][]string{}
		for _, j := range payload {
			data = append(data, []string{j.ID, j.Name, strconv.FormatBool(j.Enabled)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Enabled"})
		table.SetRowLine(true)
		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
