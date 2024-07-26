package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	apiURL    = "https://api.coingecko.com/api/v3/simple/price?include_24hr_vol=true&include_24hr_change=true&ids=%s&vs_currencies=usd"
	searchURL = "https://api.coingecko.com/api/v3/search?query=%s"
)

// ?include_24hr_vol=true
func GetPrices(coins []string) (map[string]map[string]float64, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf(apiURL, strings.Join(coins, ","))

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}
	fmt.Println("prices log---------")

	fmt.Println(prices)
	return prices, nil
}

func SearchCoin(query string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf(searchURL, query)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var searchResult struct {
		Coins []struct {
			ID string `json:"id"`
		} `json:"coins"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return "", err
	}

	if len(searchResult.Coins) > 0 {
		return searchResult.Coins[0].ID, nil
	}

	return "", nil
}
