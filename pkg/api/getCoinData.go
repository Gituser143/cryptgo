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

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

// API Documentation can be found at https://docs.coincap.io/

// CoinData Holds data pertaining to a single coin.
// This is used to serve per coin details.
// It additionally holds a map of favourite coins.
type CoinData struct {
	Type         string
	PriceHistory []float64
	MinPrice     float64
	MaxPrice     float64
	Details      CoinDetails
	Favourites   map[string]float64
}

// CoinAsset holds Asset data for a single coin
type CoinAsset struct {
	Data      Asset `json:"data"`
	TimeStamp uint  `json:"timestamp"`
}

type CoinDetails struct {
	Name           string
	Symbol         string
	Rank           string
	BlockTime      string
	MarketCap      float64
	Website        string
	Explorers      [][]string
	ATH            float64
	ATHDate        string
	ATL            float64
	ATLDate        string
	TotalVolume    float64
	ChangePercents [][]string
	TotalSupply    float64
	CurrentSupply  float64
	LastUpdate     string
}

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

		IDs := []string{}

		for id := range favourites {
			IDs = append(IDs, id)
		}

		perPage := len(IDs)

		coinDataPointer, err := geckoClient.CoinsMarket(vsCurrency, IDs, order, perPage, page, sparkline, priceChangePercentage)
		if err != nil {
			finalErr = err
			return
		}

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

	geckoClient := gecko.NewClient(nil)

	return utils.LoopTick(ctx, time.Duration(3)*time.Second, func(errChan chan error) {
		var finalErr error = nil

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

// GetCoinAsset fetches asset data for a coin specified by id
// and sends the data on dataChannel
func GetCoinAsset(ctx context.Context, id string, dataChannel chan CoinData) error {
	// Init client
	geckoClient := gecko.NewClient(nil)
	localization := false
	tickers := false
	marketData := true
	communityData := false
	developerData := false
	sparkline := false

	return utils.LoopTick(ctx, time.Duration(10)*time.Second, func(errChan chan error) {
		var finalErr error = nil

		defer func() {
			if finalErr != nil {
				errChan <- finalErr
			}
		}()

		coinData, err := geckoClient.CoinsID(id, localization, tickers, marketData, communityData, developerData, sparkline)
		if err != nil {
			finalErr = err
			return
		}

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

		totalSupply := 0.0
		if coinData.MarketData.TotalSupply != nil {
			totalSupply = *coinData.MarketData.TotalSupply
		}

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
				change = DOWN_ARROW + " " + change[1:]
			} else {
				change = UP_ARROW + " " + change
			}
			changePercents[i][1] = change
		}

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
		var finalErr error = nil

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
