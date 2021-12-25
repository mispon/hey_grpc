package main

import (
	"log"

	"github.com/hey_grpc/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
