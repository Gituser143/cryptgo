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

package utilitywidgets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

// Currency holds data of a currency
type Currency struct {
	ID             string `json:"id"`
	Symbol         string `json:"symbol"`
	CurrencySymbol string `json:"currencySymbol"`
	Type           string `json:"type"`
	RateUSD        string `json:"rateUSD"`
}

// CurrencyData is used to hold data of a currency when fetched from the API
type CurrencyData struct {
	Data      Currency `json:"data"`
	Timestamp uint     `json:"timestamp"`
}

// AllCurrencyData holds details of currencies when all are fetched from the API
type AllCurrencyData struct {
	Data      []Currency `json:"data"`
	Timestamp uint       `json:"timestamp"`
}

// CurrencyTable is a widget used to display currencyies, symbols and rates
type CurrencyTable struct {
	*widgets.Table
}

// NewCurrencyPage creates, initialises and returns a pointer to an instance of CurrencyTable
func NewCurrencyPage() *CurrencyTable {
	c := &CurrencyTable{
		Table: widgets.NewTable(),
	}

	c.Table.Title = " Select Currency "
	c.Table.Header = []string{"Currency", "Symbol", "Type", "USD rate"}
	c.Table.CursorColor = ui.ColorCyan
	c.Table.ShowCursor = true
	c.Table.ColWidths = []int{5, 5, 5, 5}
	c.Table.ColResizer = func() {
		x := c.Table.Inner.Dx()
		c.Table.ColWidths = []int{
			x / 4,
			x / 4,
			x / 4,
			x / 4,
		}
	}
	return c
}

func (c *CurrencyTable) Resize(termWidth, termHeight int) {
	textWidth := 50

	textHeight := len(c.Table.Rows) + 3
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2
	if x < 0 {
		x = 0
		textWidth = termWidth
	}
	if y < 0 {
		y = 0
		textHeight = termHeight
	}

	c.Table.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget
func (c *CurrencyTable) Draw(buf *ui.Buffer) {
	c.Table.Draw(buf)
}

// UpdateAll fetches rates of all currencies and updates them as rows in the table
func (c *CurrencyTable) UpdateAll() {
	url := "https://api.coincap.io/v2/rates"
	method := "GET"

	rows := [][]string{}

	// init client
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

	data := AllCurrencyData{}

	// Read response
	err = json.NewDecoder(res.Body).Decode(&data)
	res.Body.Close()
	if err != nil {
		return
	}

	// Iterate over currencies
	for _, currency := range data.Data {
		// Get currency rate
		rate, err := strconv.ParseFloat(currency.RateUSD, 64)
		if err != nil {
			continue
		}

		// Aggregate data
		row := []string{
			currency.Symbol,
			currency.CurrencySymbol,
			currency.Type,
			fmt.Sprintf("%.4f", rate),
		}

		rows = append(rows, row)
	}

	// Update table rows and sort alphabetically
	c.Table.Rows = rows
	utils.SortData(c.Table.Rows, 0, true, "CURRENCY")
}

// Update Rows fetches popular currency details and updates the table
func (c *CurrencyTable) UpdateRows() {
	currencies := []string{
		"united-states-dollar",
		"euro",
		"japanese-yen",
		"british-pound-sterling",
		"indian-rupee",
		"australian-dollar",
		"canadian-dollar",
		"chinese-yuan-renminbi",
	}

	var wg sync.WaitGroup
	var m sync.Mutex

	// init client
	client := &http.Client{}
	method := "GET"

	rows := [][]string{}

	// iterate over currencies
	for _, currency := range currencies {
		wg.Add(1)
		// make concurrent requests
		go func(c string, wg *sync.WaitGroup, m *sync.Mutex) {
			defer wg.Done()
			url := fmt.Sprintf("https://api.coincap.io/v2/rates/%s", c)

			// Create request
			req, err := http.NewRequest(method, url, nil)
			if err != nil {
				return
			}

			// Get response
			res, err := client.Do(req)
			if err != nil {
				return
			}
			defer res.Body.Close()

			data := CurrencyData{}

			// Read response
			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				return
			}

			// Get rate
			rate, err := strconv.ParseFloat(data.Data.RateUSD, 64)
			if err != nil {
				return
			}

			// Aggregate data
			row := []string{
				data.Data.Symbol,
				data.Data.CurrencySymbol,
				data.Data.Type,
				fmt.Sprintf("%.4f", rate),
			}

			m.Lock()
			rows = append(rows, row)
			m.Unlock()
		}(currency, &wg, &m)
	}

	wg.Wait()

	// Update table rows
	c.Table.Rows = rows
	utils.SortData(c.Table.Rows, 0, true, "CURRENCY")
}
