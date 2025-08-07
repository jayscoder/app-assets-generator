package main

import (
	"app-assets-generator/cmd"
	"log"
)

func main() {
	// Execute root command
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}