package coin

import (
	"context"
	"fmt"
	"strconv"
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

func DisplayCoin(ctx context.Context, id string, interval *string, dataChannel chan api.CoinData, priceChannel chan string, uiEvents <-chan ui.Event) error {
	defer ui.Clear()

	myPage := NewCoinPage()

	currency := "USD $"
	currencyVal := 1.0
	selectCurrency := false
	currencyWidget := c.NewCurrencyPage()

	selectedTable := myPage.FavouritesTable
	previousKey := ""

	updateUI := func() {
		// Get Terminal Dimensions adn clear the UI
		w, h := ui.TerminalDimensions()

		// Adjust Suuply chart Bar graph values
		myPage.SupplyChart.BarGap = ((w / 3) - (2 * myPage.SupplyChart.BarWidth)) / 2

		myPage.Grid.SetRect(0, 0, w, h)

		ui.Clear()
		if selectCurrency {
			currencyWidget.Resize(w, h)
			ui.Render(currencyWidget)
		} else {
			ui.Render(myPage.Grid)
		}
	}

	updateUI()

	t := time.NewTicker(time.Duration(1) * time.Second)
	tick := t.C

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents:
			switch e.ID {
			case "<Escape>":
				if !selectCurrency {
					selectedTable = nil
					selectCurrency = false
					return fmt.Errorf("UI Closed")
				}
			case "q", "<C-c>":
				return fmt.Errorf("coin UI Closed")
			case "<Resize>":
				updateUI()
			case "c":
				selectedTable.ShowCursor = false
				selectCurrency = true
				selectedTable = currencyWidget.Table
				selectedTable.ShowCursor = true
				currencyWidget.UpdateRows()
				updateUI()

			case "C":
				selectedTable.ShowCursor = false
				selectCurrency = true
				selectedTable = currencyWidget.Table
				selectedTable.ShowCursor = true
				currencyWidget.UpdateAll()
				updateUI()
			}
			if selectCurrency {
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
					if currencyWidget.SelectedRow < len(currencyWidget.Rows) {
						row := currencyWidget.Rows[currencyWidget.SelectedRow]
						currency = fmt.Sprintf("%s %s", row[0], row[1])
						currencyVal, err = strconv.ParseFloat(row[3], 64)
						if err != nil {
							currencyVal = 0
							currency = "USD $"
						}
					}

					selectedTable.ShowCursor = false
					selectedTable = myPage.FavouritesTable
					selectCurrency = false

				case "<Escape>":
					selectedTable.ShowCursor = false
					selectedTable = myPage.FavouritesTable
					selectCurrency = false
				}
				if selectCurrency {
					ui.Render(currencyWidget)
				}
			}

		case data := <-priceChannel:
			p, _ := strconv.ParseFloat(data, 64)
			myPage.PriceBox.Rows[0][0] = fmt.Sprintf("%.2f %s", p/currencyVal, currency)
			if !selectCurrency {
				ui.Render(myPage.PriceBox)
			}

		case data := <-dataChannel:
			switch data.Type {

			case "HISTORY":
				// Update History graph
				price := data.PriceHistory
				myPage.ValueGraph.Data["Value"] = price
				myPage.ValueGraph.Labels["Max"] = fmt.Sprintf("%.2f %s", utils.MaxFloat64(price...)/currencyVal, currency)
				myPage.ValueGraph.Labels["Min"] = fmt.Sprintf("%.2f %s", utils.MinFloat64(price...)/currencyVal, currency)

			case "ASSET":
				// Update Details table
				myPage.DetailsTable.Header = []string{"Name", data.CoinAssetData.Data.Name}
				rows := [][]string{
					{"Symbol", data.CoinAssetData.Data.Symbol},
					{"Rank", data.CoinAssetData.Data.Rank},
					{"Market Cap USD", data.CoinAssetData.Data.MarketCapUsd},
					{"VWAP 24Hr", data.CoinAssetData.Data.Vwap24Hr},
					{"Explorer", data.CoinAssetData.Data.Explorer},
				}

				p, err := strconv.ParseFloat(data.CoinAssetData.Data.PriceUsd, 64)
				if err == nil {
					myPage.ValueGraph.Labels["Value"] = fmt.Sprintf("%.2f %s", p/currencyVal, currency)
				}

				myPage.DetailsTable.Rows = rows

				// Update Volume Guage
				vol, err1 := strconv.ParseFloat(data.CoinAssetData.Data.VolumeUsd24Hr, 64)
				mCap, err2 := strconv.ParseFloat(data.CoinAssetData.Data.MarketCapUsd, 64)
				if err1 == nil && err2 == nil {
					percent := int((vol / mCap) * 100)
					if percent <= 100 && percent >= 0 {
						myPage.VolumeGauge.Percent = percent
					}
				}

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

		case <-tick:
			updateUI()
		}
	}
}
