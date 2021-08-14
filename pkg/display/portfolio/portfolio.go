package portfolio

import (
	"context"
	"fmt"
	"time"

	"github.com/Gituser143/cryptgo/pkg/api"
	ui "github.com/gizak/termui/v3"
)

func DisplayPortfolio(ctx context.Context, dataChannel chan api.PortfolioData) error {

	// Initialise UI
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialise termui: %v", err)
	}
	defer ui.Close()

	myPage := NewPortfolioPage()

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
		case <-dataChannel:
		case <-tick: // Refresh UI
			updateUI()
		}
	}

}
