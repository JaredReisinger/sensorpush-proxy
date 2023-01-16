package main

import (
	"log"

	"github.com/jaredreisinger/sensorpush-proxy/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// func init() {
// 	cmd.SetVars(version, commit, date)
// }

func main() {
	// TODO: only output version on demand or when starting proxy?
	log.Printf("sensorpush-proxy %s (%s : %s)", version, commit, date)
	cmd.Execute()
}
