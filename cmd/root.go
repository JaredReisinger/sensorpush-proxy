package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

// SensorPushAuth is the authentication information needed for SensorPush.
type SensorPushAuth struct {
	Username string `asp:"username,,,password for SensorPush"`
	Password string `asp:"password,,,username for SensorPush"`
}

// RootConfig is the root-level only config (common for all subcommands)
type RootConfig struct {
	SensorPush struct {
		// in case we want to add more SensorPush config?
		SensorPushAuth `mapstructure:",squash"`
	}
}

var rootCmd = &cobra.Command{
	Use: "sensorpush-proxy",
	// Short: "",
	// Long: "",

	Run: func(cmd *cobra.Command, args []string) {
		a := cmd.Context().Value(asp.ContextKey).(asp.Asp[RootConfig])
		// cfg := a.Config()
		// log.Printf("got config: %#v", cfg)
		a.Command().Help()
	},
}

var rootDefaults = RootConfig{}

// provided by main!
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func SetVars(versionX string, commitX string, dateX string) {
	// log.Printf()
	version = versionX
	commit = commitX
	date = dateX

	rootCmd.Version = fmt.Sprintf("%s (%.7s : %s)", version, commit, date)
}

var aspOptions = []asp.Option{
	asp.WithEnvPrefix("SPP_"),
	asp.WithDefaultConfigName(".sensorpush-proxy"),
}

// Execute is the main entrypoint into the sensorpush-proxy CLI.
func Execute() {
	// The cobra docs show the flags getting set in init() (pre-Execute), and a
	// cobra.OnInitialize() handler to manage any config file (which is called
	// by cobra prior to each/any command's Execute() call).  Each subcommand
	// *could* have distinct config/flags... I need to think about how that
	// ought to be represented.  It would be great if the env prefix and/or
	// config file name persisted to the subcommands by default.
	err := asp.Attach(rootCmd, rootDefaults, aspOptions...)
	cobra.CheckErr(err)

	// if err := rootCmd.ExecuteContext(ctx); err != nil {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
