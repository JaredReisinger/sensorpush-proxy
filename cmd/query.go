package cmd

import (
	"log"

	"github.com/jaredreisinger/asp"
	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "show all of the available device IDs for the user",
	Run:   query,
}

func query(cmd *cobra.Command, args []string) {
	cfg, err := asp.Get[RootConfig](cmd)
	if err != nil {
		log.Fatalf("unable to get config: %+v", err)
	}

	client, err := sensorpush.NewClient(cfg.SensorPush.Username, cfg.SensorPush.Password)
	if err != nil {
		log.Fatalf("unable to create client: %+v", err)
	}

	gateways, err := client.Gateways()
	if err != nil {
		log.Printf("unable to get gateways: %+v", err)
	}

	log.Print("Gateways")
	for _, gateway := range *gateways {
		log.Printf("  %s", gateway.Name)
		log.Printf("    ID  : %q", gateway.ID)
		log.Printf("    Version: %s", gateway.Version)
		log.Printf("    Paired: %t", gateway.Paired)
	}

	sensors, err := client.Sensors()
	if err != nil {
		log.Printf("unable to get sensors: %+v", err)
	}

	log.Print("Sensors")
	for _, sensor := range *sensors {
		log.Printf("  %s", sensor.Name)
		log.Printf("    ID  : %q", sensor.ID)
		log.Printf("    Type: %s", sensor.Type)
	}

}
