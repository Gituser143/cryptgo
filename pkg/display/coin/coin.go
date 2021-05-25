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

package coin

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	c "github.com/Gituser143/cryptgo/pkg/display/currency"
	"github.com/Gituser143/cryptgo/pkg/display/portfolio"
	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

// DisplayCoin displays the per coin values and details along with a favourites table. It uses the same uiEvents channel as the root page
func DisplayCoin(
	ctx context.Context,
	id string,
	coinIDs map[string]string,
	intervalChannel chan string,
	dataChannel chan api.CoinData,
	priceChannel chan string,
	uiEvents <-chan ui.Event) error {

	defer ui.Clear()

	// Init Coin page
	myPage := NewCoinPage()

	// variables for currency
	currency := "USD $"
	currencyVal := 1.0
	selectCurrency := false
	currencyWidget := c.NewCurrencyPage()

	// Selection of default table
	selectedTable := myPage.IntervalTable
	selectedTable.ShowCursor = true

	previousKey := ""

	// variables to sort favourites table
	favSortIdx := -1
	favSortAsc := false
	favHeader := []string{
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
	}

	// Initialise portfolio
	favourites := utils.GetFavourites()
	portfolioMap := utils.GetPortfolio()
	defer utils.SaveMetadata(favourites, currency, portfolioMap)

	// Initiliase Portfolio Table
	portfolioTable := portfolio.NewPortfolioPage()
	portfolioSelected := false

	// intervals holds interval mappings to be used in the call to the history API
	intervals := map[string]string{
		"1  min":  "m1",
		"5  min":  "m5",
		"15 min":  "m15",
		"30 min":  "m30",
		"1  hour": "h1",
		"2  hour": "h2",
		"6  hour": "h6",
		"12 hour": "h12",
		"1  day":  "d1",
	}

	// Initialise help menu
	help := widgets.NewHelpMenu()
	help.SelectHelpMenu("COIN")
	helpSelected := false

	// UpdateUI to refresh UI
	updateUI := func() {
		// Get Terminal Dimensions
		w, h := ui.TerminalDimensions()

		// Adjust Suuply chart Bar graph values
		myPage.SupplyChart.BarGap = ((w / 3) - (2 * myPage.SupplyChart.BarWidth)) / 2

		myPage.Grid.SetRect(0, 0, w, h)

		// Clear UI
		ui.Clear()

		// Render required widgets
		if helpSelected {
			help.Resize(w, h)
			ui.Render(help)
		} else if portfolioSelected {
			portfolioTable.Resize(w, h)
			ui.Render(portfolioTable)
		} else if selectCurrency {
			currencyWidget.Resize(w, h)
			ui.Render(currencyWidget)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	// Render empty UI
	updateUI()

	// Create ticker to periodically refresh UI
	t := time.NewTicker(time.Duration(1) * time.Second)
	tick := t.C

	for {
		select {
		case <-ctx.Done(): // Context cancelled, exit
			return ctx.Err()

		case e := <-uiEvents: // keyboard events
			switch e.ID {
			case "<Escape>":
				if !helpSelected && !selectCurrency && !portfolioSelected {
					selectCurrency = false
					return fmt.Errorf("UI Closed")
				}
			case "q", "<C-c>":
				ui.Clear()
				return fmt.Errorf("coin UI Closed")

			case "<Resize>":
				updateUI()

			case "?":
				helpSelected = !helpSelected
				updateUI()

			case "c":
				if !helpSelected && !portfolioSelected {
					selectCurrency = true
					currencyWidget.UpdateRows()
					updateUI()
				}

			case "C":
				if !helpSelected && !portfolioSelected {
					selectCurrency = true
					currencyWidget.UpdateAll()
				}
				updateUI()

			case "f":
				if !helpSelected && !portfolioSelected && !selectCurrency {
					selectedTable.ShowCursor = false
					selectedTable = myPage.FavouritesTable
				}

			case "F":
				if !helpSelected && !portfolioSelected && !selectCurrency {
					selectedTable.ShowCursor = false
					selectedTable = myPage.IntervalTable
				}

			case "P":
				if !helpSelected && !selectCurrency {
					portfolioTable.UpdateRows(portfolioMap, currency, currencyVal)
					portfolioSelected = !portfolioSelected
					updateUI()
				}

			}
			if helpSelected {
				switch e.ID {
				case "?":
					updateUI()
				case "<Escape>":
					helpSelected = false
					updateUI()
				case "j", "<Down>":
					help.List.ScrollDown()
					ui.Render(help)
				case "k", "<Up>":
					help.List.ScrollUp()
					ui.Render(help)
				}
			} else if portfolioSelected {
				switch e.ID {
				case "<Escape>":
					portfolioSelected = false
					updateUI()
				case "j", "<Down>":
					portfolioTable.ScrollDown()
					ui.Render(portfolioTable)
				case "k", "<Up>":
					portfolioTable.ScrollUp()
					ui.Render(portfolioTable)
				case "e":
					id := ""
					symbol := ""

					// Get ID and symbol
					if portfolioTable.SelectedRow < len(portfolioTable.Rows) {
						row := portfolioTable.Rows[portfolioTable.SelectedRow]
						symbol = row[1]
					}

					id = coinIDs[symbol]

					if id != "" {
						inputStr := widgets.DrawEdit(uiEvents, symbol)
						amt, err := strconv.ParseFloat(inputStr, 64)

						if err == nil {
							if amt > 0 {
								portfolioMap[id] = amt
							} else {
								delete(portfolioMap, id)
							}
						}
					}

					portfolioTable.UpdateRows(portfolioMap, currency, currencyVal)
				}
			} else if selectCurrency {
				switch e.ID {
				case "j", "<Down>":
					currencyWidget.ScrollDown()
				case "k", "<Up>":
					currencyWidget.ScrollUp()
				case "<C-d>":
					currencyWidget.ScrollHalfPageDown()
				case "<C-u>":
					currencyWidget.ScrollHalfPageUp()
				case "<C-f>":
					currencyWidget.ScrollPageDown()
				case "<C-b>":
					currencyWidget.ScrollPageUp()
				case "g":
					if previousKey == "g" {
						currencyWidget.ScrollTop()
					}
				case "<Home>":
					currencyWidget.ScrollTop()
				case "G", "<End>":
					currencyWidget.ScrollBottom()
				case "<Enter>":
					var err error

					// Update Currency
					if currencyWidget.SelectedRow < len(currencyWidget.Rows) {
						row := currencyWidget.Rows[currencyWidget.SelectedRow]

						// Get currency and rate
						currency = fmt.Sprintf("%s %s", row[0], row[1])
						currencyVal, err = strconv.ParseFloat(row[3], 64)
						if err != nil {
							currencyVal = 0
							currency = "USD $"
						}

						// Update currency fields
						favHeader[1] = fmt.Sprintf("Price (%s)", currency)
					}
					selectCurrency = false

				case "<Escape>":
					selectCurrency = false
				}
				if selectCurrency {
					ui.Render(currencyWidget)
				}
			} else {
				if selectedTable == myPage.FavouritesTable {
					myPage.FavouritesTable.ShowCursor = true
					switch e.ID {
					case "j", "<Down>":
						myPage.FavouritesTable.ScrollDown()
					case "k", "<Up>":
						myPage.FavouritesTable.ScrollUp()
					case "<C-d>":
						myPage.FavouritesTable.ScrollHalfPageDown()
					case "<C-u>":
						myPage.FavouritesTable.ScrollHalfPageUp()
					case "<C-f>":
						myPage.FavouritesTable.ScrollPageDown()
					case "<C-b>":
						myPage.FavouritesTable.ScrollPageUp()
					case "g":
						if previousKey == "g" {
							myPage.FavouritesTable.ScrollTop()
						}
					case "<Home>":
						myPage.FavouritesTable.ScrollTop()
					case "G", "<End>":
						myPage.FavouritesTable.ScrollBottom()

					// Sort Ascending
					case "1", "2":
						idx, _ := strconv.Atoi(e.ID)
						favSortIdx = idx - 1
						myPage.FavouritesTable.Header = append([]string{}, favHeader...)
						myPage.FavouritesTable.Header[favSortIdx] = favHeader[favSortIdx] + " " + UP_ARROW
						favSortAsc = true
						utils.SortData(myPage.FavouritesTable.Rows, favSortIdx, favSortAsc, "FAVOURITES")

					// Sort Descending
					case "<F1>", "<F2>":
						myPage.FavouritesTable.Header = append([]string{}, favHeader...)
						idx, _ := strconv.Atoi(e.ID[2:3])
						favSortIdx = idx - 1
						myPage.FavouritesTable.Header[favSortIdx] = favHeader[favSortIdx] + " " + DOWN_ARROW
						favSortAsc = false
						utils.SortData(myPage.FavouritesTable.Rows, favSortIdx, favSortAsc, "FAVOURITES")

					}
				} else {
					myPage.IntervalTable.ShowCursor = true

					switch e.ID {
					case "j", "<Down>":
						myPage.IntervalTable.ScrollDown()
					case "k", "<Up>":
						myPage.IntervalTable.ScrollUp()
					case "<C-d>":
						myPage.IntervalTable.ScrollHalfPageDown()
					case "<C-u>":
						myPage.IntervalTable.ScrollHalfPageUp()
					case "<C-f>":
						myPage.IntervalTable.ScrollPageDown()
					case "<C-b>":
						myPage.IntervalTable.ScrollPageUp()
					case "g":
						if previousKey == "g" {
							myPage.IntervalTable.ScrollTop()
						}
					case "<Home>":
						myPage.IntervalTable.ScrollTop()
					case "G", "<End>":
						myPage.IntervalTable.ScrollBottom()
					case "<Enter>":
						// Update interval
						if myPage.IntervalTable.SelectedRow < len(myPage.IntervalTable.Rows) {
							// Get interval
							row := myPage.IntervalTable.Rows[myPage.IntervalTable.SelectedRow]
							val := row[0]
							myPage.ValueGraph.Data["Value"] = []float64{}

							// Send new interval for GetCoinHistory
							intervalChannel <- intervals[val]
						}
					}
				}

				ui.Render(myPage.Grid)
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		case data := <-priceChannel:
			// Update live price
			p, _ := strconv.ParseFloat(data, 64)
			myPage.PriceBox.Rows[0][0] = fmt.Sprintf("%.2f %s", p/currencyVal, currency)
			if !selectCurrency && !helpSelected && !portfolioSelected {
				ui.Render(myPage.PriceBox)
			}

		case data := <-dataChannel:
			switch data.Type {

			case "FAVOURITES":
				// Update favorites table
				rows := [][]string{}
				for symbol, price := range data.Favourites {
					p := fmt.Sprintf("%.2f", price/currencyVal)
					rows = append(rows, []string{symbol, p})
				}
				myPage.FavouritesTable.Header[1] = fmt.Sprintf("Price (%s)", currency)
				myPage.FavouritesTable.Rows = rows

			case "HISTORY":
				// Update History graph
				price := data.PriceHistory

				// Set value, min 7 max price
				myPage.ValueGraph.Data["Value"] = price
				myPage.ValueGraph.Labels["Max"] = fmt.Sprintf("%.2f %s", data.MaxPrice/currencyVal, currency)
				myPage.ValueGraph.Labels["Min"] = fmt.Sprintf("%.2f %s", data.MinPrice/currencyVal, currency)

			case "ASSET":
				// Update Details table
				myPage.DetailsTable.Header = []string{"Name", data.CoinAssetData.Data.Name}

				// Get Market Cap
				mCapStr := ""
				mCap, err := strconv.ParseFloat(data.CoinAssetData.Data.MarketCapUsd, 64)
				if err == nil {
					mCapVals, units := utils.RoundValues(mCap/currencyVal, 0)
					mCapStr = fmt.Sprintf("%.2f %s %s", mCapVals[0], units, currency)
				}

				// Get Volume Weighted Average price
				vwapStr := ""
				vwap, err := strconv.ParseFloat(data.CoinAssetData.Data.Vwap24Hr, 64)
				if err == nil {
					vwapStr = fmt.Sprintf("%.2f %s", vwap/currencyVal, currency)
				}

				// Get 24 Hr Volume
				vol, volErr := strconv.ParseFloat(data.CoinAssetData.Data.VolumeUsd24Hr, 64)
				volStr := ""
				if volErr == nil {
					vals, units := utils.RoundValues(vol/currencyVal, 0)
					volStr = fmt.Sprintf("%.2f %s %s", vals[0], units, currency)
				}

				// Aggregate data
				rows := [][]string{
					{"Symbol", data.CoinAssetData.Data.Symbol},
					{"Rank", data.CoinAssetData.Data.Rank},
					{"Market Cap", mCapStr},
					{"VWAP 24Hr", vwapStr},
					{"Explorer", data.CoinAssetData.Data.Explorer},
					{"Volume Used 24 Hr", volStr},
				}

				// Update value label in history graph
				p, err := strconv.ParseFloat(data.CoinAssetData.Data.PriceUsd, 64)
				if err == nil {
					myPage.ValueGraph.Labels["Value"] = fmt.Sprintf("%.2f %s", p/currencyVal, currency)
				}

				myPage.DetailsTable.Rows = rows

				// Update Volume Guage
				if volErr == nil {
					if mCap > 0 {
						percent := int((vol / mCap) * 100)
						if percent <= 100 && percent >= 0 {
							myPage.VolumeGauge.Percent = percent
						}
					}
				}

				// Get supply and Max supply
				supply, err1 := strconv.ParseFloat(data.CoinAssetData.Data.Supply, 64)
				maxSupply, err2 := strconv.ParseFloat(data.CoinAssetData.Data.MaxSupply, 64)

				if err1 == nil && err2 == nil {
					supplyVals, units := utils.RoundValues(supply, maxSupply)
					myPage.SupplyChart.Data = supplyVals
					myPage.SupplyChart.Title = fmt.Sprintf(" Supply (%s) ", units)
				}

				// Update Price Box Change %
				change := "NA"
				c, err := strconv.ParseFloat(data.CoinAssetData.Data.ChangePercent24Hr, 64)
				if err == nil {
					if c < 0 {
						change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -1*c)
					} else {
						change = fmt.Sprintf("%s %.2f", UP_ARROW, c)
					}
				}
				myPage.PriceBox.Rows[0][1] = change
			}

			// Sort favourites table
			if favSortIdx != -1 {
				utils.SortData(myPage.FavouritesTable.Rows, favSortIdx, favSortAsc, "FAVOURITES")

				if favSortAsc {
					myPage.FavouritesTable.Header[favSortIdx] = favHeader[favSortIdx] + " " + UP_ARROW
				} else {
					myPage.FavouritesTable.Header[favSortIdx] = favHeader[favSortIdx] + " " + DOWN_ARROW
				}
			} else {
				utils.SortData(myPage.FavouritesTable.Rows, 0, true, "FAVOURITES")
			}

		case <-tick: // Refresh UI
			updateUI()
		}
	}
}
