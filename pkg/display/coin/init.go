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
			x / 2,
			x / 2,
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

	// Initialise Price Box
	page.PriceBox.BorderStyle.Fg = ui.ColorCyan
	page.PriceBox.TitleStyle.Fg = ui.ColorClear
	page.PriceBox.Header = []string{"Price", "Change %"}
	page.PriceBox.ColResizer = func() {
		x := page.PriceBox.Inner.Dx()
		page.PriceBox.ColWidths = []int{
			x / 2,
			x / 2,
		}
	}

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
}
