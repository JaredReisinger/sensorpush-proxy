package cmd

import (
	"github.com/jaredreisinger/asp"
	"github.com/jaredreisinger/sensorpush-proxy/pkg/core"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Show all of the available device IDs for the user",
	RunE:  query,
}

type queryConfig struct {
}

func init() {
	err := asp.Attach(queryCmd, queryConfig{}, aspOptions...)
	cobra.CheckErr(err)

	rootCmd.AddCommand(queryCmd)
}

func query(cmd *cobra.Command, args []string) error {
	username, password, err := getSensorPushAuth(cmd)
	if err != nil {
		return err
	}

	return core.Query(username, password)
}
