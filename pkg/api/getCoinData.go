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
	"sync"
	"time"

	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/gorilla/websocket"
)

// API Documentation can be found at https://docs.coincap.io/

// CoinData Holds data pertaining to a single coin.
// This is used to serve per coin details.
// It additionally holds a map of favourite coins.
type CoinData struct {
	Type          string
	PriceHistory  []float64
	MinPrice      float64
	MaxPrice      float64
	CoinAssetData CoinAsset
	Price         string
	Favourites    map[string]float64
}

// CoinAsset holds Asset data for a single coin
type CoinAsset struct {
	Data      Asset `json:"data"`
	TimeStamp uint  `json:"timestamp"`
}

// GetFavouritePrices gets coin prices for coins specified by favourites.
// This data is returned on the dataChannel.
func GetFavouritePrices(ctx context.Context, favourites map[string]bool, dataChannel chan CoinData) error {
	method := "GET"

	// Init Client
	client := &http.Client{}

	return utils.LoopTick(ctx, time.Duration(1)*time.Second, func() error {

		var wg sync.WaitGroup
		var m sync.Mutex

		favouriteData := make(map[string]float64)

		// Iterate over favorite coins
		for id := range favourites {
			wg.Add(1)
			go func(id string, wg *sync.WaitGroup, m *sync.Mutex) {

				price := 0.0
				data := CoinAsset{}

				// In case errorred in later stages, below deferred function call
				// is done to make sure coin does not vanish from favourites
				defer func() {
					m.Lock()
					favouriteData[data.Data.Symbol] = price
					m.Unlock()
				}()

				defer wg.Done()

				url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id)

				// Create Request
				req, err := http.NewRequestWithContext(ctx, method, url, nil)
				if err != nil {
					return
				}

				// Send Request
				res, err := client.Do(req)
				if err != nil {
					return
				}
				defer res.Body.Close()

				// Read response
				err = json.NewDecoder(res.Body).Decode(&data)
				if err != nil {
					return
				}

				// Get price
				price, err = strconv.ParseFloat(data.Data.PriceUsd, 64)
				if err != nil {
					price = 0
					return
				}

			}(id, &wg, &m)
		}

		wg.Wait()

		// Aggregate data
		coinData := CoinData{
			Type:       "FAVOURITES",
			Favourites: favouriteData,
		}

		// Send data
		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- coinData:
		}

		return nil
	})
}

// GetCoinHistory gets price history of a coin specified by id, for an interval
// received through the interval channel.
func GetCoinHistory(ctx context.Context, id string, intervalChannel chan string, dataChannel chan CoinData) error {
	method := "GET"

	// Init Client
	client := &http.Client{}

	// Set Default Interval to 1 day
	i := "d1"

	return utils.LoopTick(ctx, time.Duration(3)*time.Second, func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case interval := <-intervalChannel:
			// Update interval
			i = interval
		default:
			break
		}

		url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s/history?interval=%s", id, i)
		data := CoinHistory{}

		// Create Request
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return err
		}

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

		// Aggregate price history
		price := []float64{}
		for _, v := range data.Data {
			p, err := strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return err
			}

			price = append(price, p)
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
			return ctx.Err()
		case dataChannel <- coinData:
		}

		return nil
	})
}

// GetCoinAsset fetches asset data for a coin specified by id
// and sends the data on dataChannel
func GetCoinAsset(ctx context.Context, id string, dataChannel chan CoinData) error {
	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s/", id)
	method := "GET"

	// Init client
	client := &http.Client{}

	// Create Request
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	return utils.LoopTick(ctx, time.Duration(3)*time.Second, func() error {
		data := CoinAsset{}

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

		// Aggregate data
		coinData := CoinData{
			Type:          "ASSET",
			CoinAssetData: data,
		}

		// Send data
		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- coinData:
		}

		return nil
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

	return utils.LoopTick(ctx, time.Duration(100*time.Millisecond), func() error {
		// Defer panic recovery for closed websocket
		var err error
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("socker read error")
			}
		}()

		err = c.ReadJSON(&msg)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- msg[id]:
		}

		return nil
	})
}
