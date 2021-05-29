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
	"sync"
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
	MaxPrices     []float64
	MinPrices     []float64
	TopCoins      []string
	AllCoinData   geckoTypes.CoinsMarket
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

func GetTopNCoinsFromCoinGecko(n int) (geckoTypes.CoinsMarket, error) {
	geckoClient := gecko.NewClient(nil)

	vsCurrency := "usd"
	ids := []string{}

	if n > 1000 {
		return nil, fmt.Errorf("page size limit is 1000")
	}

	perPage := n
	page := 1

	sparkline := false

	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{pcp.PCP1h, pcp.PCP24h, pcp.PCP7d, pcp.PCP14d, pcp.PCP30d, pcp.PCP200d, pcp.PCP1y}

	order := geckoTypes.OrderTypeObject.MarketCapDesc
	coinDataPointer, err := geckoClient.CoinsMarket(vsCurrency, ids, order, perPage, page, sparkline, priceChangePercentage)

	if err != nil {
		return nil, err
	}

	coinData := *coinDataPointer

	return coinData, nil
}

func GetTopNCoinIdsFromCoinGecko(n int) (map[int][]string, error) {

	coinData, err := GetTopNCoinsFromCoinGecko(n)
	if err != nil {
		return nil, err
	}

	data := make(map[int][]string)

	for i := 0; i < n; i += 1 {
		coinId := coinData[i].ID
		coinName := coinData[i].Name
		data[i] = []string{coinId, coinName}
	}

	return data, nil
}

func GetPercentageChangeForDuration(coinData geckoTypes.CoinsMarketItem, duration string) float64 {

	m := map[string]*float64{
		"1h":   coinData.PriceChangePercentage1hInCurrency,
		"24h":  coinData.PriceChangePercentage24hInCurrency,
		"7d":   coinData.PriceChangePercentage7dInCurrency,
		"14d":  coinData.PriceChangePercentage14dInCurrency,
		"30d":  coinData.PriceChangePercentage30dInCurrency,
		"200d": coinData.PriceChangePercentage200dInCurrency,
		"1y":   coinData.PriceChangePercentage1yInCurrency,
	}

	if percentageDuration, isPresent := m[duration]; isPresent && percentageDuration != nil {
		return *percentageDuration
	}
	return coinData.PriceChangePercentage24h
}

// Get Assets contacts the 'https://api.coingecko.com/api/v3/coins/markets' endpoint
// to get asset information of top 100 coins. It then serves this information through the
// dataChannel
func GetAssets(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {

	return utils.LoopTick(ctx, time.Duration(10)*time.Second, func(errChan chan error) {
		var finalErr error = nil
		data := AssetData{}

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		if *sendData {

			coinsData, err := GetTopNCoinsFromCoinGecko(100)
			data.AllCoinData = coinsData
			if err != nil {
				finalErr = err
				return
			}

			// Send Data
			select {
			case <-ctx.Done():
				finalErr = ctx.Err()
				return
			case dataChannel <- data:
			}
		} else {
			select {
			case <-ctx.Done():
				finalErr = ctx.Err()
				return
			default:
			}
		}
	})
}

// GetTopCoinData first fetches the top 3 ranked coins by contacting the
// 'api.coincap.io/v2/assets' endpoint. Following which, the history of each of
// these coin is queried to the endpoint
// 'api.coincap.io/v2/assets/{id}/history'. This history data is served
//  on the dataChannel
func GetTopCoinData(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {

	// Init Client
	client := &http.Client{}
	geckoClient := gecko.NewClient(client)

	return utils.LoopTick(ctx, time.Duration(1)*time.Minute, func(errChan chan error) {
		var finalErr error = nil
		data := AssetData{}

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		if *sendData {

			topIDs, err := GetTopNCoinIdsFromCoinGecko(3)
			if err != nil {
				finalErr = err
				return
			}

			topCoinData := make([][]float64, 3)
			topCoins := make([]string, 3)
			maxPrices := make([]float64, 3)
			minPrices := make([]float64, 3)

			var wg sync.WaitGroup
			var m sync.Mutex

			for i, coin := range topIDs {
				id := coin[0]
				name := coin[1]

				wg.Add(1)

				go func(id, name string, index int, wg *sync.WaitGroup, m *sync.Mutex) {
					defer wg.Done()
					data, err := geckoClient.CoinsIDMarketChart(id, "usd", "7")
					if err != nil {
						finalErr = err
						return
					}

					price := []float64{}
					prices := *data.Prices
					max := float64(prices[0][1])
					min := float64(prices[1][1])
					for _, val := range *data.Prices {
						p := float64(val[1])
						if p > max {
							max = p
						}
						if p < min {
							min = p
						}
						price = append(price, p)
					}

					// Clean prices
					for i, val := range price {
						price[i] = val - min
					}

					m.Lock()
					maxPrices[index] = max
					minPrices[index] = min
					topCoinData[index] = price
					topCoins[index] = name
					m.Unlock()
				}(id, name, i, &wg, &m)
			}

			wg.Wait()

			// Aggregate data
			data.MaxPrices = maxPrices
			data.MinPrices = minPrices
			data.TopCoinData = topCoinData
			data.TopCoins = topCoins
			data.IsTopCoinData = true

			// Send data
			select {
			case <-ctx.Done():
				finalErr = ctx.Err()
				return
			case dataChannel <- data:
			}
		} else {
			select {
			case <-ctx.Done():
				finalErr = ctx.Err()
				return
			default:
			}
		}
	})
}

func GetTopNCoinSymbolToIDMap(n int) (map[string]string, error) {

	coinToSymbolMap := map[string]string{}

	url := fmt.Sprintf("https://api.coincap.io/v2/assets?limit=%d", n)
	method := "GET"

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Init Client
	client := &http.Client{}
	data := AssetData{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	for _, coin := range data.Data {
		coinToSymbolMap[coin.Symbol] = coin.Id
	}

	return coinToSymbolMap, nil
}
