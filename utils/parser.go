package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"skyup-parser/models"
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

	//Find each table
	doc.Find(".board__body").Each(func(i int, boardBody *goquery.Selection) {
		items := boardBody.Find("[data-content='yesterday']")
		plane := Plane{}

		items.Find(".board__item").Each(func(i int, boardItem *goquery.Selection) {
			flightNumber := boardItem.Find(".board__flight p").First().Text()

			boardItem.Find(".flight-board__cont").Each(func(i int, flightPlaneInfo *goquery.Selection) {
				flightPlaneInfo.Find(".flight-board__text").Each(func(i int, planeInfo *goquery.Selection) {
					if 0 == i {
						plane.Model = planeInfo.Find(".flight-board__info").Text()
					}

					if 1 == i {
						plane.FirstFly = planeInfo.Find(".flight-board__info").Text()
					}

					if 2 == i {
						plane.Age = planeInfo.Find(".flight-board__info").Text()
					}

					if 3 == i {
						plane.Places = planeInfo.Find(".flight-board__info").Text()
					}
				})
			})

			departureTrafficHubName := ""
			departureTrafficHub := ""
			arrivalTrafficHubName := ""
			arrivalTrafficHub := ""

			boardItem.Find(".board__cities").Each(func(i int, cities *goquery.Selection) {
				cities.Find(".board__city").Each(func(i int, city *goquery.Selection) {
					if 0 == i {
						departure := strings.TrimSpace(city.Find(".board__city--current").Text())
						departureHub := departure[len(departure)-4:]

						departureTrafficHubName = strings.TrimSpace(departure[:strings.IndexByte(departure, '(')])
						departureTrafficHub = strings.TrimSpace(departureHub[:strings.IndexByte(departureHub, ')')])
					}

					if 1 == i {
						arrival := strings.TrimSpace(city.Find(".board__city--current").Text())
						arrivalHub := arrival[len(arrival)-4:]

						arrivalTrafficHubName = strings.TrimSpace(arrival[:strings.IndexByte(arrival, '(')])
						arrivalTrafficHub = strings.TrimSpace(arrivalHub[:strings.IndexByte(arrivalHub, ')')])
					}
				})
			})

			departureTime := ""
			arrivalTime := ""

			boardItem.Find(".board__departure").Each(func(i int, departure *goquery.Selection) {
				departureTime = departure.Find(".board__time").Text()
			})

			boardItem.Find(".board__arrival").Each(func(i int, arrival *goquery.Selection) {
				arrivalTime = arrival.Find(".board__time").Text()
			})

			boardStatus := strings.TrimSpace(boardItem.Find(".board__status").First().Text())
			isCharter := boardItem.Find(".board__status--charter").Text()

			flight := Flight{
				Number:              strings.TrimSpace(flightNumber),
				Info:                plane,
				ArrivalTrafficHub:   TrafficHub{Code: arrivalTrafficHub, Name: arrivalTrafficHubName},
				DepartureTrafficHub: TrafficHub{Code: departureTrafficHub, Name: departureTrafficHubName},
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
