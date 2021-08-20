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
	"strings"
	"time"

	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/gorilla/websocket"
	gecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
)

// GetFavouritePrices gets coin prices for coins specified by favourites.
// This data is returned on the dataChannel.
func GetFavouritePrices(ctx context.Context, favourites map[string]bool, dataChannel chan CoinData) error {

	// Init Client
	geckoClient := gecko.NewClient(nil)

	// Set Parameters
	vsCurrency := "usd"
	order := geckoTypes.OrderTypeObject.MarketCapDesc
	page := 1
	sparkline := true
	priceChangePercentage := []string{}

	return utils.LoopTick(ctx, time.Duration(10)*time.Second, func(errChan chan error) {

		var finalErr error

		favouriteData := make(map[string]float64)

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		// Get Coin IDs
		IDs := []string{}
		for id := range favourites {
			IDs = append(IDs, id)
		}

		perPage := len(IDs)

		// Fetch Data
		coinDataPointer, err := geckoClient.CoinsMarket(vsCurrency, IDs, order, perPage, page, sparkline, priceChangePercentage)
		if err != nil {
			finalErr = err
			return
		}

		// Set Prices
		for _, val := range *coinDataPointer {
			symbol := strings.ToUpper(val.Symbol)
			favouriteData[symbol] = val.CurrentPrice
		}

		// Aggregate data
		coinData := CoinData{
			Type:       "FAVOURITES",
			Favourites: favouriteData,
		}

		// Send data
		select {
		case <-ctx.Done():
			finalErr = ctx.Err()
			return
		case dataChannel <- coinData:
		}

	})
}

// GetCoinHistory gets price history of a coin specified by id, for an interval
// received through the interval channel.
// The default interval is set as 24 Hours
func GetCoinHistory(ctx context.Context, id string, intervalChannel chan string, dataChannel chan CoinData) error {

	intervalToDuration := map[string]string{
		"24hr": "1",
		"7d":   "7",
		"14d":  "14",
		"30d":  "30",
		"90d":  "90",
		"180d": "180",
		"1yr":  "365",
		"5yr":  "1825",
	}

	// Set Default Interval to 1 day
	i := "24hr"

	// Init Client
	geckoClient := gecko.NewClient(nil)

	return utils.LoopTick(ctx, time.Duration(3)*time.Second, func(errChan chan error) {
		var finalErr error

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		select {
		case <-ctx.Done():
			finalErr = ctx.Err()
			return
		case interval := <-intervalChannel:
			// Update interval
			i = interval
		default:
			break
		}

		// Get interval duration and fetch data
		intervalDuration := intervalToDuration[i]
		data, err := geckoClient.CoinsIDMarketChart(id, "usd", intervalDuration)
		if err != nil {
			finalErr = err
			return
		}

		// Aggregate price history
		price := []float64{}
		for _, v := range *data.Prices {
			price = append(price, float64(v[1]))
		}

		// Set max and min
		min := utils.MinFloat64(price...)
		max := utils.MaxFloat64(price...)

		// Clean price for graphs
		for i, val := range price {
			price[i] = val - min
		}

		// Aggregate data
		coinData := CoinData{
			Type:         "HISTORY",
			PriceHistory: price,
			MinPrice:     min,
			MaxPrice:     max,
		}

		// Send Data
		select {
		case <-ctx.Done():
			finalErr = ctx.Err()
			return
		case dataChannel <- coinData:
		}
	})
}

