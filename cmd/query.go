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
	a := cmd.Context().Value(asp.ContextKey).(asp.Asp[RootConfig])
	cfg := a.Config()
	// cfg := getRootConfig(cmd)
	// log.Printf("got config: %#v", cfg)

	// TODO: check flags!

	client, err := sensorpush.NewClient(cfg.SensorPush.Username, cfg.SensorPush.Password)
	if err != nil {
		log.Fatalf("unable to create client: %+v", err)
	}

	// TODO: should lastSample/SuccessfulCall really be channels?  We don't
	// really want to server/read from them *while* updates are happening...
	// var lastSample sensorpush.Sample
	// lastSuccessfulCall := time.Now()

	// appCtx, appCancel := context.WithCancel(context.Background())

	gateways, err := client.Gateways()
	if err != nil {
		log.Printf("unable to get gateways: %+v", err)
	}

	log.Printf("got gateways: %+v", gateways)

	sensors, err := client.Sensors()
	if err != nil {
		log.Printf("unable to get sensors: %+v", err)
	}

	// log.Printf("got sensors: %+v", sensors)

	for _, sensor := range *sensors {
		log.Printf("%s", sensor.Name)
		log.Printf("  ID  : %q", sensor.ID)
		log.Printf("  Type: %s", sensor.Type)
	}

}
