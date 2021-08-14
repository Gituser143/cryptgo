package portfolio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	uw "github.com/Gituser143/cryptgo/pkg/display/utilitywidgets"
	"github.com/Gituser143/cryptgo/pkg/utils"
	ui "github.com/gizak/termui/v3"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

func DisplayPortfolio(ctx context.Context, dataChannel chan api.AssetData, sendData *bool) error {

	// Initialise UI
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialise termui: %v", err)
	}
	defer ui.Close()

	myPage := NewPortfolioPage()

	// currency variables
	currencyWidget := uw.NewCurrencyPage()
	currencyID := utils.GetCurrency()
	currencyID, currency, currencyVal := currencyWidget.Get(currencyID)

	// get portfolio details
	portfolioMap := utils.GetPortfolio()

	// get performers map
	performersMap := GetEmptyPerformers()

	// get favourites
	favourites := utils.GetFavourites()

	defer func() {
		utils.SaveMetadata(favourites, currencyID, portfolioMap)
	}()

	// UpdateUI to refresh UI
	updateUI := func() {
		// Get Terminal Dimensions
		w, h := ui.TerminalDimensions()
		myPage.Grid.SetRect(0, 0, w, h)

		// Clear UI
		ui.Clear()

		// Render required widgets
		// switch utilitySelected {
		// case "HELP":
		// 	help.Resize(w, h)
		// 	ui.Render(help)
		// case "PORTFOLIO":
		// 	portfolioTable.Resize(w, h)
		// 	ui.Render(portfolioTable)
		// case "CURRENCY":
		// 	currencyWidget.Resize(w, h)
		// 	ui.Render(currencyWidget)
		// case "CHANGE":
		// 	changePercentWidget.Resize(w, h)
		// 	ui.Render(changePercentWidget)
		// default:
		ui.Render(myPage.Grid)
		// }
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
			case "q", "<C-c>":
				return fmt.Errorf("UI Closed")

			case "<Resize>":
				updateUI()
			}
		case data := <-dataChannel:
			rows := [][]string{}

			// Update currency headers
			myPage.CoinTable.Header[2] = fmt.Sprintf("Price (%s)", currency)
			myPage.CoinTable.Header[5] = fmt.Sprintf("Balance (%s)", currency)

			portfolioTotal := 0.0
			durations := []string{"1h", "24h", "7d", "30d", "1y"}
			balanceMap := map[string]float64{}

			// Iterate over coin assets
			for _, val := range data.AllCoinData {
				if portfolioHolding, ok := portfolioMap[val.ID]; ok {
					// Get coin price
					price := fmt.Sprintf("%.2f", val.CurrentPrice/currencyVal)

					// Get change %
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
						"holdingPercent",
					})

					portfolioTotal += balanceFloat
					balanceMap[symbol] = balanceFloat

					for _, duration := range durations {
						val := api.GetPercentageChangeForDuration(val, duration)

						if val > performersMap[duration].BestVal {
							performersMap[duration] = Performer{
								BestVal:   val,
								BestCoin:  symbol,
								WorstVal:  performersMap[duration].WorstVal,
								WorstCoin: performersMap[duration].WorstCoin,
							}
						}

						if val < performersMap[duration].WorstVal {
							performersMap[duration] = Performer{
								BestVal:   performersMap[duration].BestVal,
								BestCoin:  performersMap[duration].BestCoin,
								WorstVal:  val,
								WorstCoin: symbol,
							}
						}
					}
				}
			}

			for i, row := range rows {
				symbol := row[1]
				rows[i][6] = fmt.Sprintf("%.2f %%", (balanceMap[symbol]/portfolioTotal)*100)
			}

			// Update coin table
			myPage.CoinTable.Rows = rows

			// Update details table
			myPage.DetailsTable.Header = []string{
				"Balance",
				fmt.Sprintf("%.2f", portfolioTotal),
			}
			myPage.DetailsTable.Rows = [][]string{
				{"Currency", currency},
				{"Coins", fmt.Sprintf("%d", len(portfolioMap))},
			}

			// Update Best Performers Table
			BestPerformerRows := [][]string{}
			WorstPerformerRows := [][]string{}

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

			myPage.BestPerformerTable.Rows = BestPerformerRows
			myPage.WorstPerformerTable.Rows = WorstPerformerRows

		case <-tick: // Refresh UI
			updateUI()
		}
	}

}
