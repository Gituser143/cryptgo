package coin

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	"github.com/Gituser143/cryptgo/pkg/utils"
	ui "github.com/gizak/termui/v3"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

func DisplayCoin(ctx context.Context, id string, interval *string, dataChannel chan api.CoinData, priceChannel chan string) error {
	defer ui.Clear()

	myPage := NewCoinPage()

	updateUI := func() {
		// Get Terminal Dimensions adn clear the UI
		w, h := ui.TerminalDimensions()

		// Adjust Suuply chart Bar graph values
		myPage.SupplyChart.BarGap = ((w / 3) - (2 * myPage.SupplyChart.BarWidth)) / 2

		myPage.Grid.SetRect(0, 0, w, h)

		ui.Render(myPage.Grid)
	}

	updateUI()

	uiEvents := ui.PollEvents()
	t := time.NewTicker(time.Duration(1) * time.Second)
	tick := t.C

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-uiEvents:
			switch e.ID {
			case "<Escape>":
				return fmt.Errorf("UI Closed")
			case "q", "<C-c>":
				return fmt.Errorf("coin UI Closed")
			case "<Resize>":
				updateUI()
			}

		case data := <-priceChannel:
			myPage.PriceBox.Rows[0][0] = data + "$"
			ui.Render(myPage.PriceBox)

		case data := <-dataChannel:
			switch data.Type {

			case "HISTORY":
				// Update History graph
				price := data.PriceHistory
				myPage.ValueGraph.Data["Value"] = price
				myPage.ValueGraph.Labels["Max"] = fmt.Sprintf("%.2f %s", utils.MaxFloat64(price...), "$")
				myPage.ValueGraph.Labels["Min"] = fmt.Sprintf("%.2f %s", utils.MinFloat64(price...), "$")

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

				price := "NA"
				p, err := strconv.ParseFloat(data.CoinAssetData.Data.PriceUsd, 64)
				if err == nil {
					price = fmt.Sprintf("%.2f", p/1)
					myPage.ValueGraph.Labels["Value"] = fmt.Sprintf("%s %s", price, "$")
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
