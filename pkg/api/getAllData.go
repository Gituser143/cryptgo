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
	"fmt"
	"time"

	"github.com/Gituser143/cryptgo/pkg/utils"
	gecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
)

func getTopNCoins(n int) (geckoTypes.CoinsMarket, error) {
	geckoClient := gecko.NewClient(nil)

	vsCurrency := "usd"
	ids := []string{}

	if n > 1000 {
		return nil, fmt.Errorf("page size limit is 1000")
	}

	perPage := n
	page := 1

	sparkline := true

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

// GetPercentageChangeForDuration returns price change percentage given a
// CoinsMarketItem and a duration, If the specified duration does not exist, 24
// Hour change percent is returned
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

// GetAssets serves data about top 100 coins for the main page
func GetAssets(ctx context.Context, dataChannel chan AssetData, sendData *bool) error {

	return utils.LoopTick(ctx, time.Duration(10)*time.Second, func(errChan chan error) {
		var finalErr error
		var data AssetData

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		if *sendData {
			// Fetch Data
			coinsData, err := getTopNCoins(100)
			if err != nil {
				finalErr = err
				return
			}

			topCoinData := make([][]float64, 3)
			topCoins := make([]string, 3)
			maxPrices := make([]float64, 3)
			minPrices := make([]float64, 3)

			// Set Prices, Max and Min
			for i := 0; i < 3; i++ {
				val := coinsData[i]
				topCoins[i] = val.Name
				topCoinData[i] = val.SparklineIn7d.Price
				maxPrices[i] = utils.MaxFloat64(topCoinData[i]...)
				minPrices[i] = utils.MinFloat64(topCoinData[i]...)

				// Clean data for graph
				for index := range topCoinData[i] {
					topCoinData[i][index] -= minPrices[i]
				}
			}

			// Aggregate data
			data = AssetData{
				AllCoinData: coinsData,
				MaxPrices:   maxPrices,
				MinPrices:   minPrices,
				TopCoinData: topCoinData,
				TopCoins:    topCoins,
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
