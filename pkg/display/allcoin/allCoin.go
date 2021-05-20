package allcoin

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	c "github.com/Gituser143/cryptgo/pkg/display/currency"
	"github.com/Gituser143/cryptgo/pkg/utils"
	ui "github.com/gizak/termui/v3"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

func DisplayAllCoins(ctx context.Context, dataChannel chan api.AssetData) error {

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialise termui: %v", err)
	}
	defer ui.Close()

	var onTop sync.Once
	var onAll sync.Once

	run := true

	currency := "USD $"
	currencyVal := 1.0
	selectCurrency := false
	currencyWidget := c.NewCurrencyPage()

	sortIdx := -1
	sortAsc := false
	header := []string{
		"Rank",
		"Symbol",
		fmt.Sprintf("Price (%s)", currency),
		"Change %",
		"Supply / MaxSupply",
	}

	previousKey := ""

	myPage := NewAllCoinPage()
	selectedTable := myPage.CoinTable

	pause := func() {
		run = !run
	}

	updateUI := func() {
		// Get Terminal Dimensions adn clear the UI
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)

		currencyWidget.Resize(w, h)
		if selectCurrency {
			ui.Clear()
			ui.Render(currencyWidget)
		} else {
			ui.Clear()
			ui.Render(myPage.Grid)
		}
	}

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(1) * time.Second)
	tick := t.C

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": // q or Ctrl-C to quit
				return fmt.Errorf("UI Closed")

			case "<Resize>":
				updateUI()

			case "s":
				pause()

			case "f":
				selectedTable.ShowCursor = false
				selectedTable = myPage.FavouritesTable

			case "F":
				selectedTable.ShowCursor = false
				selectedTable = myPage.CoinTable

			case "c":
				selectCurrency = true
				selectedTable = currencyWidget.Table
				selectedTable.ShowCursor = true
				currencyWidget.UpdateRows()
				updateUI()
			}
			if selectCurrency {
				switch e.ID {
				case "<Enter>":
					var err error
					if currencyWidget.SelectedRow < len(currencyWidget.Rows) {
						row := currencyWidget.Rows[currencyWidget.SelectedRow]
						currency = fmt.Sprintf("%s %s", row[0], row[1])
						currencyVal, err = strconv.ParseFloat(row[2], 64)
						if err != nil {
							currencyVal = 0
							currency = "USD $"
						}
						header[2] = fmt.Sprintf("Price (%s)", currency)
					}

					selectedTable = myPage.CoinTable
					selectCurrency = false
					updateUI()

				case "<Escape>":
					selectedTable = myPage.CoinTable
					selectCurrency = false
					updateUI()
				}
			}
			if selectedTable != nil {
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

					// Sort Ascending
				case "1", "2", "3", "4":
					idx, _ := strconv.Atoi(e.ID)
					sortIdx = idx - 1
					myPage.CoinTable.Header = append([]string{}, header...)
					myPage.CoinTable.Header[sortIdx] = header[sortIdx] + " " + UP_ARROW
					sortAsc = true
					utils.SortData(myPage.CoinTable.Rows, sortIdx, sortAsc, "COINS")

				// Sort Descending
				case "<F1>", "<F2>", "<F3>", "<F4>":
					myPage.CoinTable.Header = append([]string{}, header...)
					idx, _ := strconv.Atoi(e.ID[2:3])
					sortIdx = idx - 1
					myPage.CoinTable.Header[sortIdx] = header[sortIdx] + " " + DOWN_ARROW
					sortAsc = false
					utils.SortData(myPage.CoinTable.Rows, sortIdx, sortAsc, "COINS")
				}

				updateUI()
				if previousKey == "g" {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		case data := <-dataChannel:
			if data.IsTopCoinData {
				for i, v := range data.TopCoinData {
					myPage.TopCoinGraphs[i].Title = " " + data.TopCoins[i] + " "
					myPage.TopCoinGraphs[i].Data["Value"] = v
					myPage.TopCoinGraphs[i].Labels["Value"] = fmt.Sprintf("%.2f", v[len(v)-1])
					myPage.TopCoinGraphs[i].Labels["Max"] = fmt.Sprintf("%.2f", utils.MaxFloat64(v...))
					myPage.TopCoinGraphs[i].Labels["Min"] = fmt.Sprintf("%.2f", utils.MinFloat64(v...))
				}
				onTop.Do(updateUI)
			} else {
				rows := [][]string{}
				for _, val := range data.Data {
					price := "NA"
					p, err := strconv.ParseFloat(val.PriceUsd, 64)
					if err == nil {
						price = fmt.Sprintf("%.2f", p/currencyVal)
						myPage.CoinTable.Header[2] = fmt.Sprintf("Price (%s)", currency)
					}

					change := "NA"
					c, err := strconv.ParseFloat(val.ChangePercent24Hr, 64)
					if err == nil {
						if c < 0 {
							change = fmt.Sprintf("%s %.2f", DOWN_ARROW, -1*c)
						} else {
							change = fmt.Sprintf("%s %.2f", UP_ARROW, c)
						}
					}

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

					rows = append(rows, []string{
						val.Rank,
						val.Symbol,
						price,
						change,
						supplyData,
					})
				}
				myPage.CoinTable.Rows = rows

				if sortIdx != -1 {
					utils.SortData(myPage.CoinTable.Rows, sortIdx, sortAsc, "COINS")
					if sortAsc {
						myPage.CoinTable.Header[sortIdx] = header[sortIdx] + " " + UP_ARROW
					} else {
						myPage.CoinTable.Header[sortIdx] = header[sortIdx] + " " + DOWN_ARROW
					}
				}

				onAll.Do(updateUI)
			}

		case <-tick:
			updateUI()
		}

	}

}
