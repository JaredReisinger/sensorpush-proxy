package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

type SensorPushAuth struct {
	Username string `asp.long:"username" asp.short:"u"`
	Password string `asp.long:"password" asp.short:"p"`
}

type ProxyConfig struct {
	Port      string            `asp.long:"port"`
	DeviceIDs map[string]string `asp.long:"device-ids" asp.desc:"sets the Proxy.DeviceIDs value, which provides symbolic names for the device IDs to fetch and return. Use\n  symbolicName=numeric-device-id,otherName=another-id\nstyle formatting for the values"`
}

type Config struct {
	SensorPush struct {
		// Username  string            `asp.long:"username" asp.short:"u"`
		// Password  string            `asp.long:"password" asp.short:"p"`
		SensorPushAuth `mapstructure:",squash"`
		// DeviceIDs      map[string]string `asp.long:"device-ids" asp.desc:"provides symbolic names for the device IDs to fetch and return.\nUse\n  symbolicName=numeric-device-id,otherName=another-id\nstyle formatting for the values"`
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
