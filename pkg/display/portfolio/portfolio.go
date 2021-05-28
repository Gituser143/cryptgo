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

package portfolio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Gituser143/cryptgo/pkg/api"
	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

var rows [][]string

type PortfolioTable struct {
	*widgets.Table
}

// NewPortfolioPage creates, initialises and returns a pointer to an instance of PortfolioTable
func NewPortfolioPage() *PortfolioTable {
	p := &PortfolioTable{
		Table: widgets.NewTable(),
	}

	p.Table.Title = " Portfolio "
	p.Table.Header = []string{"Coin", "Symbol", "Price", "Holding", "Balance"}
	p.Table.Rows = rows
	p.Table.CursorColor = ui.ColorCyan
	p.Table.ShowCursor = true
	p.Table.ColWidths = []int{5, 5, 5, 5, 5}
	p.Table.ColResizer = func() {
		x := p.Table.Inner.Dx()
		p.Table.ColWidths = []int{
			x / 5,
			x / 5,
			x / 5,
			x / 5,
			x / 5,
		}
	}
	return p
}

func (p *PortfolioTable) Resize(termWidth, termHeight int) {
	textWidth := 100

	textHeight := len(p.Table.Rows) + 3
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

	p.Table.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget
func (p *PortfolioTable) Draw(buf *ui.Buffer) {
	p.Table.Draw(buf)
}

// Update Portfolio data
func (p *PortfolioTable) UpdateRows(portfolio map[string]float64, currency string, currencyVal float64) {

	var wg sync.WaitGroup
	var m sync.Mutex

	client := &http.Client{}
	method := "GET"

	rows := [][]string{}
	sum := 0.0
	for coin, amt := range portfolio {
		wg.Add(1)
		go func(coin string, amt float64, wg *sync.WaitGroup, m *sync.Mutex) {
			defer wg.Done()

			url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", coin)

			req, err := http.NewRequest(method, url, nil)
			if err != nil {
				return
			}

			res, err := client.Do(req)
			if err != nil {
				return
			}

			defer res.Body.Close()

			data := api.CoinAsset{}

			err = json.NewDecoder(res.Body).Decode(&data)

			if err != nil {
				return
			}

			p, err := strconv.ParseFloat(data.Data.PriceUsd, 64)
			if err != nil {
				return
			}

			row := []string{
				data.Data.Name,
				data.Data.Symbol,
				fmt.Sprintf("%.2f", p/currencyVal),
				fmt.Sprintf("%.6f", amt),
				fmt.Sprintf("%.4f", p*amt/currencyVal),
			}

			m.Lock()
			sum += p * amt / currencyVal
			rows = append(rows, row)
			m.Unlock()

		}(coin, amt, &wg, &m)
	}

	wg.Wait()

	p.Header[2] = fmt.Sprintf("Price (%s)", currency)
	p.Header[4] = fmt.Sprintf("Balance (%s)", currency)
	p.Rows = rows
	p.Title = fmt.Sprintf(" Portfolio: %.4f %s ", sum, currency)
	utils.SortData(p.Rows, 4, false, "PORTFOLIO")
}
