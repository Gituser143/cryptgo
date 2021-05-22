package allcoin

import (
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

// AllCoinPage holds UI items for the home page
type AllCoinPage struct {
	Grid            *ui.Grid
	CoinTable       *widgets.Table
	TopCoinGraphs   []*widgets.LineGraph
	FavouritesTable *widgets.Table
}

// NewAllCoinPage creates, initialises and returns a pointer to an instance of AllCoinPage
func NewAllCoinPage() *AllCoinPage {
	coinGraphs := []*widgets.LineGraph{}
	for i := 0; i < 3; i++ {
		coinGraphs = append(coinGraphs, widgets.NewLineGraph())
	}

	page := &AllCoinPage{
		Grid:            ui.NewGrid(),
		CoinTable:       widgets.NewTable(),
		TopCoinGraphs:   coinGraphs,
		FavouritesTable: widgets.NewTable(),
	}

	page.InitAllCoin()

	return page
}

// InitAllCoin initialises the widgets of an AllCoinPage
func (page *AllCoinPage) InitAllCoin() {
	// Initialise CoinTable
	page.CoinTable.Title = " Coins "
	page.CoinTable.BorderStyle.Fg = ui.ColorCyan
	page.CoinTable.TitleStyle.Fg = ui.ColorClear
	page.CoinTable.Header = []string{"Rank", "Symbol", "Price", "Change %", "Supply / MaxSupply"}
	page.CoinTable.ColResizer = func() {
		x := page.CoinTable.Inner.Dx()
		page.CoinTable.ColWidths = []int{
			ui.MaxInt(8, x/5),
			ui.MaxInt(8, x/5),
			ui.MaxInt(15, x/5),
			ui.MaxInt(5, x/5),
			ui.MaxInt(20, x/5),
		}
	}
	page.CoinTable.ShowCursor = true
	page.CoinTable.CursorColor = ui.ColorCyan

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

	// Initialise Top Coin Graphs
	for i := 0; i < 3; i++ {
		page.TopCoinGraphs[i].TitleStyle = ui.NewStyle(ui.ColorClear)
		page.TopCoinGraphs[i].HorizontalScale = 1
		page.TopCoinGraphs[i].LineColors["Max"] = ui.ColorGreen
		page.TopCoinGraphs[i].LineColors["Min"] = ui.ColorRed
		page.TopCoinGraphs[i].LineColors["Value"] = ui.ColorBlue
		page.TopCoinGraphs[i].BorderStyle.Fg = ui.ColorCyan
		page.TopCoinGraphs[i].Data["Max"] = []float64{}
		page.TopCoinGraphs[i].Data["Min"] = []float64{}
	}

	// Set Grid layout
	w, h := ui.TerminalDimensions()
	page.Grid.Set(
		ui.NewRow(0.33,
			ui.NewCol(0.33, page.TopCoinGraphs[0]),
			ui.NewCol(0.33, page.TopCoinGraphs[1]),
			ui.NewCol(0.34, page.TopCoinGraphs[2]),
		),
		ui.NewRow(0.67,
			ui.NewCol(0.33, page.FavouritesTable),
			ui.NewCol(0.67, page.CoinTable),
		),
	)

	page.Grid.SetRect(0, 0, w, h)

}
