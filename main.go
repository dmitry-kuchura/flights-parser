package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TrafficHub struct {
	name string
	code string
}

type Plane struct {
	model    string
	firstFly string
	age      string
	places   string
}

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

func goGet(data []byte) {
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
						plane.model = planeInfo.Text()
					}

					if 1 == i {
						plane.firstFly = planeInfo.Text()
					}

					if 2 == i {
						plane.age = planeInfo.Text()
					}

					if 3 == i {
						plane.places = planeInfo.Text()
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
				number:              flightNumber,
				info:                plane,
				arrivalTrafficHub:   TrafficHub{code: arrivalTrafficHub, name: arrivalTrafficHubName[:strings.IndexByte(arrivalTrafficHubName, '(')]},
				departureTrafficHub: TrafficHub{code: departureTrafficHub, name: departureTrafficHubName[:strings.IndexByte(departureTrafficHubName, '(')]},
				arrivalTime:         arrivalTime,
				departureTime:       departureTime,
				boardStatus:         boardStatus,
				isCharter:           len(isCharter) > 0,
			}

			flights = append(flights, flight)
		})
	})

	fmt.Println("####### flights = ", len(flights), flights)
}

func main() {
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

	goGet(html)
}
