package config

import (
	"log"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/asp"
	"github.com/spf13/cobra"
)

type CommonFields struct {
	FirstName string
	LastName  string
}

type SensorPushAuth struct {
	Username string `asp.long:"username" asp.short:"u"`
	Password string `asp.long:"password" asp.short:"p"`
}

type Config struct {
	// CommonFields //`mapstructure:",squash"` // <==== this is the needed tag!
	// More         string

	SensorPush struct {
		// Username string `asp.long:"username" asp.short:"u"`
		// Password string `asp.long:"password" asp.short:"p"`
		SensorPushAuth `mapstructure:",squash"`
		DeviceIDs      map[string]string `asp.long:"device-ids" asp.desc:"provides symbolic names for the device IDs to fetch and return.\nUse\n  symbolicName=numeric-device-id,otherName=another-id\nstyle formatting for the values"`
	}

	// Dummy  []int
	// Buffer []byte

	// When    time.Time
	// Extent  time.Duration
	// Extents []time.Duration

	// Host     string `asp.desc:"The host to use for something..."`
	// HostName string `asp:"hostname,"`
	// Nested   struct {
	// 	Inner string `asp:""`
	// } `asp.env:"NNN"`
	// // private   string `asp.env:""`
	// ALongName string `asp:"long-name,l,LONG_NAME"`
	// Numbers   []int
}

var Default = Config{
	// CommonFields: CommonFields{
	// 	FirstName: "Mia",
	// },
}

func Init() (config *Config, err error) {
	config = &Config{
		// Host: "default",
		// private: "foo",
	}
	a, err := asp.Attach(&cobra.Command{}, *config)
	// asp.Command().Usage()
	a.Debug()
	x := a.Config()
	log.Printf("loaded config: %#v", x)
	// asp.Execute(func(cfg *Config, args []string) {
	// 	log.Printf("inside inner func!!! %v", cfg)
	// })
	return
}
