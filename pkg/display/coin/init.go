package coin

import (
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
)

type CoinPage struct {
	Grid            *ui.Grid
	FavouritesTable *widgets.Table
	ValueGraph      *widgets.LineGraph
	PriceBox        *widgets.Table
	VolumeGauge     *w.Gauge
	DetailsTable    *widgets.Table
	SupplyChart     *widgets.BarChart
}

func NewCoinPage() *CoinPage {
	page := &CoinPage{
		Grid:            ui.NewGrid(),
		FavouritesTable: widgets.NewTable(),
		ValueGraph:      widgets.NewLineGraph(),
		PriceBox:        widgets.NewTable(),
		VolumeGauge:     w.NewGauge(),
		DetailsTable:    widgets.NewTable(),
		SupplyChart:     widgets.NewBarChart(),
	}
	page.InitCoin()

	return page
}

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
	page.FavouritesTable.ShowCursor = true

	// Initialise Value Graph
	page.ValueGraph.TitleStyle = ui.NewStyle(ui.ColorClear)
	page.ValueGraph.HorizontalScale = 1
	page.ValueGraph.LineColors["Max"] = ui.ColorGreen
	page.ValueGraph.LineColors["Min"] = ui.ColorRed
	page.ValueGraph.LineColors["Value"] = ui.ColorBlue
	page.ValueGraph.BorderStyle.Fg = ui.ColorCyan
	page.ValueGraph.Data["Max"] = []float64{}
	page.ValueGraph.Data["Min"] = []float64{}

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
		ui.NewCol(0.33, page.FavouritesTable),
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
