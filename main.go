package main

import (
	"log"

	"github.com/websublime/courier/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
