package main

import (
	"time"

	"skyup/utils"
)

func startWork() {
	utils.Parse()
}

func main() {
	pollInterval := 30000

	timerCh := time.Tick(time.Duration(pollInterval) * time.Millisecond)

	for range timerCh {
		startWork()
	}
}
