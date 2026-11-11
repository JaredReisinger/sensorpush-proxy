package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

// rootConfig is the root-level only config (common for all subcommands)
type rootConfig struct {
	// Using a named sub-struct makes the config file better
	SensorPush SensorPush `asp.long:""`
}

type SensorPush struct {
	Username string `asp:"username,,,username for SensorPush"`
	Password string `asp:"password,,,password for SensorPush,true"`
}

var rootCmd = &cobra.Command{
	Use: "sensorpush-proxy",
	// Short: "",
	Long: `
A credential-hiding proxy for the SensorPush API. Intended to be used to help
provide sensor data for frontend websites.
`,
}

var aspOptions = []asp.Option{
	asp.WithEnvPrefix("SPP_"),
	asp.WithDefaultConfigName(".sensorpush-proxy"),
}

func init() {
	cobra.EnableTraverseRunHooks = true
	err := asp.Attach(rootCmd, rootConfig{}, aspOptions...)
	cobra.CheckErr(err)
}

// provided by main!
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func SetVars(versionX string, commitX string, dateX string) {
	version = versionX
	commit = commitX
	date = dateX

	rootCmd.Version = fmt.Sprintf("%s (%.7s : %s)", version, commit, date)
}

// Execute is the main entrypoint into the sensorpush-proxy CLI.
func Execute() {
	// if err := rootCmd.ExecuteContext(ctx); err != nil {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
