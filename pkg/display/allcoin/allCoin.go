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

package allcoin

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	"github.com/Gituser143/cryptgo/pkg/display/coin"
	c "github.com/Gituser143/cryptgo/pkg/display/currency"
	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
	"golang.org/x/sync/errgroup"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

// DisplayAllCoins displays the main page with top coin prices, favourites and
// general coin asset data
func DisplayAllCoins(ctx context.Context, dataChannel chan api.AssetData, sendData *bool) error {

	// Initialise UI
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialise termui: %v", err)
	}
	defer ui.Close()

	// Variables for currency
	currency := "USD $"
	currencyVal := 1.0
	selectCurrency := false
	currencyWidget := c.NewCurrencyPage()

	// variables for sorting CoinTable
	coinSortIdx := -1
	coinSortAsc := false
	coinHeader := []string{
		"Rank",
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
		"Change %",
		"Supply / MaxSupply",
	}

	// variables for sorting FavouritesTable
	favSortIdx := -1
	favSortAsc := false
	favHeader := []string{
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
	}

	previousKey := ""

	// coinIDs tracks symbol: id
	coinIDs := make(map[string]string)

	// Initalise page and set selected table
	myPage := NewAllCoinPage()
	selectedTable := myPage.CoinTable

	// Initialise favourites and portfolio
	portfolio := utils.GetPortfolio()
	favourites := utils.GetFavourites()
	defer utils.SaveMetadata(favourites, currency, portfolio)

	// Initialise Help Menu
	help := widgets.NewHelpMenu()
	help.SelectHelpMenu("ALL")
	helpSelected := false

	// Pause function to pause sending and receiving of data
	pause := func() {
		*sendData = !(*sendData)
	}

	// UpdateUI to refresh UI
	updateUI := func() {
		// Get Terminal Dimensions
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)

		// Clear UI
		ui.Clear()

		// Render required widgets
		if helpSelected {
			help.Resize(w, h)
			ui.Render(help)
		} else if selectCurrency {
			currencyWidget.Resize(w, h)
			ui.Render(currencyWidget)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	// Render Empty UI
	updateUI()

	// Create Channel to get keyboard events
	uiEvents := ui.PollEvents()

	// Create ticker to periodically refresh UI
	t := time.NewTicker(time.Duration(1) * time.Second)
	tick := t.C

	for {
		select {
		case <-ctx.Done(): // Context cancelled, exit
			return ctx.Err()

		case e := <-uiEvents: // keyboard events
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				return fmt.Errorf("UI Closed")

			case "<Resize>":
				updateUI()

			case "p":
				pause()

			case "?":
				helpSelected = !helpSelected
				updateUI()

			case "f":
				if !helpSelected {
					selectedTable.ShowCursor = false
					selectedTable = myPage.FavouritesTable
				}

			case "F":
				if !helpSelected {
					selectedTable.ShowCursor = false
					selectedTable = myPage.CoinTable
				}

			case "c":
				if !helpSelected {
					selectedTable.ShowCursor = false
					selectCurrency = true
					selectedTable.ShowCursor = true
					currencyWidget.UpdateRows()
					updateUI()
				}

			case "C":
				if !helpSelected {
					selectedTable.ShowCursor = false
					selectCurrency = true
					selectedTable.ShowCursor = true
					currencyWidget.UpdateAll()
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
						coinHeader[2] = fmt.Sprintf("Price (%s)", currency)
						favHeader[1] = fmt.Sprintf("Price (%s)", currency)
					}

					selectCurrency = false

				case "<Escape>":
					selectCurrency = false
				}
				if selectCurrency {
					ui.Render(currencyWidget)
				}
			} else if selectedTable != nil {
				selectedTable.ShowCursor = true

				switch e.ID {
				case "j", "<Down>":
					selectedTable.ScrollDown()
				case "k", "<Up>":
					selectedTable.ScrollUp()
				case "<C-d>":
					selectedTable.ScrollHalfPageDown()
				case "<C-u>":
					selectedTable.ScrollHalfPageUp()
				case "<C-f>":
					selectedTable.ScrollPageDown()
				case "<C-b>":
					selectedTable.ScrollPageUp()
				case "g":
					if previousKey == "g" {
						selectedTable.ScrollTop()
					}
				case "<Home>":
					selectedTable.ScrollTop()
				case "G", "<End>":
					selectedTable.ScrollBottom()

				case "s": // Add coin to favourites
					id := ""
					symbol := ""

					// Get ID and symbol
					if selectedTable == myPage.CoinTable {
						if myPage.CoinTable.SelectedRow < len(myPage.CoinTable.Rows) {
							row := myPage.CoinTable.Rows[myPage.CoinTable.SelectedRow]
							symbol = row[1]
						}
					} else {
						if myPage.FavouritesTable.SelectedRow < len(myPage.FavouritesTable.Rows) {
							row := myPage.FavouritesTable.Rows[myPage.FavouritesTable.SelectedRow]
							symbol = row[0]
						}
					}
					id = coinIDs[symbol]
					favourites[id] = true

				case "S": // Remove coin from favourites
					id := ""
					symbol := ""

					// Get ID and symbol
					if selectedTable == myPage.CoinTable {
						if myPage.CoinTable.SelectedRow < len(myPage.CoinTable.Rows) {
							row := myPage.CoinTable.Rows[myPage.CoinTable.SelectedRow]
							symbol = row[1]
						}
					} else {
						if myPage.FavouritesTable.SelectedRow < len(myPage.FavouritesTable.Rows) {
							row := myPage.FavouritesTable.Rows[myPage.FavouritesTable.SelectedRow]
							symbol = row[0]
						}
					}
					id = coinIDs[symbol]

					delete(favourites, id)

				case "<Enter>": // Serve per coin details
					// pause UI and data send
					pause()

					id := ""
					symbol := ""

					// Get ID and symbol
					if selectedTable == myPage.CoinTable {
						if myPage.CoinTable.SelectedRow < len(myPage.CoinTable.Rows) {
							row := myPage.CoinTable.Rows[myPage.CoinTable.SelectedRow]
							symbol = row[1]
						}
					} else {
						if myPage.FavouritesTable.SelectedRow < len(myPage.FavouritesTable.Rows) {
							row := myPage.FavouritesTable.Rows[myPage.FavouritesTable.SelectedRow]
							symbol = row[0]
						}
					}
					id = coinIDs[symbol]

					if id != "" {
						// Create new errorgroup for coin page
						eg, coinCtx := errgroup.WithContext(ctx)
						coinDataChannel := make(chan api.CoinData)
						coinPriceChannel := make(chan string)
						intervalChannel := make(chan string)

						// Clear UI
						ui.Clear()

						// Serve Coin Price History
						eg.Go(func() error {
							err := api.GetCoinHistory(
								coinCtx,
								id,
								intervalChannel,
								coinDataChannel,
							)
							return err
						})

						// Serve Coin Asset data
						eg.Go(func() error {
							err := api.GetCoinAsset(coinCtx, id, coinDataChannel)
							return err
						})

						// Serve favourie coin prices
						eg.Go(func() error {
							err := api.GetFavouritePrices(coinCtx,
								favourites,
								coinDataChannel,
							)
							return err
						})

						// Serve Live price of coin
						eg.Go(func() error {
							err := api.GetLivePrice(coinCtx, id, coinPriceChannel)
							return err
						})

						// Serve Visuals for ccoin
						eg.Go(func() error {
							err := coin.DisplayCoin(
								coinCtx,
								id,
								intervalChannel,
								coinDataChannel,
								coinPriceChannel,
								uiEvents,
							)
							return err
						})

						if err := eg.Wait(); err != nil {
							if err.Error() != "UI Closed" {
								// Unpause
								pause()
								return err
							}
						}

					}
					// unpause data send and receive
					pause()
					updateUI()
				}

				if selectedTable == myPage.CoinTable {
					switch e.ID {
					// Sort Ascending
					case "1", "2", "3", "4":
						idx, _ := strconv.Atoi(e.ID)
						coinSortIdx = idx - 1
						myPage.CoinTable.Header = append([]string{}, coinHeader...)
						myPage.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + UP_ARROW
						coinSortAsc = true
						utils.SortData(myPage.CoinTable.Rows, coinSortIdx, coinSortAsc, "COINS")

					// Sort Descending
					case "<F1>", "<F2>", "<F3>", "<F4>":
						myPage.CoinTable.Header = append([]string{}, coinHeader...)
						idx, _ := strconv.Atoi(e.ID[2:3])
						coinSortIdx = idx - 1
						myPage.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + DOWN_ARROW
						coinSortAsc = false
						utils.SortData(myPage.CoinTable.Rows, coinSortIdx, coinSortAsc, "COINS")

					}
				} else if selectedTable == myPage.FavouritesTable {
					switch e.ID {
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
				}

				ui.Render(myPage.Grid)
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		case data := <-dataChannel:
			if data.IsTopCoinData {
				// Update Top Coin data
				for i, v := range data.TopCoinData {
					// Set title to coin name
					myPage.TopCoinGraphs[i].Title = " " + data.TopCoins[i] + " "

					// Update value graphs
					myPage.TopCoinGraphs[i].Data["Value"] = v

					// Set value, max & min values
					myPage.TopCoinGraphs[i].Labels["Value"] = fmt.Sprintf("%.2f %s", v[len(v)-1]/currencyVal, currency)
					myPage.TopCoinGraphs[i].Labels["Max"] = fmt.Sprintf("%.2f %s", utils.MaxFloat64(v...)/currencyVal, currency)
					myPage.TopCoinGraphs[i].Labels["Min"] = fmt.Sprintf("%.2f %s", utils.MinFloat64(v...)/currencyVal, currency)
				}
			} else {
				rows := [][]string{}
				favouritesData := [][]string{}

				// Update currency headers
				myPage.CoinTable.Header[2] = fmt.Sprintf("Price (%s)", currency)
				myPage.FavouritesTable.Header[1] = fmt.Sprintf("Price (%s)", currency)

				// Iterate over coin assets
				for _, val := range data.Data {
					// Get coin price
					price := "NA"
					p, err := strconv.ParseFloat(val.PriceUsd, 64)
					if err == nil {
						price = fmt.Sprintf("%.2f", p/currencyVal)
					}

					// Get change %
					change := "NA"
					c, err := strconv.ParseFloat(val.ChangePercent24Hr, 64)
					if err == nil {
						if c < 0 {
							change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -1*c)
						} else {
							change = fmt.Sprintf("%s %.2f", UP_ARROW, c)
						}
					}

					// Get supply and Max supply
					s, err1 := strconv.ParseFloat(val.Supply, 64)
					ms, err2 := strconv.ParseFloat(val.MaxSupply, 64)

					units := ""
					var supplyVals []float64
					supplyData := ""

					if err1 == nil && err2 == nil {
						supplyVals, units = utils.RoundValues(s, ms)
						supplyData = fmt.Sprintf("%.2f%s / %.2f%s", supplyVals[0], units, supplyVals[1], units)
					} else {
						if err1 != nil {
							supplyVals, units = utils.RoundValues(s, ms)
							supplyData = fmt.Sprintf("NA / %.2f%s", supplyVals[1], units)
						} else {
							supplyVals, units = utils.RoundValues(s, ms)
							supplyData = fmt.Sprintf("%.2f%s / NA", supplyVals[0], units)
						}
					}

					// Aggregate data
					rows = append(rows, []string{
						val.Rank,
						val.Symbol,
						price,
						change,
						supplyData,
					})

					// Update new coin ids
					if _, ok := coinIDs[val.Symbol]; !ok {
						coinIDs[val.Symbol] = val.Id
					}

					// Aggregate favourite data
					if _, ok := favourites[val.Id]; ok {
						favouritesData = append(favouritesData, []string{
							val.Symbol,
							price,
						})
					}
				}

				myPage.CoinTable.Rows = rows
				myPage.FavouritesTable.Rows = favouritesData

				// Sort CoinTable data
				if coinSortIdx != -1 {
					utils.SortData(myPage.CoinTable.Rows, coinSortIdx, coinSortAsc, "COINS")

					if coinSortAsc {
						myPage.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + UP_ARROW
					} else {
						myPage.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + DOWN_ARROW
					}
				}

				// Sort FavouritesTable Data
				if favSortIdx != -1 {
					utils.SortData(myPage.FavouritesTable.Rows, favSortIdx, favSortAsc, "FAVOURITES")

					if favSortAsc {
						myPage.FavouritesTable.Header[favSortIdx] = favHeader[favSortIdx] + " " + UP_ARROW
					} else {
						myPage.FavouritesTable.Header[favSortIdx] = favHeader[favSortIdx] + " " + DOWN_ARROW
					}
				}
			}

		case <-tick: // Refresh UI
			if *sendData {
				updateUI()
			}
		}

	}

}
