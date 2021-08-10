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
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gituser143/cryptgo/pkg/utils"
)

func NewCurencyIDMap() CurrencyIDMap {
	c := make(CurrencyIDMap)
	return c
}

func (c *CurrencyIDMap) Populate() {
	url := "https://api.coincap.io/v2/rates"
	method := "GET"

	client := &http.Client{}

	// Create Request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	// Send Request and get response
	res, err := client.Do(req)
	if err != nil {
		res.Body.Close()
		return
	}

	data := utils.AllCurrencyData{}

	// Read response
	err = json.NewDecoder(res.Body).Decode(&data)
	res.Body.Close()
	if err != nil {
		return
	}

	// Iterate over currencies
	for _, curr := range data.Data {
		currencyID := curr.ID
		rate, err := strconv.ParseFloat(curr.RateUSD, 64)
		if err == nil {

			(*c)[currencyID] = Currency{
				Symbol:  fmt.Sprintf("%s %s", curr.Symbol, curr.CurrencySymbol),
				RateUSD: rate,
			}
		}
	}
}
