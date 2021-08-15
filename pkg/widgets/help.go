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
	ui "github.com/gizak/termui/v3"
)

var allKeybindings = [][]string{
	{"Quit: q or <C-c>"},
	{""},
	{"Table Navigation"},
	{"  - k and <Up>: up"},
	{"  - j and <Down>: down"},
	{"  - <C-u>: half page up"},
	{"  - <C-d>: half page down"},
	{"  - <C-b>: full page up"},
	{"  - <C-f>: full page down"},
	{"  - gg and <Home>: jump to top"},
	{"  - G and <End>: jump to bottom"},
	{"  - f: focus favourites table"},
	{"  - F: focus coin table"},
	{""},
	{"Sorting"},
	{"  - Use column number to sort ascending."},
	{"  - Use <F-column number> to sort descending."},
	{"  - Eg: 1 to sort ascending on 1st Col and F1 for descending"},
	{""},
	{"Actions"},
	{"  - c: Select Currency (from popular list)"},
	{"  - C: Select Currency (from full list)"},
	{"  - e: Add/Edit coin to Portfolio"},
	{"  - P: View portfolio"},
	{"  - s: Star, save to favourites"},
	{"  - S: UnStar,remove from favourites"},
	{"  - <Enter>: View Coin Information"},
	{"  - %: Select Duration for Percentage Change"},
	{""},
	{"To close this prompt: <Esc>"},
}

var coinKeybindings = [][]string{
	{"Quit: q or <C-c>"},
	{""},
	{"Table Navigation"},
	{"  - d Change Interval Duration"},
	{"  - k and <Up>: up"},
	{"  - j and <Down>: down"},
	{"  - <C-u>: half page up"},
	{"  - <C-d>: half page down"},
	{"  - <C-b>: full page up"},
	{"  - <C-f>: full page down"},
	{"  - gg and <Home>: jump to top"},
	{"  - G and <End>: jump to bottom"},
	{"  - f: focus favourites table"},
	{"  - F: focus interval table"},
	{""},
	{"Sorting"},
	{"  - Use column number to sort ascending."},
	{"  - Use <F-column number> to sort descending."},
	{"  - Eg: 1 to sort ascending on 1st Col and F1 for descending"},
	{""},
	{""},
	{"To close this prompt: <Esc>"},
}

var portfolioKeybindings = [][]string{
	{"Quit: q or <C-c>"},
	{""},
	{"Table Navigation"},
	{"  - k and <Up>: up"},
	{"  - j and <Down>: down"},
	{"  - <C-u>: half page up"},
	{"  - <C-d>: half page down"},
	{"  - <C-b>: full page up"},
	{"  - <C-f>: full page down"},
	{"  - gg and <Home>: jump to top"},
	{"  - G and <End>: jump to bottom"},
	{""},
	{"Sorting"},
	{"  - Use column number to sort ascending."},
	{"  - Use <F-column number> to sort descending."},
	{"  - Eg: 1 to sort ascending on 1st Col and F1 for descending"},
	{""},
	{"Actions"},
	{"  - c: Select Currency (from popular list)"},
	{"  - C: Select Currency (from full list)"},
	{"  - e: Add/Edit coin to Portfolio"},
	{"  - <Enter>: View Coin Information"},
	{""},
	{"To close this prompt: <Esc>"},
}

// HelpMenu is a wrapper widget around a List meant
// to display the help menu for a command
type HelpMenu struct {
	*Table
	Keybindings [][]string
}

// NewHelpMenu is a constructor for the HelpMenu type
func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		Table: NewTable(),
	}
}

// Resize resizes the widget based on specified width
// and height
func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 50
	for _, line := range help.Keybindings {
		if textWidth < len(line[0]) {
			textWidth = len(line[0]) + 2
		}
	}
	textHeight := len(help.Keybindings) + 3
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

	help.Table.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget
func (help *HelpMenu) Draw(buf *ui.Buffer) {
	help.Table.Title = " Keybindings "
	help.Table.Rows = help.Keybindings
	help.Table.BorderStyle.Fg = ui.ColorCyan
	help.Table.BorderStyle.Bg = ui.ColorClear
	help.Table.ColResizer = func() {
		x := help.Table.Inner.Dx()
		help.Table.ColWidths = []int{x}
	}
	help.Table.Draw(buf)
}

// SelectHelpMenu selects the appropriate text
// based on the command for which the help page
// is needed
func (help *HelpMenu) SelectHelpMenu(page string) {
	help.IsHelp = true
	switch page {
	case "ALL":
		help.Keybindings = allKeybindings
	case "COIN":
		help.Keybindings = coinKeybindings
	case "PORTFOLIO":
		help.Keybindings = portfolioKeybindings
	}
}
