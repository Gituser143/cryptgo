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
	selectedTable := myPage.DetailsTable
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
	// intervals := map[string]string{
	// 	"1  min":  "m1",
	// 	"5  min":  "m5",
	// 	"15 min":  "m15",
	// 	"30 min":  "m30",
	// 	"1  hour": "h1",
	// 	"2  hour": "h2",
	// 	"6  hour": "h6",
	// 	"12 hour": "h12",
	// 	"1  day":  "d1",
	// }

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
			case "<Escape>", "q", "<C-c>":
				if !helpSelected && !selectCurrency && !portfolioSelected {
					selectCurrency = false
					ui.Clear()
					return fmt.Errorf("UI Closed")
				}

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
					selectedTable = myPage.DetailsTable
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
					help.Table.ScrollDown()
					ui.Render(help)
				case "k", "<Up>":
					help.Table.ScrollUp()
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
					myPage.DetailsTable.ShowCursor = true

					switch e.ID {
					case "j", "<Down>":
						myPage.DetailsTable.ScrollDown()
					case "k", "<Up>":
						myPage.DetailsTable.ScrollUp()
					case "<C-d>":
						myPage.DetailsTable.ScrollHalfPageDown()
					case "<C-u>":
						myPage.DetailsTable.ScrollHalfPageUp()
					case "<C-f>":
						myPage.DetailsTable.ScrollPageDown()
					case "<C-b>":
						myPage.DetailsTable.ScrollPageUp()
					case "g":
						if previousKey == "g" {
							myPage.DetailsTable.ScrollTop()
						}
					case "<Home>":
						myPage.DetailsTable.ScrollTop()
					case "G", "<End>":
						myPage.DetailsTable.ScrollBottom()
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
			myPage.PriceBox.Rows = []string{
				fmt.Sprintf("%.2f %s", p/currencyVal, currency),
			}
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

			case "DETAILS":
				// Update Details table
				myPage.DetailsTable.Header = []string{"Name", data.Details.Name}

				marketCapVals, units := utils.RoundValues(data.Details.MarketCap, 0)
				marketCap := fmt.Sprintf("%.2f %s %s", marketCapVals[0]/currencyVal, units, currency)

				ATHVals, units := utils.RoundValues(data.Details.ATH, 0)
				ATH := fmt.Sprintf("%.2f %s %s", ATHVals[0]/currencyVal, units, currency)

				ATLVals, units := utils.RoundValues(data.Details.ATL, 0)
				ATL := fmt.Sprintf("%.2f %s %s", ATLVals[0]/currencyVal, units, currency)

				TotalVolVals, units := utils.RoundValues(data.Details.TotalVolume, 0)
				TotalVolume := fmt.Sprintf("%.2f %s %s", TotalVolVals[0]/currencyVal, units, currency)

				rows := [][]string{
					{"Symbol", data.Details.Symbol},
					{"Rank", data.Details.Rank},
					{"BlockTime (min)", data.Details.BlockTime},
					{"MarketCap", marketCap},
					{"ATH", ATH},
					{"ATHDate", data.Details.ATHDate},
					{"ATL", ATL},
					{"ATLDate", data.Details.ATLDate},
					{"TotalVolume", TotalVolume},
					{"LastUpdate", data.Details.LastUpdate},
				}

				myPage.DetailsTable.Rows = rows

				// Get Change Percents
				myPage.ChangesTable.Rows = data.Details.ChangePercents

				// Get supply and Max supply
				supply := data.Details.CurrentSupply
				maxSupply := data.Details.TotalSupply

				supplyVals, units := utils.RoundValues(supply, maxSupply)
				myPage.SupplyChart.Data = supplyVals
				myPage.SupplyChart.Title = fmt.Sprintf(" Supply (%s) ", units)

				// Get Explorers
				myPage.ExplorerTable.Rows = data.Details.Explorers

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
