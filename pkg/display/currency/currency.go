package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

var rows [][]string

type Currency struct {
	ID             string `json:"id"`
	Symbol         string `json:"symbol"`
	CurrencySymbol string `json:"currencySymbol"`
	Type           string `json:"type"`
	RateUSD        string `json:"rateUSD"`
}

type CurrencyData struct {
	Data      Currency `json:"data"`
	Timestamp uint     `json:"timestamp"`
}

type CurrencyTable struct {
	*widgets.Table
}

func NewCurrencyPage() *CurrencyTable {
	c := &CurrencyTable{
		Table: widgets.NewTable(),
	}

	c.Table.Title = " Select Currency "
	c.Table.Header = []string{"Currency", "Symbol", "USD rate"}
	c.Table.Rows = rows
	c.Table.CursorColor = ui.ColorCyan
	c.Table.ColWidths = []int{5, 5, 5}
	c.Table.ColResizer = func() {
		x := c.Table.Inner.Dx()
		c.Table.ColWidths = []int{
			x / 3,
			x / 3,
			x / 3,
		}
	}
	return c
}

func (c *CurrencyTable) Resize(termWidth, termHeight int) {
	textWidth := 50

	textHeight := len(c.Table.Rows) + 3
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2
	if x < 0 {
		x = 0
		textWidth = termWidth
	}
	if y < 0 {
		y = 0
		textHeight = termHeight
	}

	c.Table.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget
func (c *CurrencyTable) Draw(buf *ui.Buffer) {
	c.Table.Draw(buf)
}

func (c *CurrencyTable) UpdateRows() {
	currencies := []string{
		"united-states-dollar",
		"euro",
		"japanese-yen",
		"british-pound-sterling",
		"indian-rupee",
		"australian-dollar",
		"canadian-dollar",
		"chinese-yuan-renminbi",
	}

	var wg sync.WaitGroup
	var m sync.Mutex

	client := &http.Client{}
	method := "GET"

	rows := [][]string{}

	for _, currency := range currencies {
		wg.Add(1)
		go func(c string, wg *sync.WaitGroup, m *sync.Mutex) {
			defer wg.Done()
			url := fmt.Sprintf("https://api.coincap.io/v2/rates/%s", c)

			req, err := http.NewRequest(method, url, nil)
			if err != nil {
				return
			}

			res, err := client.Do(req)
			if err != nil {
				return
			}
			defer res.Body.Close()

			data := CurrencyData{}

			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				return
			}

			rate, err := strconv.ParseFloat(data.Data.RateUSD, 64)
			if err != nil {
				return
			}

			row := []string{
				data.Data.Symbol,
				data.Data.CurrencySymbol,
				fmt.Sprintf("%.4f", rate),
			}

			m.Lock()
			rows = append(rows, row)
			m.Unlock()
		}(currency, &wg, &m)
	}

	wg.Wait()

	c.Table.Rows = rows
}
