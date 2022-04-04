package query

import (
	"net/http"
	"encoding/json"
	"os"
)

type ValueTracker interface {
	Current(ticker string) float32
}

type FakeValueTracker struct {
	ValueMap map[string]float32
}

func (valueTracker FakeValueTracker) Current(ticker string) float32 {
	return valueTracker.ValueMap[ticker]
}

type FinnHubValueTracker struct {
	apiKey string
}

type FinnHubResponse struct {
	Value float32 `json:"c"`
}

func (valueTracker FinnHubValueTracker) Current(ticker string) float32 {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://finnhub.io/api/v1/quote?symbol=" + ticker, nil)
	req.Header.Set("X-Finnhub-Token", os.Getenv("FINNHUB_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		//error
	}

	defer resp.Body.Close()

	result := FinnHubResponse{}

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Value
}
