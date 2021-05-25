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
	"github.com/gizak/termui/v3/widgets"
)

var allKeybindings = []string{
	"Quit: q or <C-c>",
	"",
	"[Table Navigation](fg:white)",
	"  - k and <Up>: up",
	"  - j and <Down>: down",
	"  - <C-u>: half page up",
	"  - <C-d>: half page down",
	"  - <C-b>: full page up",
	"  - <C-f>: full page down",
	"  - gg and <Home>: jump to top",
	"  - G and <End>: jump to bottom",
	"  - f: focus favourites table",
	"  - F: focus coin table",
	"",
	"[Sorting](fg:white)",
	"  - Use column number to sort ascending.",
	"  - Use <F-column number> to sort descending.",
	"  - Eg: 1 to sort ascedning on 1st Col and F1 for descending",
	"",
	"[Actions](fg:white)",
	"  - c: Select Currency (from popular list)",
	"  - C: Select Currency (from full list)",
	"  - e: Add/Edit coin to Portfolio",
	"  - s: Star, save to favourites",
	"  - S: UnStar,remove from favourites",
	"  - <Enter>: View Coin Information",
	"",
	"[To close this prompt: <Esc>](fg:white)",
}

var coinKeybindings = []string{
	"Quit: q or <C-c>",
	"",
	"[Table Navigation](fg:white)",
	"  - k and <Up>: up",
	"  - j and <Down>: down",
	"  - <C-u>: half page up",
	"  - <C-d>: half page down",
	"  - <C-b>: full page up",
	"  - <C-f>: full page down",
	"  - gg and <Home>: jump to top",
	"  - G and <End>: jump to bottom",
	"  - f: focus favourites table",
	"  - F: focus interval table",
	"",
	"[Sorting](fg:white)",
	"  - Use column number to sort ascending.",
	"  - Use <F-column number> to sort descending.",
	"  - Eg: 1 to sort ascedning on 1st Col and F1 for descending",
	"",
	"[Actions (Interval Table)](fg:white)",
	"  - c: Select Currency (from popular list)",
	"  - C: Select Currency (from full list)",
	"  - e: Add/Edit coin to Portfolio",
	"  - <Enter>: Set Interval",
	"",
	"[To close this prompt: <Esc>](fg:white)",
}

// HelpMenu is a wrapper widget around a List meant
// to display the help menu for a command
type HelpMenu struct {
	*widgets.List
	Keybindings []string
}

// NewHelpMenu is a constructor for the HelpMenu type
func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		List: widgets.NewList(),
	}
}

// Resize resizes the widget based on specified width
// and height
func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 50
	for _, line := range help.Keybindings {
		if textWidth < len(line) {
			textWidth = len(line) + 2
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

	help.List.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw puts the required text into the widget
func (help *HelpMenu) Draw(buf *ui.Buffer) {
	help.List.Title = " Keybindings "

	help.List.Rows = help.Keybindings
	help.List.TextStyle = ui.NewStyle(ui.ColorYellow)
	help.List.WrapText = false
	help.List.Draw(buf)
}

// SelectHelpMenu selects the appropriate text
// based on the command for which the help page
// is needed
func (help *HelpMenu) SelectHelpMenu(page string) {
	switch page {
	case "ALL":
		help.Keybindings = allKeybindings
	case "COIN":
		help.Keybindings = coinKeybindings
	}
}
