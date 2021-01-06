package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/biohuns/speed-wi-fi-status/api"
	"github.com/jackpal/gateway"
)

func main() {
	ip, err := gateway.DiscoverGateway()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Detected Gateway IP: %s", ip)

	for {
		res, err := http.Get(
			fmt.Sprintf("http://%s/api/monitoring/statistics_3days", ip),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		day3 := new(api.Statistics3Days)
		if err := xml.Unmarshal(buf, day3); err != nil {
			log.Fatal(err)
		}

		log.Printf("Yesterday: %s", convertToHumanReadableSizeString(
			day3.ToYesterdayUpload+day3.ToYesterdayDownload,
		))
		log.Printf("Today: %s", convertToHumanReadableSizeString(
			day3.ToTodayUpload+day3.ToTodayDownload,
		))

		time.Sleep(time.Second)
	}
}

func convertToHumanReadableSizeString(size int64) string {
	s := float64(size)
	units := []string{"Bytes", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; s > 1024; i++ {
		s /= 1024
	}
	return fmt.Sprintf("%.2f %s", s, units[i])
}