// GetCoinDetails fetches details for a coin specified by id
// and sends the data on dataChannel
func GetCoinDetails(ctx context.Context, id string, dataChannel chan CoinData) error {
	// Init client
	geckoClient := gecko.NewClient(nil)

	// Set Parameters
	localization := false
	tickers := false
	marketData := true
	communityData := false
	developerData := false
	sparkline := false

	return utils.LoopTick(ctx, time.Duration(10)*time.Second, func(errChan chan error) {
		var finalErr error

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		// Fetch Data
		coinData, err := geckoClient.CoinsID(id, localization, tickers, marketData, communityData, developerData, sparkline)
		if err != nil {
			finalErr = err
			return
		}

		// Get Explorer links
		explorerLinks := [][]string{}
		for key, val := range *coinData.Links {
			if key == "blockchain_site" {
				sites := val.([]interface{})
				for _, site := range sites {
					siteStr := site.(string)
					if siteStr != "" {
						explorerLinks = append(explorerLinks, []string{siteStr})
					}
				}
			}
		}

		// Get Total Supply if coin has it
		totalSupply := 0.0
		if coinData.MarketData.TotalSupply != nil {
			totalSupply = *coinData.MarketData.TotalSupply
		}

		// Get Change Percents
		changePercents := [][]string{
			{"24H", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage24h)},
			{"7D", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage7d)},
			{"14D", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage14d)},
			{"30D", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage30d)},
			{"60D", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage60d)},
			{"200D", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage200d)},
			{"1Y", fmt.Sprintf("%.2f", coinData.MarketData.PriceChangePercentage1y)},
		}

		for i, row := range changePercents {
			change := row[1]
			if string(change[0]) == "-" {
				change = utils.DownArrow + " " + change[1:]
			} else {
				change = utils.UpArrow + " " + change
			}
			changePercents[i][1] = change
		}

		// Get ATH, ATL and Last update times
		timeLayout := "2006-01-02T15:04:05.000Z"
		tATHDate, err := time.Parse(timeLayout, coinData.MarketData.ATHDate["usd"])
		if err != nil {
			finalErr = err
			return
		}

		tATLDate, err := time.Parse(timeLayout, coinData.MarketData.ATLDate["usd"])
		if err != nil {
			finalErr = err
			return
		}

		tUpdate, err := time.Parse(timeLayout, coinData.LastUpdated)
		if err != nil {
			finalErr = err
			return
		}

		data := CoinDetails{
			Name:           coinData.Name,
			Symbol:         strings.ToUpper(coinData.Symbol),
			Rank:           fmt.Sprintf("%d", coinData.MarketCapRank),
			BlockTime:      fmt.Sprintf("%d", coinData.BlockTimeInMin),
			MarketCap:      coinData.MarketData.MarketCap["usd"],
			Website:        "",
			Explorers:      explorerLinks,
			ATH:            coinData.MarketData.ATH["usd"],
			ATHDate:        tATHDate.Format(time.RFC822),
			ATL:            coinData.MarketData.ATL["usd"],
			ATLDate:        tATLDate.Format(time.RFC822),
			High24:         coinData.MarketData.High24["usd"],
			Low24:          coinData.MarketData.Low24["usd"],
			TotalVolume:    coinData.MarketData.TotalVolume["usd"],
			ChangePercents: changePercents,
			TotalSupply:    totalSupply,
			CurrentSupply:  coinData.MarketData.CirculatingSupply,
			LastUpdate:     tUpdate.Format(time.RFC822),
		}

		// Aggregate data
		CoinDetails := CoinData{
			Type:    "DETAILS",
			Details: data,
		}

		// Send data
		select {
		case <-ctx.Done():
			finalErr = ctx.Err()
			return
		case dataChannel <- CoinDetails:
		}
	})
}

// GetLivePrice uses a websocket to stream realtime prices of a coin specified
// by id. The prices are sent on the dataChannel
func GetLivePrice(ctx context.Context, id string, dataChannel chan string) error {
	url := fmt.Sprintf("wss://ws.coincap.io/prices?assets=%s", id)
	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	msg := make(map[string]string)

	return utils.LoopTick(ctx, time.Duration(100*time.Millisecond), func(errChan chan error) {
		var finalErr error

		// Defer panic recovery for closed websocket
		defer func() {
			if e := recover(); e != nil {
				finalErr = fmt.Errorf("socket read error")
			}
		}()

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		err = c.ReadJSON(&msg)
		if err != nil {
			finalErr = err
			return
		}

		select {
		case <-ctx.Done():
			finalErr = ctx.Err()
			return
		case dataChannel <- msg[id]:
		}
	})
}
