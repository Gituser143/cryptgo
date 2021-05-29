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

package changePercentageDuration

import (
	"github.com/Gituser143/cryptgo/pkg/widgets"
	ui "github.com/gizak/termui/v3"
)

var rows [][]string = [][]string{{"1 Hour"}, {"24 Hours"}, {"7 Days"}, {"14 Days"}, {"30 Days"}, {"200 Days"}, {"1 Year"}}

var DurationMap map[string]string = map[string]string{
	"1 Hour":   "1h",
	"24 Hours": "24h",
	"7 Days":   "7d",
	"14 Days":  "14d",
	"30 Days":  "30d",
	"200 Days": "200d",
	"1 Year":   "1y",
}

type ChangePercentageDurationTable struct {
	*widgets.Table
}

// NewCurrencyPage creates, initialises and returns a pointer to an instance of CurrencyTable
func NewChangePercentPage() *ChangePercentageDurationTable {
	c := &ChangePercentageDurationTable{
		Table: widgets.NewTable(),
	}

	c.Table.Title = " Select Duration for Percentage Change "
	c.Table.Header = []string{"Duration"}
	c.Table.Rows = rows
	c.Table.CursorColor = ui.ColorCyan
	c.Table.ShowCursor = true
	c.Table.ColWidths = []int{5}
	c.Table.ColResizer = func() {
		x := c.Table.Inner.Dx()
		c.Table.ColWidths = []int{
			x,
		}
	}
	return c
}

func (c *ChangePercentageDurationTable) Resize(termWidth, termHeight int) {
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
func (c *ChangePercentageDurationTable) Draw(buf *ui.Buffer) {
	c.Table.Draw(buf)
}
