/*
Copyright © 2021 Bhargav SNV bhargavsnv100@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gituser143/cryptgo/pkg/utils"
	gecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
)

// API Documentation can be found at https://docs.coincap.io/

// Asset holds details of a single coin
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

// AssetData is used to hold details of multiple coins and the price history
// of top ranked coins along with their names
type AssetData struct {
	IsTopCoinData bool
	Data          []Asset `json:"data"`
	TimeStamp     uint    `json:"timestamp"`
	TopCoinData   [][]float64
	TopCoins      []string
}

// CoinPrice holds the price of a coin at a given time
type CoinPrice struct {
	Price     string `json:"priceUsd"`
	Timestamp uint   `json:"time"`
}

// CoinHistory holds a slice of CoinPrice, as history of coin value
type CoinHistory struct {
	Data      []CoinPrice `json:"data"`
	Timestamp uint        `json:"timestamp"`
}

func GetTopNCoinsFromCoinGecko(n int) ([]string, error) {
	geckoClient := gecko.NewClient(nil)

	vsCurrency := "usd"
	ids := []string{}

	if n > 1000 {
		return nil, fmt.Errorf("page size limit is 1000")
	}

	perPage := n
	page := 1

	sparkline := false
	priceChangePercentage := []string{}

	order := geckoTypes.OrderTypeObject.MarketCapDesc
	coinDataPointer, err := geckoClient.CoinsMarket(vsCurrency, ids, order, perPage, page, sparkline, priceChangePercentage)

	if err != nil {
		return nil, err
	}

	coinData := *coinDataPointer

	topNIds := []string{}

	for i := 0; i < n; i += 1 {
		coinId := coinData[i].ID
		topNIds = append(topNIds, coinId)
	}

	return topNIds, nil
}

// Get Assets contacts the 'api.coincap.io/v2/assets' endpoint to get asset
// information of top 100 coins. It then serves this information through the
// dataChannel
func GetAssets(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {
	url := "https://api.coincap.io/v2/assets"
	method := "GET"

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Init Client
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

			// Send Data
			select {
			case <-ctx.Done():
				return ctx.Err()
			case dataChannel <- data:
			}
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		return nil
	})
}

// GetTopCoinData first fetches the top 3 ranked coins by contacting the
// 'api.coincap.io/v2/assets' endpoint. Following which, the history of each of
// these coin is queried to the endpoint
// 'api.coincap.io/v2/assets/{id}/history'. This history data is served
//  on the dataChannel
func GetTopCoinData(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {

	topThreeIds, err := GetTopNCoinsFromCoinGecko(3)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.coincap.io/v2/assets?ids=%s,%s,%s", topThreeIds[0], topThreeIds[1], topThreeIds[2])
	method := "GET"

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Init Client
	client := &http.Client{}

	return utils.LoopTick(ctx, time.Duration(5)*time.Second, func() error {
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

				// Fetch History
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

				// Aggregate price
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

			// Aggregate data
			data.TopCoinData = topCoinData
			data.TopCoins = topCoins
			data.IsTopCoinData = true

			// Send data
			select {
			case <-ctx.Done():
				return ctx.Err()
			case dataChannel <- data:
			}
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		return nil
	})
}
