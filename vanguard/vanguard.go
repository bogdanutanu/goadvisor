package vanguard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	uri    *url.URL
}

// Vanguard Client for Vanguard UK funds
type Vanguard interface {
	FetchLastNDays(days uint16) valuations
	Decide() bool // this method should belong to a generic interface
}

// New creates a new Vanguard client
func New(vanguardURL string) (Vanguard, error) {
	uri, err := url.Parse(vanguardURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing '%s' as URL for Vanguard: %w", vanguardURL, err)
	}

	return &vanguard{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		uri: uri,
	}, nil
}

func (v *vanguard) FetchLastNDays(days uint16) valuations {

	now := time.Now()
	nowDateStr := now.Format(dateLayout)
	twoHundredDaysAgo := now.AddDate(0, 0, -200)
	twoHundredDaysAgoStr := twoHundredDaysAgo.Format(dateLayout)

	vars := fmt.Sprintf("portId:9244,issueType:S,startDate:%s,endDate:%s", twoHundredDaysAgoStr, nowDateStr)
	v.uri.Query().Add("vars", vars)
	log.Infof("Will query '%s'", v.uri.String())
	resp, err := v.client.Get(v.uri.String())
	if err != nil {
		log.Errorf("Error making GET '%s': %v", v.uri.String(), err)
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
		sma.Add(vals[i].NavPrice)
	}

	return sma.Avg()
}

func (v *vanguard) Decide() bool {
	vals := v.FetchLastNDays(365)
	mAvg200 := computeMovingAverage(vals, 200)
	mAvg50 := computeMovingAverage(vals, 50)
	log.Infof("Mavg200: %f Mavg50: %f", mAvg200, mAvg50)

	return mAvg50 > mAvg200
}
