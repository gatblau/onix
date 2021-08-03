package main

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/cmd"
	"log"
)

func main() {
	rootCmd := cmd.InitialiseRootCmd()

	if err := rootCmd.Cmd.Execute(); err != nil {
		log.Printf("Error occured: %s", err)
	}
}
