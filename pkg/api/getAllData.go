package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gituser143/cryptgo/pkg/utils"
)

type Asset struct {
	Id                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
	Explorer          string `json:"explorer"`
}

type AssetData struct {
	IsTopCoinData bool
	Data          []Asset `json:"data"`
	TimeStamp     uint    `json:"timestamp"`
	TopCoinData   [][]float64
	TopCoins      []string
}

type CoinPrice struct {
	Price     string `json:"priceUsd"`
	Timestamp uint   `json:"time"`
}

type CoinHistory struct {
	Data      []CoinPrice `json:"data"`
	Timestamp uint        `json:"timestamp"`
}

func GetAssets(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {
	url := "https://api.coincap.io/v2/assets"
	method := "GET"

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	return utils.LoopTick(ctx, time.Duration(1)*time.Second, func() error {
		data := AssetData{}

		if *sendData {
			// Send Request
			res, err := client.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			// Read response
			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case dataChannel <- data:
			}
		}

		return nil
	})
}

func GetTopCoinData(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {
	url := "https://api.coincap.io/v2/assets?limit=3"
	method := "GET"

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	return utils.LoopTick(ctx, time.Duration(10)*time.Second, func() error {
		data := AssetData{}

		if *sendData {

			// Send Request
			res, err := client.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			// Read response
			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				return err
			}

			topCoinData := make([][]float64, 3)
			topCoins := make([]string, 3)

			for i, val := range data.Data {
				historyUrl := fmt.Sprintf("https://api.coincap.io/v2/assets/%s/history?interval=d1", val.Id)

				// Create Request
				req, err := http.NewRequest(method, historyUrl, nil)
				if err != nil {
					return err
				}

				res, err := client.Do(req)
				if err != nil {
					return err
				}
				defer res.Body.Close()

				historyData := CoinHistory{}

				// Read response
				err = json.NewDecoder(res.Body).Decode(&historyData)
				if err != nil {
					return err
				}

				price := []float64{}
				for _, v := range historyData.Data {
					p, err := strconv.ParseFloat(v.Price, 64)
					if err != nil {
						return err
					}

					price = append(price, p)
				}

				topCoinData[i] = price
				topCoins[i] = val.Name
			}

			data.TopCoinData = topCoinData
			data.TopCoins = topCoins
			data.IsTopCoinData = true

			select {
			case <-ctx.Done():
				return ctx.Err()
			case dataChannel <- data:
			}
		}

		return nil
	})
}
