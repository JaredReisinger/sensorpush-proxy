package cmd

import (
	"fmt"

	"github.com/jaredreisinger/asp"
	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
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

	client, err := sensorpush.NewClient(username, password)
	if err != nil {
		return err
	}

	gateways, err := client.Gateways()
	if err != nil {
		return err
	}

	fmt.Print("Gateways\n")
	for _, gateway := range *gateways {
		fmt.Printf("  %s\n", gateway.Name)
		fmt.Printf("    ID     : %s\n", gateway.ID)
		fmt.Printf("    Version: %s\n", gateway.Version)
		fmt.Printf("    Paired : %t\n", gateway.Paired)
	}

	sensors, err := client.Sensors()
	if err != nil {
		return err
	}

	fmt.Print("Sensors\n")
	for _, sensor := range *sensors {
		fmt.Printf("  %s\n", sensor.Name)
		fmt.Printf("    ID  : %s\n", sensor.ID)
		fmt.Printf("    Type: %s\n", sensor.Type)
	}

	return nil
}
