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
	"context"
	"fmt"
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

// DisplayPortfolio serves the prtfolio page
func DisplayPortfolio(ctx context.Context, dataChannel chan api.AssetData, sendData *bool) error {

	// Initialise UI
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialise termui: %v", err)
	}
	defer ui.Close()

	// Initialise page
	page := newPortfolioPage()
	selectedTable := page.CoinTable
	utilitySelected := ""

	// Variables for CoinIDs
	coinIDMap := api.NewCoinIDMap()
	coinIDMap.Populate()

	// currency variables
	currencyWidget := uw.NewCurrencyPage()
	currencyID := utils.GetCurrencyID()
	currencyID, currency, currencyVal := currencyWidget.Get(currencyID)

	// get portfolio details
	portfolioMap := utils.GetPortfolio()

	// get performers map
	performersMap := getEmptyPerformers()

	// get favourites
	favourites := utils.GetFavourites()

	// Save metadata back to disk
	defer func() {
		utils.SaveMetadata(favourites, currencyID, portfolioMap)
	}()

	// Initialise help menu
	help := widgets.NewHelpMenu()
	help.SelectHelpMenu("PORTFOLIO")

	// Variables for sorting CoinTable
	coinSortIdx := -1
	coinSortAsc := false
	coinHeader := []string{
		"Rank",
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
		"Change % (1d)",
		"Holding",
		fmt.Sprintf("Balance (%s)", currency),
		"Holding %",
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
		page.Grid.SetRect(0, 0, w, h)

		// Clear UI
		ui.Clear()

		// Render required widgets
		switch utilitySelected {
		case "HELP":
			help.Resize(w, h)
			ui.Render(help)
		case "CURRENCY":
			currencyWidget.Resize(w, h)
			ui.Render(currencyWidget)
		default:
			ui.Render(page.Grid)
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

		case e := <-uiEvents:
			switch e.ID {

			// handle button events
			case "q", "<C-c>":
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

			case "c":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = currencyWidget.Table
					selectedTable.ShowCursor = true
					currencyWidget.UpdateRows(false)
					utilitySelected = "CURRENCY"
				}

			case "C":
				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = currencyWidget.Table
					selectedTable.ShowCursor = true
					currencyWidget.UpdateRows(true)
					utilitySelected = "CURRENCY"
				}

			case "e":
				switch utilitySelected {
				case "":
					id := ""
					symbol := ""

					// Get ID and symbol
					if selectedTable == page.CoinTable {
						if page.CoinTable.SelectedRow < len(page.CoinTable.Rows) {
							row := page.CoinTable.Rows[page.CoinTable.SelectedRow]
							symbol = row[1]
						}
					}

					coinIDs := coinIDMap[symbol]

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
					// Update Currency
					if currencyWidget.SelectedRow < len(currencyWidget.Rows) {
						row := currencyWidget.Rows[currencyWidget.SelectedRow]

						// Get currency and rate
						currencyID = row[0]
						currencyID, currency, currencyVal = currencyWidget.Get(currencyID)

						// Update currency fields
						coinHeader[2] = fmt.Sprintf("Price (%s)", currency)
						coinHeader[5] = fmt.Sprintf("Balance (%s)", currency)
					}
					utilitySelected = ""

				case "":

					// pause UI and data send
					pause()

					symbol := ""

					// Get ID and symbol
					if selectedTable == page.CoinTable {
						if page.CoinTable.SelectedRow < len(page.CoinTable.Rows) {
							row := page.CoinTable.Rows[page.CoinTable.SelectedRow]
							symbol = row[1]
						}
					}

					coinIDs := coinIDMap[symbol]

					coinCapID := coinIDs.CoinCapID
					coinGeckoID := coinIDs.CoinGeckoID

					if coinGeckoID != "" {
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
								coinGeckoID,
								intervalChannel,
								coinDataChannel,
							)
							return err
						})

						// Serve Coin Asset data
						eg.Go(func() error {
							err := api.GetCoinDetails(coinCtx, coinGeckoID, coinDataChannel)
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
						if coinCapID != "" {
							eg.Go(func() error {
								api.GetLivePrice(coinCtx, coinCapID, coinPriceChannel)
								// Send NA to indicate price is not being updated
								go func() {
									coinPriceChannel <- "NA"
								}()
								return nil
							})
						}

						utils.SaveMetadata(favourites, currencyID, portfolioMap)

						// Serve Visuals for coin
						eg.Go(func() error {
							err := coin.DisplayCoin(
								coinCtx,
								coinGeckoID,
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

						currencyID = utils.GetCurrencyID()
						currencyID, currency, currencyVal = currencyWidget.Get(currencyID)

					}

					// unpause data send and receive
					pause()
					updateUI()
					utilitySelected = ""
				}

				if utilitySelected == "" {
					selectedTable.ShowCursor = false
					selectedTable = page.CoinTable
					selectedTable.ShowCursor = true
				}

			// Handle Navigations
			case "<Escape>":
				utilitySelected = ""
				selectedTable = page.CoinTable
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

			// handle sorting
			case "1", "2", "3", "4", "5", "6", "7":
				// Sort Ascending
				if utilitySelected == "" {
					idx, _ := strconv.Atoi(e.ID)
					coinSortIdx = idx - 1
					page.CoinTable.Header = append([]string{}, coinHeader...)
					page.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + UP_ARROW
					coinSortAsc = true
					utils.SortData(page.CoinTable.Rows, coinSortIdx, coinSortAsc, "PORTFOLIO")
				}

			case "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>":
				// Sort Descending
				if utilitySelected == "" {
					page.CoinTable.Header = append([]string{}, coinHeader...)
					idx, _ := strconv.Atoi(e.ID[2:3])
					coinSortIdx = idx - 1
					page.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + DOWN_ARROW
					coinSortAsc = false
					utils.SortData(page.CoinTable.Rows, coinSortIdx, coinSortAsc, "PORTFOLIO")
				}
			}

			updateUI()
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		case data := <-dataChannel:
			rows := [][]string{}

			// Update currency headers
			page.CoinTable.Header[2] = fmt.Sprintf("Price (%s)", currency)
			page.CoinTable.Header[5] = fmt.Sprintf("Balance (%s)", currency)

			// variables to calculate holding %
			balanceMap := map[string]float64{}
			portfolioTotal := 0.0
			durations := []string{"1h", "24h", "7d", "30d", "1y"}

			// Iterate over coin assets
			for _, val := range data.AllCoinData {
				// Get coins in portfolio
				if portfolioHolding, ok := portfolioMap[val.ID]; ok {
					// Get coin details
					price := fmt.Sprintf("%.2f", val.CurrentPrice/currencyVal)

					change := "NA"
					percentageChange := api.GetPercentageChangeForDuration(val, "24h")
					if percentageChange < 0 {
						change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -percentageChange)
					} else {
						change = fmt.Sprintf("%s %.2f", UP_ARROW, percentageChange)
					}

					rank := fmt.Sprintf("%d", val.MarketCapRank)
					symbol := strings.ToUpper(val.Symbol)
					holding := fmt.Sprintf("%.5f", portfolioHolding)
					balanceFloat := val.CurrentPrice / currencyVal * portfolioHolding
					balance := fmt.Sprintf("%.2f", balanceFloat)

					// Aggregate data
					rows = append(rows, []string{
						rank,
						symbol,
						price,
						change,
						holding,
						balance,
						"holdingPercent", // calculated after total balance is calculated
					})

					// Calculate portfolio total
					portfolioTotal += balanceFloat

					// Keep track of a coin's balance
					balanceMap[symbol] = balanceFloat

					// Calculate best and worst performers
					for _, duration := range durations {
						val := api.GetPercentageChangeForDuration(val, duration)

						if val > performersMap[duration].BestVal {
							performersMap[duration] = performer{
								BestVal:   val,
								BestCoin:  symbol,
								WorstVal:  performersMap[duration].WorstVal,
								WorstCoin: performersMap[duration].WorstCoin,
							}
						}

						if val < performersMap[duration].WorstVal {
							performersMap[duration] = performer{
								BestVal:   performersMap[duration].BestVal,
								BestCoin:  performersMap[duration].BestCoin,
								WorstVal:  val,
								WorstCoin: symbol,
							}
						}
					}
				}
			}

			// Update portfolio holding % values
			for i, row := range rows {
				symbol := row[1]
				rows[i][6] = fmt.Sprintf("%.2f", (balanceMap[symbol]/portfolioTotal)*100)
			}

			// Update coin table
			page.CoinTable.Rows = rows

			// Update details table
			page.DetailsTable.Header = []string{
				"Balance",
				fmt.Sprintf("%.2f", portfolioTotal),
			}
			page.DetailsTable.Rows = [][]string{
				{"Currency", currency},
				{"Coins", fmt.Sprintf("%d", len(portfolioMap))},
			}

			// Update Best Performers Table
			BestPerformerRows := [][]string{}
			WorstPerformerRows := [][]string{}

			// Format best and worst performer data
			for _, duration := range durations {
				change := ""
				if performersMap[duration].BestVal < 0 {
					change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -performersMap[duration].BestVal)
				} else {
					change = fmt.Sprintf("%s %.2f", UP_ARROW, performersMap[duration].BestVal)
				}

				BestPerformerRows = append(BestPerformerRows, []string{
					duration,
					performersMap[duration].BestCoin,
					change,
				})

				if performersMap[duration].WorstVal < 0 {
					change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -performersMap[duration].WorstVal)
				} else {
					change = fmt.Sprintf("%s %.2f", UP_ARROW, performersMap[duration].WorstVal)
				}

				WorstPerformerRows = append(WorstPerformerRows, []string{
					duration,
					performersMap[duration].WorstCoin,
					change,
				})

			}

			page.BestPerformerTable.Rows = BestPerformerRows
			page.WorstPerformerTable.Rows = WorstPerformerRows

			// Sort CoinTable data
			if coinSortIdx != -1 {
				utils.SortData(page.CoinTable.Rows, coinSortIdx, coinSortAsc, "PORTFOLIO")

				if coinSortAsc {
					page.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + UP_ARROW
				} else {
					page.CoinTable.Header[coinSortIdx] = coinHeader[coinSortIdx] + " " + DOWN_ARROW
				}
			}

		case <-tick: // Refresh UI
			updateUI()
		}
	}

}
