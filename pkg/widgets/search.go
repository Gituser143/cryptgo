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

package widgets

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
)

const (
	FULL_BLOCK  = "█"
	EMPTY_BLOCK = " "
)

type SearchMenu struct {
	SearchString string
	*Table
	IsFull             bool
	SymbolDoesNotExist bool
	SearchData         [][]string
}

// NewSearchMenu is a constructor for the SearchMenu type
func NewSearchMenu() *SearchMenu {
	return &SearchMenu{
		Table: NewTable(),
	}
}

// Reset resets a search menu to its default values
func (search *SearchMenu) Reset() {
	search.SearchString = ""
	search.SearchData = [][]string{}
	search.IsFull = true
	search.SymbolDoesNotExist = false
	search.SelectedRow = 0
	search.Header = []string{}
}

func (search *SearchMenu) Resize(termWidth, termHeight int) {
	x1, y1 := termWidth/4, termHeight/4
	x2, y2 := 3*termWidth/4, 3*termHeight/4
	search.Table.SetRect(x1, y1, x2, y2)
}

func (search *SearchMenu) Draw(buf *ui.Buffer) {
	search.Table.Title = " Search "

	if search.SymbolDoesNotExist {
		search.Table.Header = []string{fmt.Sprintf(" Coin with symbol %s does not exist", search.SearchString)}
	} else if len(search.SearchData) > 0 {
		search.Table.Header = []string{fmt.Sprintf(" %v results", len(search.SearchData))}
	}

	search.IsFull = !search.IsFull

	suffix := EMPTY_BLOCK
	if search.IsFull && len(search.SearchString) > 0 {
		suffix = FULL_BLOCK
	}

	input := [][]string{{fmt.Sprintf(" ~ %s%s", search.SearchString, suffix), "", ""}, {"SYMBOL", "NAME", "PRICE"}}

	search.Table.Rows = append(input, search.SearchData...)

	search.Table.BorderStyle.Fg = ui.ColorCyan
	search.Table.BorderStyle.Bg = ui.ColorClear

	search.Table.RowStyle.Fg = ui.ColorCyan
	if search.SymbolDoesNotExist {
		search.Table.RowStyle.Fg = ui.ColorRed
	}

	search.Table.RowStyle.Bg = ui.ColorClear
	search.Table.ColResizer = func() {
		x := search.Table.Inner.Dx()
		search.Table.ColWidths = []int{x / 4, x / 2, x / 4}
	}
	search.Table.Draw(buf)
}
