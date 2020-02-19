package models

type Flight struct {
	number              string
	info                Plane
	arrivalTrafficHub   TrafficHub
	departureTrafficHub TrafficHub
	arrivalTime         string
	departureTime       string
	boardStatus         string
	isCharter           bool
}
