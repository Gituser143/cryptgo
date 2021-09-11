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

// Widget borrowed from termbox-go demo
// https://github.com/nsf/termbox-go/blob/master/_demos/editbox.go

package widgets

import (
	"fmt"
	"unicode/utf8"

	ui "github.com/gizak/termui/v3"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func runeAdvanceLen(r rune, pos int) int {
	if r == '\t' {
		return tabstopLength - pos%tabstopLength
	}
	return runewidth.RuneWidth(r)
}

func voffsetCoffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset++
		voffset += runeAdvanceLen(r, voffset)
	}
	return
}

func byteSliceGrow(s []byte, desiredCap int) []byte {
	if cap(s) < desiredCap {
		ns := make([]byte, len(s), desiredCap)
		copy(ns, s)
		return ns
	}
	return s
}

func byteSliceRemove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferredHorizontalThreshold = 5
const tabstopLength = 8

// EditBox helps user provide input through a text box like widget
type EditBox struct {
	text          []byte
	lineVoffset   int
	cursorBoffset int // cursor offset in bytes
	cursorVoffset int // visual cursor offset in termbox cells
	cursorCoffset int // cursor offset in unicode code points
}

// Draw helps draw the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	const colred = termbox.ColorRed

	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.lineVoffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstopLength
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, arrowRight,
				colred, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.lineVoffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if eb.lineVoffset != 0 {
		termbox.SetCell(x, y, arrowLeft, colred, coldef)
	}
}

// AdjustVOffset adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferredHorizontalThreshold
	maxHorizontalThreshold := (width - 1) / 2
	if ht > maxHorizontalThreshold {
		ht = maxHorizontalThreshold
	}

	threshold := width - 1
	if eb.lineVoffset != 0 {
		threshold = width - ht
	}
	if eb.cursorVoffset-eb.lineVoffset >= threshold {
		eb.lineVoffset = eb.cursorVoffset + (ht - width + 1)
	}

	if eb.lineVoffset != 0 && eb.cursorVoffset-eb.lineVoffset < ht {
		eb.lineVoffset = eb.cursorVoffset - ht
		if eb.lineVoffset < 0 {
			eb.lineVoffset = 0
		}
	}
}

func (eb *EditBox) moveCursorTo(boffset int) {
	eb.cursorBoffset = boffset
	eb.cursorVoffset, eb.cursorCoffset = voffsetCoffset(eb.text, boffset)
}

func (eb *EditBox) runeUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursorBoffset:])
}

func (eb *EditBox) runeBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursorBoffset])
}

func (eb *EditBox) moveCursorOneRuneBackward() {
	if eb.cursorBoffset == 0 {
		return
	}
	_, size := eb.runeBeforeCursor()
	eb.moveCursorTo(eb.cursorBoffset - size)
}

func (eb *EditBox) moveCursorOneRuneForward() {
	if eb.cursorBoffset == len(eb.text) {
		return
	}
	_, size := eb.runeUnderCursor()
	eb.moveCursorTo(eb.cursorBoffset + size)
}

func (eb *EditBox) moveCursorToBeginningOfTheLine() {
	eb.moveCursorTo(0)
}

func (eb *EditBox) moveCursorToEndOfTheLine() {
	eb.moveCursorTo(len(eb.text))
}

func (eb *EditBox) deleteRuneBackward() {
	if eb.cursorBoffset == 0 {
		return
	}

	eb.moveCursorOneRuneBackward()
	_, size := eb.runeUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.cursorBoffset, eb.cursorBoffset+size)
}

func (eb *EditBox) deleteRuneForward() {
	if eb.cursorBoffset == len(eb.text) {
		return
	}
	_, size := eb.runeUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.cursorBoffset, eb.cursorBoffset+size)
}

func (eb *EditBox) deleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursorBoffset]
}

func (eb *EditBox) insertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.cursorBoffset, buf[:n])
	eb.moveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of lineVoffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) cursorX() int {
	return eb.cursorVoffset - eb.lineVoffset
}

var editBox EditBox

const editBoxWidth = 30

func redrawAll(symbol string) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	midy := h / 2
	midx := (w - editBoxWidth) / 2

	// unicode box drawing chars around the edit box
	if runewidth.EastAsianWidth {
		termbox.SetCell(midx-1, midy, '|', coldef, coldef)
		termbox.SetCell(midx+editBoxWidth, midy, '|', coldef, coldef)
		termbox.SetCell(midx-1, midy-1, '+', coldef, coldef)
		termbox.SetCell(midx-1, midy+1, '+', coldef, coldef)
		termbox.SetCell(midx+editBoxWidth, midy-1, '+', coldef, coldef)
		termbox.SetCell(midx+editBoxWidth, midy+1, '+', coldef, coldef)
		fill(midx, midy-1, editBoxWidth, 1, termbox.Cell{Ch: '-'})
		fill(midx, midy+1, editBoxWidth, 1, termbox.Cell{Ch: '-'})
	} else {
		termbox.SetCell(midx-1, midy, '│', coldef, coldef)
		termbox.SetCell(midx+editBoxWidth, midy, '│', coldef, coldef)
		termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
		termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
		termbox.SetCell(midx+editBoxWidth, midy-1, '┐', coldef, coldef)
		termbox.SetCell(midx+editBoxWidth, midy+1, '┘', coldef, coldef)
		fill(midx, midy-1, editBoxWidth, 1, termbox.Cell{Ch: '─'})
		fill(midx, midy+1, editBoxWidth, 1, termbox.Cell{Ch: '─'})
	}

	editBox.Draw(midx, midy, editBoxWidth, 1)
	termbox.SetCursor(midx+editBox.cursorX(), midy)

	title := " Enter Symbol/Name of coin "
	if symbol != "" {
		title = fmt.Sprintf(" Enter Amount in %s ", symbol)
	}
	tbprint(midx, midy-1, coldef, coldef, title)
	tbprint(midx, midy+2, coldef, coldef, "ESC to Close")
	tbprint(midx, midy+3, coldef, coldef, "Enter to Save")

	termbox.Flush()
}

var arrowLeft = '←'
var arrowRight = '→'

func init() {
	if runewidth.EastAsianWidth {
		arrowLeft = '<'
		arrowRight = '>'
	}
}

// DrawEdit draws an editbox and returns input passed to the box
func DrawEdit(ev <-chan ui.Event, symbol string) string {
	termbox.SetInputMode(termbox.InputEsc)

	redrawAll(symbol)
	defer termbox.HideCursor()
	for {
		for e := range ev {
			switch e.ID {
			case "<Escape>":
				return ""
			case "<Enter>":
				return string(editBox.text)
			case "<Left>":
				editBox.moveCursorOneRuneBackward()
			case "<Right>":
				editBox.moveCursorOneRuneForward()
			case "<C-<Backspace>>", "<Backspace>":
				editBox.deleteRuneBackward()
			case "<Delete>", "<C-d>":
				editBox.deleteRuneForward()
			case "<Tab>":
				editBox.insertRune('\t')
			case "<Space>":
				editBox.insertRune(' ')
			case "<C-k>":
				editBox.deleteTheRestOfTheLine()
			case "<Home>":
				editBox.moveCursorToBeginningOfTheLine()
			case "<End>":
				editBox.moveCursorToEndOfTheLine()
			default:
				if len(e.ID) == 1 && []rune(e.ID)[0] != 0 {
					editBox.insertRune([]rune(e.ID)[0])
				}
			}
			redrawAll(symbol)
		}
	}
}
