package main

import (
	"log"
	"time"

	"skyup-parser/utils"

	"github.com/joho/godotenv"
)

func syncFlights() {
	pollInterval := 10000

	timerCh := time.Tick(time.Duration(pollInterval) * time.Millisecond)

	for range timerCh {
		utils.Parse()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	syncFlights()
}
