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

type CoinData struct {
	Type          string
	PriceHistory  []float64
	CoinAssetData CoinAsset
	Price         string
	Favourites    [][]string
}

type CoinAsset struct {
	Data      Asset `json:"data"`
	TimeStamp uint  `json:"timestamp"`
}

func GetFavouritePrices(ctx context.Context, favourites map[string]bool, dataChannel chan CoinData) error {
	method := "GET"

	client := &http.Client{}

	return utils.LoopTick(ctx, time.Duration(1)*time.Second, func() error {

		var wg sync.WaitGroup
		var m sync.Mutex

		favouriteData := [][]string{}

		for id := range favourites {
			wg.Add(1)
			go func(id string, wg *sync.WaitGroup, m *sync.Mutex) {
				data := CoinAsset{}
				url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id)

				// Create Request
				req, err := http.NewRequest(method, url, nil)
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

				m.Lock()
				favouriteData = append(favouriteData, []string{
					data.Data.Symbol,
					data.Data.PriceUsd,
				})
				m.Unlock()
				wg.Done()
			}(id, &wg, &m)
		}

		wg.Wait()

		coinData := CoinData{
			Type:       "FAVOURITES",
			Favourites: favouriteData,
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- coinData:
		}

		return nil
	})
}

func GetCoinHistory(ctx context.Context, id string, interval *string, dataChannel chan CoinData) error {
	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s/history?interval=%s", id, *interval)
	method := "GET"

	client := &http.Client{}

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	return utils.LoopTick(ctx, time.Duration(3)*time.Second, func() error {
		data := CoinHistory{}

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

		price := []float64{}
		for _, v := range data.Data {
			p, err := strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return err
			}

			price = append(price, p)
		}

		coinData := CoinData{
			Type:         "HISTORY",
			PriceHistory: price,
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- coinData:
		}

		return nil
	})
}

func GetCoinAsset(ctx context.Context, id string, dataChannel chan CoinData) error {
	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s/", id)
	method := "GET"

	client := &http.Client{}

	// Create Request
	req, err := http.NewRequest(method, url, nil)
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

		coinData := CoinData{
			Type:          "ASSET",
			CoinAssetData: data,
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- coinData:
		}

		return nil
	})
}

func GetLivePrice(ctx context.Context, id string, dataChannel chan string) error {
	url := fmt.Sprintf("wss://ws.coincap.io/prices?assets=%s", id)
	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	msg := make(map[string]string)

	for {
		err := c.ReadJSON(&msg)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- msg[id]:
		}
	}
}
