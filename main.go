package main

import (
	"skyup/utils"
	"time"
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
