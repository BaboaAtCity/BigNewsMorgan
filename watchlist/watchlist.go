package watchlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/BaboaAtCity/BigNewsMorgan/coingecko"
)

const watchlistFile = "watchlist.json"

var (
	watchlist []string
	mutex     sync.RWMutex
)

func init() {
	loadWatchlist()
}

func loadWatchlist() {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := ioutil.ReadFile(watchlistFile)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, initialize with default values
			watchlist = []string{"bitcoin", "ethereum", "solana", "ripple"}
			saveWatchlist()
		} else {
			fmt.Printf("Error reading watchlist file: %v\n", err)
		}
		return
	}

	err = json.Unmarshal(data, &watchlist)
	if err != nil {
		fmt.Printf("Error unmarshaling watchlist: %v\n", err)
	}
}

func saveWatchlist() {
	data, err := json.Marshal(watchlist)
	if err != nil {
		fmt.Printf("Error marshaling watchlist: %v\n", err)
		return
	}

	err = ioutil.WriteFile(watchlistFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing watchlist file: %v\n", err)
	}
}

func Get() []string {
	mutex.RLock()
	defer mutex.RUnlock()
	return watchlist
}

func Add(query string) string {
	coinID, err := coingecko.SearchCoin(query)
	if err != nil {
		return fmt.Sprintf("Error searching for coin: %v", err)
	}

	if coinID == "" {
		return fmt.Sprintf("Could not find a coin matching '%s'", query)
	}

	mutex.Lock()
	defer mutex.Unlock()

	for _, coin := range watchlist {
		if coin == coinID {
			return fmt.Sprintf("%s is already in the watchlist", strings.Title(coinID))
		}
	}

	watchlist = append(watchlist, coinID)
	saveWatchlist()
	return fmt.Sprintf("Added %s to the watchlist", strings.Title(coinID))
}

func Remove(query string) string {
	coinID, err := coingecko.SearchCoin(query)
	if err != nil {
		return fmt.Sprintf("Error searching for coin: %v", err)
	}

	if coinID == "" {
		return fmt.Sprintf("Could not find a coin matching '%s'", query)
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, coin := range watchlist {
		if coin == coinID {
			watchlist = append(watchlist[:i], watchlist[i+1:]...)
			saveWatchlist()
			return fmt.Sprintf("Removed %s from the watchlist", strings.Title(coinID))
		}
	}

	return fmt.Sprintf("%s is not in the watchlist", strings.Title(coinID))
}

func FormatPrices(prices map[string]map[string]float64) string {
	mutex.RLock()
	defer mutex.RUnlock()

	var result strings.Builder
	result.WriteString("Current prices:\n")
	for _, coin := range watchlist {
		if priceData, ok := prices[coin]; ok {
			if price, hasPrice := priceData["usd"]; hasPrice {
				if volume, hasVolume := priceData["usd_24h_vol"]; hasVolume {
					if change, hasChange := priceData["usd_24h_change"]; hasChange {
						// Both price and volume are available
						volume = volume / 1000000
						result.WriteString(fmt.Sprintf("\n%s: $%.2f \n(24h Change: %.2f%%) (24h Vol(mill): $%.2f) \n",
							strings.Title(coin), price, change, volume))
					} else {
						result.WriteString(fmt.Sprintf("%s: $%.2f (24h Change: N/A)\n",
							strings.Title(coin), price))
					}

				} else {
					// Only price is available
					result.WriteString(fmt.Sprintf("%s: $%.2f (24h Vol: N/A)\n",
						strings.Title(coin), price))
				}
			} else {
				// Neither price nor volume is available
				result.WriteString(fmt.Sprintf("%s: Price not available\n",
					strings.Title(coin)))
			}
		} else {
			// Coin data not found
			result.WriteString(fmt.Sprintf("%s: Data not available\n",
				strings.Title(coin)))
		}
	}
	return result.String()
}
