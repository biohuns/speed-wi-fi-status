package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/biohuns/speed-wi-fi-status/api"
	"github.com/jackpal/gateway"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func main() {
	ip, err := gateway.DiscoverGateway()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Detected Gateway IP: %s", ip)

	client := api.NewClient(ip.String())

	var status *walk.TextLabel

	mux := sync.Mutex{}
	reloadFunc := func() {
		mux.Lock()
		defer mux.Unlock()
		t, err := client.GetStatistics()
		if err != nil {
			log.Printf("[ERROR] Fetch Error: %s", err)

			// Retry Discover Gateway
			ip, err := gateway.DiscoverGateway()
			if err != nil {
				log.Printf("[ERROR] Discovet Gateway Error: %s", err)
			}
			log.Printf("Detected Gateway IP: %s", ip)
			client.SetHost(ip.String())

			return
		}

		status.SetText(fmt.Sprintf(
			"Current Month: %8s / %8s\nUntil Yesterday: %8s / %8s\nUntil Today: %8s / %8s",
			humanReadable(t.CurrentMonthUpload+t.CurrentMonthDownload),
			humanReadable(t.MaxLimit),
			humanReadable(t.UntilYesterdayUpload3Days+t.UntilYesterdayDownload3Days),
			humanReadable(t.MaxLimit3Days),
			humanReadable(t.UntilTodayUpload3Days+t.UntilTodayDownload3Days),
			humanReadable(t.MaxLimit3Days),
		))
		log.Print("[INFO] Fetch Complete")
	}

	declarative.MainWindow{
		Title: "Speed Wi-Fi Status",
		// MinSize: declarative.Size{
		// 	Height: 300,
		// 	Width:  300,
		// },
		Size: declarative.Size{
			Height: 300,
			Width:  300,
		},
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				AssignTo:  &status,
				Alignment: declarative.AlignHFarVNear,
				// Font: declarative.Font{
				// 	Family: "ＭＳ Ｐゴシック",
				// },
				TextAlignment: declarative.AlignHFarVNear,
			},
			declarative.PushButton{
				Text:      "Reload",
				OnClicked: reloadFunc,
			},
		},
	}.Run()

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

// label0.SetText(fmt.Sprintf(
// 	"Last Updated At: %s",
// 	time.Now().Format("2006/01/02 15:04:05"),
// ))
// label1.SetText(fmt.Sprintf(
// 	"Current Month: %s / %s",
// 	humanReadable(t.CurrentMonthUpload+t.CurrentMonthDownload),
// 	humanReadable(t.MaxLimit),
// ))
// label2.SetText(fmt.Sprintf(
// 	"Until Yesterday: %s / %s",
// 	humanReadable(t.UntilYesterdayUpload3Days+t.UntilYesterdayDownload3Days),
// 	humanReadable(t.MaxLimit3Days),
// ))
// label3.SetText(fmt.Sprintf(
// 	"Until Today: %s / %s",
// 	humanReadable(t.UntilTodayUpload3Days+t.UntilTodayDownload3Days),
// 	humanReadable(t.MaxLimit3Days),
// ))
