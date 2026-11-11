package core

import (
	"fmt"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
)

func Query(
	sensorPushUsername string,
	sensorPushPassword string,
) error {
	client, err := sensorpush.NewClient(sensorPushUsername, sensorPushPassword)
	if err != nil {
		return err
	}

	return query(client)
}

type gatewaySensorer interface {
	Gateways() (*sensorpush.Gateways, error)
	Sensors() (*sensorpush.Sensors, error)
}

func query(
	client gatewaySensorer,
) error {
	gateways, err := client.Gateways()
	if err != nil {
		return err
	}

	if gateways != nil {
		fmt.Print("Gateways\n")
		for _, gateway := range *gateways {
			fmt.Printf("  %s\n", gateway.Name)
			fmt.Printf("    ID     : %s\n", gateway.ID)
			fmt.Printf("    Version: %s\n", gateway.Version)
			fmt.Printf("    Paired : %t\n", gateway.Paired)
		}
	}

	sensors, err := client.Sensors()
	if err != nil {
		return err
	}

	if sensors != nil {
		fmt.Print("Sensors\n")
		for _, sensor := range *sensors {
			fmt.Printf("  %s\n", sensor.Name)
			fmt.Printf("    ID  : %s\n", sensor.ID)
			fmt.Printf("    Type: %s\n", sensor.Type)
		}
	}

	return nil
}
