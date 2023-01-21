package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

// SensorPushAuth is the authentication information needed for SensorPush.
type SensorPushAuth struct {
	Username string `asp.long:"username" asp.short:"u"`
	Password string `asp.long:"password" asp.short:"p"`
}

// ProxyConfig is the configuration for the HTTP proxy that sensorpush-proxy
// provides.
type ProxyConfig struct {
	Port      string            `asp.long:"port"`
	DeviceIDs map[string]string `asp.long:"device-ids" asp.desc:"sets the Proxy.DeviceIDs value, which provides symbolic names for\nthe device IDs to fetch and return. Use\n  symbolicName=numeric-device-id,otherName=another-id\nstyle formatting for the values"`
}

// Config is the all-up configuration for sensorpush-proxy.
type Config struct {
	SensorPush struct {
		SensorPushAuth `mapstructure:",squash"`
	}

	Proxy ProxyConfig
}

var rootCmd = &cobra.Command{
	Use: "sensorpush-proxy",
	// Short: "",
	// Long: "",

	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		config := getConfig(cmd)
		log.Printf("got config: %#v", config)
	},
}

var configDefaults = Config{
	Proxy: ProxyConfig{
		Port: ":5375",
	},
}

// func SetVars(version string, commit string, date string) {
// 	log.Printf()
// }

// Execute is the main entrypoint into the sensorpush-proxy CLI.
func Execute() {
	a, err := asp.Attach(rootCmd, configDefaults, asp.WithEnvPrefix[Config]("SPP_"))
	cobra.CheckErr(err)

	ctx := context.WithValue(context.Background(), asp.ContextKey, a)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getConfig(cmd *cobra.Command) *Config {
	a := cmd.Context().Value(asp.ContextKey).(asp.Asp[Config])
	config := a.Config()
	return config
}
