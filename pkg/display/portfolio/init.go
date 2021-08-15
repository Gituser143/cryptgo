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
	"math"

	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

// portfolioPage holds UI items for the portfolio page
type portfolioPage struct {
	Grid                *ui.Grid
	DetailsTable        *widgets.Table
	CoinTable           *widgets.Table
	BestPerformerTable  *widgets.Table
	WorstPerformerTable *widgets.Table
}

// performer holds best and worst perfomer details
type performer struct {
	BestVal   float64
	BestCoin  string
	WorstVal  float64
	WorstCoin string
}

// getEmptyPerformers returns a map with required durations as keys and empty performers
func getEmptyPerformers() map[string]performer {
	m := map[string]performer{
		"1h":  {BestVal: math.Inf(-1), WorstVal: math.Inf(1)},
		"24h": {BestVal: math.Inf(-1), WorstVal: math.Inf(1)},
		"7d":  {BestVal: math.Inf(-1), WorstVal: math.Inf(1)},
		"30d": {BestVal: math.Inf(-1), WorstVal: math.Inf(1)},
		"1y":  {BestVal: math.Inf(-1), WorstVal: math.Inf(1)},
	}

	return m
}

func newPortfolioPage() *portfolioPage {
	page := &portfolioPage{
		Grid:                ui.NewGrid(),
		DetailsTable:        widgets.NewTable(),
		CoinTable:           widgets.NewTable(),
		BestPerformerTable:  widgets.NewTable(),
		WorstPerformerTable: widgets.NewTable(),
	}

	page.init()

	return page
}

func (page *portfolioPage) init() {
	// Initialise Details table
	page.DetailsTable.Title = " Details "
	page.DetailsTable.BorderStyle.Fg = ui.ColorCyan
	page.DetailsTable.TitleStyle.Fg = ui.ColorClear
	page.DetailsTable.Header = []string{"Balance", ""}
	page.DetailsTable.ColResizer = func() {
		x := page.DetailsTable.Inner.Dx()
		page.DetailsTable.ColWidths = []int{
			x / 2,
			x / 2,
		}
	}
	page.DetailsTable.ShowCursor = false
	page.DetailsTable.CursorColor = ui.ColorCyan

	// Initialise CoinTable
	page.CoinTable.Title = " Coins "
	page.CoinTable.BorderStyle.Fg = ui.ColorCyan
	page.CoinTable.TitleStyle.Fg = ui.ColorClear
	page.CoinTable.Header = []string{"Rank", "Symbol", "Price", "Change % (1d)", "Holding", "Balance", "Holding %"}
	page.CoinTable.ColResizer = func() {
		x := page.CoinTable.Inner.Dx()
		page.CoinTable.ColWidths = []int{
			ui.MaxInt(5, 1*(x/10)),
			ui.MaxInt(5, 1*(x/10)),
			2 * (x / 10),
			2 * (x / 10),
			ui.MaxInt(5, 1*(x/10)),
			2 * (x / 10),
			ui.MaxInt(5, 1*(x/10)),
		}
	}
	page.CoinTable.ShowCursor = true
	page.CoinTable.CursorColor = ui.ColorCyan
	page.CoinTable.ChangeCol[3] = true

	// Initialise Best Performer Table
	page.BestPerformerTable.Title = " Best Performers "
	page.BestPerformerTable.BorderStyle.Fg = ui.ColorCyan
	page.BestPerformerTable.TitleStyle.Fg = ui.ColorClear
	page.BestPerformerTable.Header = []string{"Time", "Coin", "Change"}
	page.BestPerformerTable.ColResizer = func() {
		x := page.BestPerformerTable.Inner.Dx()
		page.BestPerformerTable.ColWidths = []int{
			3 * x / 10,
			3 * x / 10,
			3 * x / 10,
		}
	}
	page.BestPerformerTable.CursorColor = ui.ColorCyan
	page.BestPerformerTable.ChangeCol[2] = true

	// Initialise Worst Performer Table
	page.WorstPerformerTable.Title = " Worst Performers "
	page.WorstPerformerTable.BorderStyle.Fg = ui.ColorCyan
	page.WorstPerformerTable.TitleStyle.Fg = ui.ColorClear
	page.WorstPerformerTable.Header = []string{"Time", "Coin", "Change"}
	page.WorstPerformerTable.ColResizer = func() {
		x := page.WorstPerformerTable.Inner.Dx()
		page.WorstPerformerTable.ColWidths = []int{
			3 * x / 10,
			3 * x / 10,
			3 * x / 10,
		}
	}
	page.WorstPerformerTable.CursorColor = ui.ColorCyan
	page.WorstPerformerTable.ChangeCol[2] = true

	// Set Grid layout
	w, h := ui.TerminalDimensions()
	page.Grid.Set(
		ui.NewRow(0.3,
			ui.NewCol(0.2, page.DetailsTable),
			ui.NewCol(0.4, page.BestPerformerTable),
			ui.NewCol(0.4, page.WorstPerformerTable),
		),
		ui.NewRow(0.7, page.CoinTable),
	)

	page.Grid.SetRect(0, 0, w, h)
}
