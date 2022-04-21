package main

import (
	"log"

	"github.com/mispon/hey_grpc/internal/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
