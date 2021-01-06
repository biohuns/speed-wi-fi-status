package api

import "encoding/xml"

type Statistics3Days struct {
	XMLName                  xml.Name `xml:"response"`
	Text                     string   `xml:",chardata"`
	ToYesterdayDownload      int64    `xml:"ToYestodayDownload"`
	ToYesterdayUpload        int64    `xml:"ToYestodayUpload"`
	ToYesterdayDuration      int64    `xml:"ToYestodayDuration"`
	ToTodayDownload          int64    `xml:"ToTodayDownload"`
	ToTodayUpload            int64    `xml:"ToTodayUpload"`
	ToTodayDuration          int64    `xml:"ToTodayDuration"`
	IsYesterdayFluxOverLimit bool     `xml:"IsYestodayFluxOverLimit"`
	LastClearTime3days       string   `xml:"LastClearTime3days"`
}
