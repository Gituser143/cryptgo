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
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

// NewCoinIDMap returns an instance of CoinIDMap
func NewCoinIDMap() CoinIDMap {
	c := make(CoinIDMap)
	return c
}

// Populate updates values into the CoinIDMap
func (c CoinIDMap) Populate() {

	var m sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	// Get CoinCapIDs
	go func(IDMap CoinIDMap, m *sync.Mutex, wg *sync.WaitGroup) {
		defer wg.Done()
		url := "https://api.coincap.io/v2/assets?limit=2000"
		method := "GET"

		client := &http.Client{}

		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return
		}

		res, err := client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		coinCapData := CoinCapData{}

		err = json.NewDecoder(res.Body).Decode(&coinCapData)
		if err != nil {
			return
		}

		for _, val := range coinCapData.Data {
			m.Lock()
			if _, ok := IDMap[val.Symbol]; ok {
				IDMap[val.Symbol] = CoinID{
					CoinCapID:   val.ID,
					CoinGeckoID: IDMap[val.Symbol].CoinGeckoID,
				}
			} else {
				IDMap[val.Symbol] = CoinID{
					CoinCapID: val.ID,
				}
			}
			m.Unlock()
		}
	}(c, &m, &wg)

	go func(IDMap CoinIDMap, m *sync.Mutex, wg *sync.WaitGroup) {
		defer wg.Done()

		coinPtr, err := getTopNCoins(250)
		if err != nil {
			return
		}

		for _, val := range coinPtr {
			symbol := strings.ToUpper(val.Symbol)
			m.Lock()
			if _, ok := IDMap[symbol]; ok {
				IDMap[symbol] = CoinID{
					CoinGeckoID: val.ID,
					CoinCapID:   IDMap[symbol].CoinCapID,
				}
			} else {
				IDMap[symbol] = CoinID{
					CoinGeckoID: val.ID,
				}
			}
			m.Unlock()
		}
	}(c, &m, &wg)

	wg.Wait()
}
