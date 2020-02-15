package vanguard

import (
	"net/http"
	"time"
)

// Data structures

type valuation struct {
	Date     string  `json:"date"`
	NavPrice float64 `json:"navPrice"`
}

type valuations []valuation

const dateLayout = "2006-01-02"

// End Data structures

type vanguard struct {
	client http.Client
	url    string
}

// Vanguard Client for Vanguard UK funds
type Vanguard interface {
	FetchLastNDays(days uint16) valuations
}

// New creates a new Vanguard client
func New(url string) Vanguard {
	return &vanguard{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		url: url,
	}
}

func (v *vanguard) FetchLastNDays(days uint16) valuations {

	now := time.Now()
	nowDateStr := now.Format(dateLayout)
	twoHundredDaysAgo := now.AddDate(0, 0, -200)
	twoHundredDaysAgoStr := twoHundredDaysAgo.Format(dateLayout)

	resp, err := v.client.Get()

	return nil
}

func computeMovingAverage(vals valuations, ndays ...uint16) map[uint16]float64 {

	return nil
}
