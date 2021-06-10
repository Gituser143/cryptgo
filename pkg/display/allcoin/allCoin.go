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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	"github.com/Gituser143/cryptgo/pkg/display/coin"
	uw "github.com/Gituser143/cryptgo/pkg/display/utilitywidgets"
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
func DisplayAllCoins(ctx context.Context, dataChannel chan api.AssetData, searchChannel chan []api.CoinSearchDetails, sendData *bool) error {

	// Initialise UI
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialise termui: %v", err)
	}
	defer ui.Close()

	// Variables for CoinIDs
	coinIDMap, m := api.NewCoinIDMap()
	go func() {
		m.Lock()
		coinIDMap.Populate()
		m.Unlock()
	}()

	// Variables for currency
	currency := "USD $"
	currencyVal := 1.0
	currencyWidget := uw.NewCurrencyPage()

	// Variables for percentage change
	changePercent := "24h"
	changePercentWidget := uw.NewChangePercentPage()

	// Initalise page and set selected table
	myPage := NewAllCoinPage()
	selectedTable := myPage.CoinTable
	utilitySelected := ""

	// Initialise favourites and portfolio
	portfolioMap := utils.GetPortfolio()
	favourites := utils.GetFavourites()
	defer utils.SaveMetadata(favourites, currency, portfolioMap)

	// Initialise Help Menu
	help := widgets.NewHelpMenu()
	help.SelectHelpMenu("ALL")

	// Initialise Search Menu
	search := widgets.NewSearchMenu()
	search.Reset()

	// Initiliase Portfolio Table
	portfolioTable := uw.NewPortfolioPage()

	// Variables for sorting CoinTable
	coinSortIdx := -1
	coinSortAsc := false
	coinHeader := []string{
		"Rank",
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
		fmt.Sprintf("Change %%(%s)", changePercent),
		"Supply / MaxSupply",
	}

	// Variables for sorting FavouritesTable
	favSortIdx := -1
	favSortAsc := false
	favHeader := []string{
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
	}

	previousKey := ""

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
		switch utilitySelected {
		case "HELP":
			help.Resize(w, h)
			ui.Render(help)
		case "PORTFOLIO":
			portfolioTable.Resize(w, h)
			ui.Render(portfolioTable)
		case "CURRENCY":
			currencyWidget.Resize(w, h)
			ui.Render(currencyWidget)
		case "CHANGE":
			changePercentWidget.Resize(w, h)
			ui.Render(changePercentWidget)
		case "SEARCH":
			search.Resize(w, h)
			ui.Render(search)
		default:
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
			// Handle Utility Selection, resize and Quit
			switch e.ID {
			case "q", "<C-c>":
				if utilitySelected == "SEARCH" && search.Table.SelectedRow == 0 && e.ID == "q" {
					break
				}
				return fmt.Errorf("UI Closed")

			case "<Resize>":
				updateUI()

			case "p":
				pause()

			case "?":
				selectedTable.ShowCursor = false
				selectedTable = help.Table
				selectedTable.ShowCursor = true
				utilitySelected = "HELP"
				updateUI()

			case "f":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = myPage.FavouritesTable
					selectedTable.ShowCursor = true
				}

			case "F":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = myPage.CoinTable
					selectedTable.ShowCursor = true
				}

			case "c":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = currencyWidget.Table
					selectedTable.ShowCursor = true
					currencyWidget.UpdateRows()
					utilitySelected = "CURRENCY"
				}

			case "C":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = currencyWidget.Table
					selectedTable.ShowCursor = true
					currencyWidget.UpdateAll()
					utilitySelected = "CURRENCY"
				}

			case "%":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = changePercentWidget.Table
					selectedTable.ShowCursor = true
					utilitySelected = "CHANGE"
				}

			case "P":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = portfolioTable.Table
					selectedTable.ShowCursor = true
					portfolioTable.UpdateRows(portfolioMap, currency, currencyVal)
					utilitySelected = "PORTFOLIO"
				}

			// Handle Navigations
			case "<Escape>":
				utilitySelected = ""
				search.Reset()
				selectedTable = myPage.CoinTable
				selectedTable.ShowCursor = true
				updateUI()
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

			// Handle Actions
			case "e":
				switch utilitySelected {
				case "PORTFOLIO":
					id := ""
					symbol := ""

					// Get ID and symbol
					if portfolioTable.SelectedRow < len(portfolioTable.Rows) {
						row := portfolioTable.Rows[portfolioTable.SelectedRow]
						symbol = row[1]
					}

					m.Lock()
					coinIDs := coinIDMap[symbol]
					m.Unlock()

					id = coinIDs.CoinGeckoID

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
				case "":
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

					m.Lock()
					coinIDs := coinIDMap[symbol]
					m.Unlock()

					id = coinIDs.CoinGeckoID

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
				}

			case "<Enter>":
				switch utilitySelected {
				case "CURRENCY":
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
					utilitySelected = ""

				case "CHANGE":
					if changePercentWidget.SelectedRow < len(changePercentWidget.Rows) {
						row := changePercentWidget.Rows[changePercentWidget.SelectedRow]

						changePercent = uw.DurationMap[row[0]]

						coinHeader[3] = fmt.Sprintf("Change %%(%s)", changePercent)
					}
					utilitySelected = ""

				case "":
					// pause UI and data send
					pause()

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
					m.Lock()
					coinIDs := coinIDMap[symbol]
					m.Unlock()

					coinCapId := coinIDs.CoinCapID
					coinGeckoId := coinIDs.CoinGeckoID

					if coinGeckoId != "" {
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
								coinGeckoId,
								intervalChannel,
								coinDataChannel,
							)
							return err
						})

						// Serve Coin Asset data
						eg.Go(func() error {
							err := api.GetCoinDetails(coinCtx, coinGeckoId, coinDataChannel)
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
						if coinCapId != "" {
							eg.Go(func() error {
								api.GetLivePrice(coinCtx, coinCapId, coinPriceChannel)
								// Send NA to indicate price is not being updated
								go func() {
									coinPriceChannel <- "NA"
								}()
								return nil
							})
						}

						// Serve Visuals for coin
						eg.Go(func() error {
							err := coin.DisplayCoin(
								coinCtx,
								coinGeckoId,
								coinIDMap,
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
					utilitySelected = ""

				case "SEARCH":
					// if currently searching
					if search.Table.SelectedRow == 0 {
						searchList := []string{}
						for symbol, coinID := range coinIDMap {
							// check if search query is part of the symbol,
							//       if coinID exists for this search query
							ok, _ := regexp.MatchString(fmt.Sprintf(".*%s.*", search.SearchString), symbol)
							if ok && coinID.CoinGeckoID != "" && coinID.CoinCapID != "" {
								searchList = append(searchList, symbol)
							}
						}
						if len(searchList) > 0 {
							// Get Searched Prices
							eg, coinCtx := errgroup.WithContext(ctx)

							eg.Go(func() error {
								err := api.GetSearchedPrices(
									coinCtx,
									searchList,
									coinIDMap,
									searchChannel,
								)
								return err
							})
						} else {
							search.SymbolDoesNotExist = true
							search.SearchData = [][]string{}
						}
					} else {
						// else user has chosen a symbol
						symbol := search.SelectedItem
						m.Lock()
						coinIDs := coinIDMap[symbol]
						m.Unlock()

						coinCapId := coinIDs.CoinCapID
						coinGeckoId := coinIDs.CoinGeckoID

						if coinGeckoId != "" {
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
									coinGeckoId,
									intervalChannel,
									coinDataChannel,
								)
								return err
							})

							// Serve Coin Asset data
							eg.Go(func() error {
								err := api.GetCoinDetails(coinCtx, coinGeckoId, coinDataChannel)
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
							if coinCapId != "" {
								eg.Go(func() error {
									api.GetLivePrice(coinCtx, coinCapId, coinPriceChannel)
									// Send NA to indicate price is not being updated
									go func() {
										coinPriceChannel <- "NA"
									}()
									return nil
								})
							}

							// Serve Visuals for coin
							eg.Go(func() error {
								err := coin.DisplayCoin(
									coinCtx,
									coinGeckoId,
									coinIDMap,
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
						utilitySelected = ""

					}

				}

				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = myPage.CoinTable
					selectedTable.ShowCursor = true
				}

			case "s":
				if utilitySelected == "" {
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
					coinIDs := coinIDMap[symbol]
					id = coinIDs.CoinGeckoID
					favourites[id] = true
				}

			case "S":
				if utilitySelected == "" {
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

					m.Lock()
					coinIDs := coinIDMap[symbol]
					m.Unlock()

					id = coinIDs.CoinGeckoID

					delete(favourites, id)
				}

			case "<C-s>":
				selectedTable.ShowCursor = false
				selectedTable = search.Table
				search.Reset()
				selectedTable.ShowCursor = true
				utilitySelected = "SEARCH"
				updateUI()
			}

			if utilitySelected == "" {
				// Handle Sorting of tables
				switch selectedTable {
				case myPage.CoinTable:
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

				case myPage.FavouritesTable:
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
			}

			if utilitySelected == "SEARCH" {
				if e.ID == "<Backspace>" && search.Table.SelectedRow == 0 {
					searchStrLen := len(search.SearchString)
					if searchStrLen > 0 {
						search.SearchString = search.SearchString[:searchStrLen-1]
					}
					search.SearchData = [][]string{}
					search.SymbolDoesNotExist = false // reset red color on change in search string
					search.Header = []string{}
				} else if len(e.ID) == 1 && len(search.SearchString) < 10 && search.Table.SelectedRow == 0 {
					// check if alphabet
					//       if length is less than max length
					//       if on search row
					search.SearchString += strings.ToUpper(e.ID)
					search.SymbolDoesNotExist = false
					search.SearchData = [][]string{}
					search.Header = []string{}
				}
			}

			updateUI()
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-searchChannel:
			searchData := [][]string{}
			for _, val := range data {
				searchData = append(searchData, []string{val.Symbol, val.Name, fmt.Sprintf("%v", val.Price)})
			}
			search.SearchData = searchData

		case data := <-dataChannel:
			if data.IsTopCoinData {
				// Update Top Coin data
				for i, v := range data.TopCoinData {
					// Set title to coin name
					myPage.TopCoinGraphs[i].Title = fmt.Sprintf(" %s (7D) ", data.TopCoins[i])

					// Update value graphs
					myPage.TopCoinGraphs[i].Data["Value"] = v

					// Set value, max & min values
					maxValue := data.MaxPrices[i] / currencyVal
					minValue := data.MinPrices[i] / currencyVal
					// Current value is last point (cleaned) in graph + minimum value
					value := (v[len(v)-1] + data.MinPrices[i]) / currencyVal

					myPage.TopCoinGraphs[i].Labels["Value"] = fmt.Sprintf("%.2f %s", value, currency)
					myPage.TopCoinGraphs[i].Labels["Max"] = fmt.Sprintf("%.2f %s", maxValue, currency)
					myPage.TopCoinGraphs[i].Labels["Min"] = fmt.Sprintf("%.2f %s", minValue, currency)
				}
			} else {
				rows := [][]string{}
				favouritesData := [][]string{}

				// Update currency headers
				myPage.CoinTable.Header[2] = fmt.Sprintf("Price (%s)", currency)
				myPage.CoinTable.Header[3] = fmt.Sprintf("Change %%(%s)", changePercent)
				myPage.FavouritesTable.Header[1] = fmt.Sprintf("Price (%s)", currency)

				// Iterate over coin assets
				for _, val := range data.AllCoinData {
					// Get coin price
					price := fmt.Sprintf("%.2f", val.CurrentPrice/currencyVal)

					// Get change %
					change := "NA"
					percentageChange := api.GetPercentageChangeForDuration(val, changePercent)
					if percentageChange < 0 {
						change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -1*percentageChange)
					} else {
						change = fmt.Sprintf("%s %.2f", UP_ARROW, percentageChange)
					}

					units := ""
					var supplyVals []float64
					supplyData := ""

					if val.CirculatingSupply != 0.00 && val.TotalSupply != 0.00 {
						supplyVals, units = utils.RoundValues(val.CirculatingSupply, val.TotalSupply)
						supplyData = fmt.Sprintf("%.2f%s / %.2f%s", supplyVals[0], units, supplyVals[1], units)
					} else {
						if val.CirculatingSupply == 0.00 {
							supplyVals, units = utils.RoundValues(val.CirculatingSupply, val.TotalSupply)
							supplyData = fmt.Sprintf("NA / %.2f%s", supplyVals[1], units)
						} else {
							supplyVals, units = utils.RoundValues(val.CirculatingSupply, val.TotalSupply)
							supplyData = fmt.Sprintf("%.2f%s / NA", supplyVals[0], units)
						}
					}

					rank := fmt.Sprintf("%d", val.MarketCapRank)

					// Aggregate data
					rows = append(rows, []string{
						rank,
						strings.ToUpper(val.Symbol),
						price,
						change,
						supplyData,
					})

					// Aggregate favourite data
					if _, ok := favourites[val.ID]; ok {
						favouritesData = append(favouritesData, []string{
							strings.ToUpper(val.Symbol),
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
