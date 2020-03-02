package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"skyup-parser/models"

	"github.com/PuerkitoBio/goquery"
)

type Plane = models.Plane
type TrafficHub = models.TrafficHub
type Flight = models.Flight

func Parse() {
	url := "https://skyup.aero/ru/board"
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	prepareData(html)
}

func prepareData(data []byte) {
	var flights []Flight

	html := string(data)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}

	// Find each table
	doc.Find(".board__body").Each(func(index int, boardBody *goquery.Selection) {
		items := boardBody.Find("[data-date='']")
		plane := Plane{}

		items.Find(".board__item").Each(func(indexItems int, boardItem *goquery.Selection) {
			flightNumber := boardItem.Find(".board__flight p").First().Text()

			boardItem.Find(".flight-popup__cont").Each(func(i int, flightPlaneInfo *goquery.Selection) {
				flightPlaneInfo.Find(".flight-popup__info").Each(func(i int, planeInfo *goquery.Selection) {
					if 0 == i {
						plane.Model = planeInfo.Text()
					}

					if 1 == i {
						plane.FirstFly = planeInfo.Text()
					}

					if 2 == i {
						plane.Age = planeInfo.Text()
					}

					if 3 == i {
						plane.Places = planeInfo.Text()
					}
				})
			})

			departureTrafficHub := boardItem.Find(".board__departure-code-lg").Text()
			departureTrafficHubName := boardItem.Find(".board__departure-city").Text()
			arrivalTrafficHub := boardItem.Find(".board__arrival-code-lg").Text()
			arrivalTrafficHubName := boardItem.Find(".board__arrival-city").Text()
			departureTime := boardItem.Find(".board__departure-time").Text()
			arrivalTime := boardItem.Find(".board__arrival-time").Text()
			boardStatus := boardItem.Find(".board__status").Text()
			isCharter := boardItem.Find(".board__charter-text").Text()

			flight := Flight{
				Number:              flightNumber,
				Info:                plane,
				ArrivalTrafficHub:   TrafficHub{Code: arrivalTrafficHub, Name: strings.TrimSpace(arrivalTrafficHubName[:strings.IndexByte(arrivalTrafficHubName, '(')])},
				DepartureTrafficHub: TrafficHub{Code: departureTrafficHub, Name: strings.TrimSpace(departureTrafficHubName[:strings.IndexByte(departureTrafficHubName, '(')])},
				ArrivalTime:         arrivalTime,
				DepartureTime:       departureTime,
				BoardStatus:         boardStatus,
				IsCharter:           len(isCharter) > 0,
			}

			flights = append(flights, flight)
		})
	})

	collection := GetConnection()

	for x := range flights {
		row := Flight{
			Number:              flights[x].Number,
			Info:                flights[x].Info,
			DepartureTrafficHub: flights[x].DepartureTrafficHub,
			ArrivalTrafficHub:   flights[x].ArrivalTrafficHub,
			DepartureTime:       flights[x].DepartureTime,
			ArrivalTime:         flights[x].ArrivalTime,
			BoardStatus:         flights[x].BoardStatus,
		}

		_, err := FindOne(collection, flights[x].Number)

		if err == nil {
			UpdateOne(collection, row.Number, row)
		} else {
			InsertOne(collection, row)
		}
	}
}
