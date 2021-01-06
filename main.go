package main

import (
	"fmt"
	"log"
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

	client := api.NewClient(ip.String())

	for {
		t, err := client.GetStatistics()
		if err != nil {
			panic(err)
		}

		log.Printf(
			"Current Month: %s / %s",
			humanReadable(t.CurrentMonthUpload+t.CurrentMonthDownload),
			humanReadable(t.MaxLimit),
		)
		log.Printf(
			"Until Yesterday: %s / %s",
			humanReadable(t.UntilYesterdayUpload3Days+t.UntilYesterdayDownload3Days),
			humanReadable(t.MaxLimit3Days),
		)
		log.Printf(
			"Until Today: %s / %s",
			humanReadable(t.UntilTodayUpload3Days+t.UntilTodayDownload3Days),
			humanReadable(t.MaxLimit3Days),
		)

		time.Sleep(3 * time.Second)
	}
}

func humanReadable(size int64) string {
	s := float64(size)
	units := []string{"Bytes", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; s > 1024; i++ {
		s /= 1024
	}
	return fmt.Sprintf("%.2f%s", s, units[i])
}
