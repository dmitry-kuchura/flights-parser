package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"skyup/api"
	"skyup/utils"

	"github.com/gorilla/mux"
)

func syncFlights() {
	pollInterval := 60000

	timerCh := time.Tick(time.Duration(pollInterval) * time.Millisecond)

	for range timerCh {
		utils.Parse()
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/flights", api.GetFlights).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
