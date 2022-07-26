package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	cfg "github.com/jaredreisinger/sensorpush-proxy/pkg/config"
	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
)

const (
	appName = "sensorpush-query"
)

func main() {
	log.Print(appName)

	config, err := cfg.Init()
	if err != nil {
		log.Fatalf("ERROR: (%T) %s", err, err.Error())
	}
	log.Printf("got config: %#v", config)

	if true {
		return
	}

	viper.SetDefault("sensorPush.username", "")
	viper.SetDefault("sensorPush.password", "")
	// viper.SetDefault("deviceId", "")
	// viper.SetDefault("port", ":5375")

	viper.SetConfigName("config")
	// viper.SetConfigType("yaml") // setting the config type takes precedence
	// over the extension, which seems wrong!

	viper.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("SENSORPUSH")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		case *viper.ConfigFileNotFoundError:
			log.Printf("no config file found... perhaps there are environment variables")
		default:
			log.Fatalf("ERROR: (%T) %s", err, err.Error())
		}
	}

	user := viper.GetString("sensorPush.username")
	pass := viper.GetString("sensorPush.password")
	// deviceID := viper.GetString("deviceId")
	// port := viper.GetString("port")

	if user == "" || pass == "" {
		log.Fatalf("one of SENSORPUSH_USERNAME (%q) or SENSORPUSH_PASSWORD (length %d) is missing", user, len(pass))
	}

	client, err := sensorpush.NewClient(user, pass)
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
