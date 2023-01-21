package main

import (
	"github.com/jaredreisinger/sensorpush-proxy/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	cmd.SetVars(version, commit, date)
}

func main() {
	cmd.Execute()
}
