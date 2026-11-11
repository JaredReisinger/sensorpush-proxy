package cmd

import (
	"time"

	"github.com/jaredreisinger/asp"
	"github.com/jaredreisinger/sensorpush-proxy/pkg/core"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Run a proxy server for specific device IDs",
	RunE:  proxy,
}

// proxyConfig is the configuration for the HTTP proxy that sensorpush-proxy
// provides.
type proxyConfig struct {
	// Using a named sub-struct makes the config file better
	Proxy ProxyConfig `asp.long:"" asp.env:""`
}

type ProxyConfig struct {
	Port         string            `asp:"port,p,,port for the proxy to listen on"`
	Sensors      map[string]string `asp:"sensors,s,,sensors to proxy; use 'symbolicName=id,otherName=another-id'\nstyle formatting for the values"`
	UpdatePeriod time.Duration     `asp:"update-period,u,,duration between updates"`
}

var proxyDefaults = proxyConfig{
	Proxy: ProxyConfig{
		Port:         ":5375",
		UpdatePeriod: 5 * time.Minute,
	},
}

func init() {
	err := asp.Attach(proxyCmd, proxyDefaults, aspOptions...)
	cobra.CheckErr(err)

	rootCmd.AddCommand(proxyCmd)
}

func proxy(cmd *cobra.Command, args []string) error {
	rootCfg, err := asp.Get[rootConfig](cmd.Root())
	if err != nil {
		return err
	}
	cfg, err := asp.Get[proxyConfig](cmd)
	if err != nil {
		return err
	}

	// slog.Info("got root config", "config", rootCfg)
	// slog.Info("got proxy config", "config", cfg)

	server, err := core.NewServer(
		cfg.Proxy.Port,
		rootCfg.SensorPush.Username,
		rootCfg.SensorPush.Password,
		cfg.Proxy.Sensors,
		cfg.Proxy.UpdatePeriod,
	)
	if err != nil {
		return err
	}

	err = server.ListenAndServe()
	return err
}
