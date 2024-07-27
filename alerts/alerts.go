package alerts

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/BaboaAtCity/BigNewsMorgan/coingecko"
)

const alertsFile = "alerts.json"

type Alert struct {
	Coin      string    `json:"coin"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type AlertsData struct {
	Alerts []Alert `json:"alerts"`
}

func SaveAlerts(alerts AlertsData) error {
	file, err := os.Create(alertsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For pretty-printing
	return encoder.Encode(alerts)
}

func LoadAlerts() (AlertsData, error) {
	var alerts AlertsData
	file, err := os.Open(alertsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return AlertsData{Alerts: []Alert{}}, nil
		}
		return alerts, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&alerts)
	return alerts, err
}

func AddAlert(coin string, price float64) error {
	coinID, err := coingecko.SearchCoin(coin)
	if err != nil {
		return err
	}

	alerts, err := LoadAlerts()
	if err != nil {
		return err
	}

	newAlert := Alert{
		Coin:      coinID,
		Price:     price,
		CreatedAt: time.Now(),
	}
	fmt.Println("newAlert")
	fmt.Println(newAlert)

	alerts.Alerts = append(alerts.Alerts, newAlert)
	return SaveAlerts(alerts)

}
