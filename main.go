package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"skyup/api"
	"skyup/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func syncFlights() {
	pollInterval := 60000

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

	router := mux.NewRouter()

	router.HandleFunc("/api/flights", api.GetFlights).Methods("GET")
	router.HandleFunc("/api/flights/departed", api.GetDepartedFlights).Methods("GET")
	router.HandleFunc("/api/flights/arrival", api.GetArrivalFlights).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
