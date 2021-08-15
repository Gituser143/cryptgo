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
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

// coinPage holds UI items for a coin page
type coinPage struct {
	Grid            *ui.Grid
	FavouritesTable *widgets.Table
	ValueGraph      *widgets.LineGraph
	DetailsTable    *widgets.Table
	ChangesTable    *widgets.Table
	PriceBox        *widgets.Table
	ExplorerTable   *widgets.Table
	SupplyChart     *widgets.BarChart
}

// newcoinPage creates, initialises and returns a pointer to an instance of coinPage
func newCoinPage() *coinPage {
	page := &coinPage{
		Grid:            ui.NewGrid(),
		FavouritesTable: widgets.NewTable(),
		ValueGraph:      widgets.NewLineGraph(),
		DetailsTable:    widgets.NewTable(),
		ChangesTable:    widgets.NewTable(),
		PriceBox:        widgets.NewTable(),
		ExplorerTable:   widgets.NewTable(),
		SupplyChart:     widgets.NewBarChart(),
	}
	page.init()

	return page
}

// init initialises the widgets of an coinPage
func (page *coinPage) init() {
	// Initialise Favourites table
	page.FavouritesTable.Title = " Favourites "
	page.FavouritesTable.BorderStyle.Fg = ui.ColorCyan
	page.FavouritesTable.TitleStyle.Fg = ui.ColorClear
	page.FavouritesTable.Header = []string{"Symbol", "Price"}
	page.FavouritesTable.ColResizer = func() {
		x := page.FavouritesTable.Inner.Dx()
		page.FavouritesTable.ColWidths = []int{
			4 * x / 10,
			6 * x / 10,
		}
	}
	page.FavouritesTable.CursorColor = ui.ColorCyan

	// Initialise Value Graph
	page.ValueGraph.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.ValueGraph.HorizontalScale = 1
	page.ValueGraph.LineColors["Max"] = ui.ColorGreen
	page.ValueGraph.LineColors["Min"] = ui.ColorRed
	page.ValueGraph.LineColors["Value"] = ui.ColorBlue
	page.ValueGraph.BorderStyle.Fg = ui.ColorCyan
	page.ValueGraph.Data["Max"] = []float64{}
	page.ValueGraph.Data["Min"] = []float64{}

	// Initialise Details Table
	page.DetailsTable.Title = " Details "
	page.DetailsTable.BorderStyle.Fg = ui.ColorCyan
	page.DetailsTable.TitleStyle.Fg = ui.ColorClear
	page.DetailsTable.ColResizer = func() {
		x := page.DetailsTable.Inner.Dx()
		page.DetailsTable.ColWidths = []int{
			4 * x / 10,
			6 * x / 10,
		}
	}

	// Initialise Change Table
	page.ChangesTable.Title = " Changes "
	page.ChangesTable.BorderStyle.Fg = ui.ColorCyan
	page.ChangesTable.BorderStyle.Bg = ui.ColorClear
	page.ChangesTable.Header = []string{"Interval", "Change"}
	page.ChangesTable.ColResizer = func() {
		x := page.ChangesTable.Inner.Dx()
		page.ChangesTable.ColWidths = []int{
			4 * x / 10,
			6 * x / 10,
		}
	}
	page.ChangesTable.ChangeCol[1] = true
	page.ChangesTable.ShowCursor = false

	// Initialise Price Box
	page.PriceBox.Title = " Live Price "
	page.PriceBox.BorderStyle.Fg = ui.ColorCyan
	page.PriceBox.TitleStyle.Fg = ui.ColorClear
	page.PriceBox.Header = []string{"Price", "24H High", "24H Low"}
	page.PriceBox.ColResizer = func() {
		x := page.PriceBox.Inner.Dx()
		page.PriceBox.ColWidths = []int{
			4 * x / 10,
			3 * x / 10,
			3 * x / 10,
		}
	}
	page.PriceBox.Rows = [][]string{{"NA", "", ""}}
	page.PriceBox.ColColor[1] = ui.ColorGreen
	page.PriceBox.ColColor[2] = ui.ColorRed

	// Initialise Explorer Table
	page.ExplorerTable.Title = " Explorers "
	page.ExplorerTable.BorderStyle.Fg = ui.ColorCyan
	page.ExplorerTable.TitleStyle.Fg = ui.ColorClear
	page.ExplorerTable.Header = []string{"Links"}
	page.ExplorerTable.ColResizer = func() {
		x := page.ExplorerTable.Inner.Dx()
		page.ExplorerTable.ColWidths = []int{x}
	}
	page.ExplorerTable.CursorColor = ui.ColorCyan

	// Initalise Bar Graph
	page.SupplyChart.Title = " Supply "
	page.SupplyChart.Data = []float64{0, 0}
	page.SupplyChart.Labels = []string{"Supply", "Max Supply"}
	page.SupplyChart.BorderStyle.Fg = ui.ColorCyan
	page.SupplyChart.TitleStyle.Fg = ui.ColorClear
	page.SupplyChart.BarWidth = 9
	page.SupplyChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorCyan}
	page.SupplyChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorClear)}
	page.SupplyChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// Set Grid layout
	w, h := ui.TerminalDimensions()
	page.Grid.Set(
		ui.NewCol(0.33,
			ui.NewRow(0.5, page.FavouritesTable),
			ui.NewRow(0.5, page.DetailsTable),
		),
		ui.NewCol(0.67,
			ui.NewRow(0.5, page.ValueGraph),
			ui.NewRow(0.5,
				ui.NewCol(0.5,
					ui.NewRow(0.4, page.PriceBox),
					ui.NewRow(0.6, page.ChangesTable),
				),
				ui.NewCol(0.5,
					ui.NewRow(0.5, page.ExplorerTable),
					ui.NewRow(0.5, page.SupplyChart),
				),
			),
		),
	)

	page.Grid.SetRect(0, 0, w, h)
}
