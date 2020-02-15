package vanguard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mxmCherry/movavg"

	log "github.com/sirupsen/logrus"
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

	url := fmt.Sprintf(v.url, twoHundredDaysAgoStr, nowDateStr)
	resp, err := v.client.Get(url)
	if err != nil {
		log.Errorf("Error making GET '%s': %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Errorf("Non OK status code when GET %s: %d %s", resp, resp.StatusCode, resp.Status)
		return nil
	}

	var vals valuations
	err = json.NewDecoder(resp.Body).Decode(&vals)
	if err != nil {
		log.Errorf("Error decoding JSON response from Vanguard: %v", err)
		return nil
	}

	return vals
}

func computeMovingAverage(vals valuations, ndays int) float64 {
	if ndays < 1 {
		log.Errorf("Fuck off, gimme at least one day")
		return -1
	}
	if len(vals) < ndays {
		log.Errorf("You are asking me to compute average for %d days but only getting me %d valuations",
			ndays, len(vals))
		return -1
	}

	sma := movavg.NewSMA(ndays)

	for i := len(vals) - ndays; i < ndays; i++ {

	}

	return sma.Avg()
}
