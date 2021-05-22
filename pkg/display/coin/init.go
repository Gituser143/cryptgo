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
	w "github.com/gizak/termui/v3/widgets"
)

// CoinPage holds UI items for a coin page
type CoinPage struct {
	Grid            *ui.Grid
	FavouritesTable *widgets.Table
	ValueGraph      *widgets.LineGraph
	IntervalTable   *widgets.Table
	PriceBox        *widgets.Table
	VolumeGauge     *w.Gauge
	DetailsTable    *widgets.Table
	SupplyChart     *widgets.BarChart
}

// NewCoinPage creates, initialises and returns a pointer to an instance of CoinPage
func NewCoinPage() *CoinPage {
	page := &CoinPage{
		Grid:            ui.NewGrid(),
		FavouritesTable: widgets.NewTable(),
		ValueGraph:      widgets.NewLineGraph(),
		IntervalTable:   widgets.NewTable(),
		PriceBox:        widgets.NewTable(),
		VolumeGauge:     w.NewGauge(),
		DetailsTable:    widgets.NewTable(),
		SupplyChart:     widgets.NewBarChart(),
	}
	page.InitCoin()

	return page
}

// InitCoin initialises the widgets of an CoinPage
func (page *CoinPage) InitCoin() {
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

	// Initialise Interval Table
	page.IntervalTable.Title = " Graph Interval "
	page.IntervalTable.BorderStyle.Fg = ui.ColorCyan
	page.IntervalTable.TitleStyle.Fg = ui.ColorClear
	page.IntervalTable.Header = []string{"Interval"}
	page.IntervalTable.Rows = [][]string{
		{"1  min"},
		{"5  min"},
		{"15 min"},
		{"30 min"},
		{"1  hour"},
		{"2  hour"},
		{"6  hour"},
		{"12 hour"},
		{"1  day"},
	}
	page.IntervalTable.ColResizer = func() {
		x := page.IntervalTable.Inner.Dx()
		page.IntervalTable.ColWidths = []int{x}
	}
	page.IntervalTable.CursorColor = ui.ColorCyan

	// Initialise Price Box
	page.PriceBox.Title = " Price & Change "
	page.PriceBox.BorderStyle.Fg = ui.ColorCyan
	page.PriceBox.TitleStyle.Fg = ui.ColorClear
	page.PriceBox.Header = []string{"Live Price", "Change %"}
	page.PriceBox.ColResizer = func() {
		x := page.PriceBox.Inner.Dx()
		page.PriceBox.ColWidths = []int{
			6 * x / 10,
			4 * x / 10,
		}
	}
	page.PriceBox.Rows = [][]string{{"", ""}}

	// Initialise Details Table
	page.DetailsTable.Title = " Details "
	page.DetailsTable.BorderStyle.Fg = ui.ColorCyan
	page.DetailsTable.TitleStyle.Fg = ui.ColorClear
	page.DetailsTable.ColResizer = func() {
		x := page.DetailsTable.Inner.Dx()
		page.DetailsTable.ColWidths = []int{
			x / 2,
			x / 2,
		}
	}

	// Initialise Volume Used gauge
	page.VolumeGauge.Title = " 24 Hr Volume Used "
	page.VolumeGauge.BorderStyle.Fg = ui.ColorCyan
	page.VolumeGauge.TitleStyle.Fg = ui.ColorClear
	page.VolumeGauge.BarColor = ui.ColorCyan
	page.VolumeGauge.Percent = 0

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
			ui.NewRow(0.5, page.IntervalTable),
		),
		ui.NewCol(0.67,
			ui.NewRow(0.5, page.ValueGraph),
			ui.NewRow(0.2,
				ui.NewCol(0.5, page.PriceBox),
				ui.NewCol(0.5, page.VolumeGauge),
			),
			ui.NewRow(0.3,
				ui.NewCol(0.5, page.DetailsTable),
				ui.NewCol(0.5, page.SupplyChart),
			),
		),
	)

	page.Grid.SetRect(0, 0, w, h)
}
