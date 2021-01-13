package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type statistics3Days struct {
	XMLName             xml.Name `xml:"response"`
	Text                string   `xml:",chardata"`
	ToYesterdayDownload int64    `xml:"ToYestodayDownload"`
	ToYesterdayUpload   int64    `xml:"ToYestodayUpload"`
	ToTodayDownload     int64    `xml:"ToTodayDownload"`
	ToTodayUpload       int64    `xml:"ToTodayUpload"`
}

type statisticsMonth struct {
	XMLName              xml.Name `xml:"response"`
	Text                 string   `xml:",chardata"`
	CurrentMonthDownload int64    `xml:"CurrentMonthDownload"`
	CurrentMonthUpload   int64    `xml:"CurrentMonthUpload"`
}

type dataLimits struct {
	XMLName              xml.Name `xml:"response"`
	Text                 string   `xml:",chardata"`
	DataLimit            string   `xml:"DataLimit"`
	Trafficmaxlimit      int64    `xml:"trafficmaxlimit"`
	Data3daysLimit       string   `xml:"Data3daysLimit"`
	Traffic3daysmaxlimit int64    `xml:"traffic3daysmaxlimit"`
}

type Client struct {
	host string
}

func NewClient(host string) *Client {
	return &Client{host: host}
}

type TrafficStats struct {
	CurrentMonthUpload          int64 `json:"current_month_upload"`
	CurrentMonthDownload        int64 `json:"current_month_download"`
	UntilYesterdayUpload3Days   int64 `json:"until_yesterday_upload_3_days"`
	UntilYesterdayDownload3Days int64 `json:"until_yesterday_download_3_days"`
	UntilTodayUpload3Days       int64 `json:"until_today_upload_3_days"`
	UntilTodayDownload3Days     int64 `json:"until_today_download_3_days"`
	MaxLimit                    int64 `json:"max_limit"`
	MaxLimit3Days               int64 `json:"max_limit_3_days"`
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) GetStatistics() (*TrafficStats, error) {
	var (
		stats3Days statistics3Days
		statsMonth statisticsMonth
		limit      dataLimits
	)

	if err := c.fetchXML(
		fmt.Sprintf("http://%s/api/monitoring/start_date", c.host),
		&limit,
	); err != nil {
		return nil, fmt.Errorf("failed to get data limit: %w", err)
	}

	if err := c.fetchXML(
		fmt.Sprintf("http://%s/api/monitoring/month_statistics", c.host),
		&statsMonth,
	); err != nil {
		return nil, fmt.Errorf("failed to get month stats: %w", err)
	}

	if err := c.fetchXML(
		fmt.Sprintf("http://%s/api/monitoring/statistics_3days", c.host),
		&stats3Days,
	); err != nil {
		return nil, fmt.Errorf("failed to get 3days stats: %w", err)
	}

	return &TrafficStats{
		CurrentMonthUpload:          statsMonth.CurrentMonthUpload,
		CurrentMonthDownload:        statsMonth.CurrentMonthDownload,
		UntilYesterdayUpload3Days:   stats3Days.ToYesterdayUpload,
		UntilYesterdayDownload3Days: stats3Days.ToYesterdayDownload,
		UntilTodayUpload3Days:       stats3Days.ToTodayUpload,
		UntilTodayDownload3Days:     stats3Days.ToTodayDownload,
		MaxLimit:                    limit.Trafficmaxlimit,
		MaxLimit3Days:               limit.Traffic3daysmaxlimit,
	}, nil
}

func (c *Client) fetchXML(url string, v interface{}) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("api: http get request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("api: failed to close response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New("not success status code returned")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("api: failed to read from response body: %w", err)
	}

	if err := xml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("api: failed to unmarshal xml: %w", err)
	}

	return nil
}
