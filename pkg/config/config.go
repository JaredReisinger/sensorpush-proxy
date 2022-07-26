package config

import (
	"log"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/asp"
)

type Config struct {
	Host     string `asp.desc:"The host to use for something..."`
	HostName string `asp:"hostname,"`
	Nested   struct {
		Inner string `asp:""`
	} `asp.env:"NNN"`
	// private   string `asp.env:""`
	ALongName string `asp:"long-name,l,LONG_NAME"`
	Numbers   []int
}

func Init() (config *Config, err error) {
	config = &Config{
		Host: "default",
		// private: "foo",
	}
	asp, err := asp.New(config, "APP_")
	// asp.Command().Usage()
	// x := asp.Config()
	// log.Printf("loaded config: %#v", x)
	asp.Execute(func(cfg *Config, args []string) {
		log.Printf("inside inner func!!! %v", cfg)
	})
	return
}
