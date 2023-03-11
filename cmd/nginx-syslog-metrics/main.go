package main

import (
	"github.com/alxark/nginx-syslog-metrics/internal"
	"log"
)

func main() {
	logger := log.Default()

	a := internal.NewApplication(logger)
	if err := a.ConfigureFlags(); err != nil {
		logger.Fatalf("failed to configure flags: %s", err.Error())
	}

	if err := a.Run(); err != nil {
		logger.Fatalf("run failed: %s", err.Error())
	}

	logger.Print("run completed")
}
