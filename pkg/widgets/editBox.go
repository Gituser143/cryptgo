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

func rune_advance_len(r rune, pos int) int {
	if r == '\t' {
		return tabstop_length - pos%tabstop_length
	}
	return runewidth.RuneWidth(r)
}

func voffset_coffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset += 1
		voffset += rune_advance_len(r, voffset)
	}
	return
}

func byte_slice_grow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byte_slice_remove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byte_slice_insert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byte_slice_grow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferred_horizontal_threshold = 5
const tabstop_length = 8

type EditBox struct {
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
}

// Draws the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	const colred = termbox.ColorRed

	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.line_voffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstop_length
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, arrowRight,
				colred, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.line_voffset
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

	if eb.line_voffset != 0 {
		termbox.SetCell(x, y, arrowLeft, colred, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferred_horizontal_threshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.line_voffset != 0 {
		threshold = width - ht
	}
	if eb.cursor_voffset-eb.line_voffset >= threshold {
		eb.line_voffset = eb.cursor_voffset + (ht - width + 1)
	}

	if eb.line_voffset != 0 && eb.cursor_voffset-eb.line_voffset < ht {
		eb.line_voffset = eb.cursor_voffset - ht
		if eb.line_voffset < 0 {
			eb.line_voffset = 0
		}
	}
}

func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursor_boffset = boffset
	eb.cursor_voffset, eb.cursor_coffset = voffset_coffset(eb.text, boffset)
}

func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursor_boffset:])
}

func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursor_boffset])
}

func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursor_boffset - size)
}

func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursor_boffset + size)
}

func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursor_boffset]
}

func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byte_slice_insert(eb.text, eb.cursor_boffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursor_voffset - eb.line_voffset
}

var edit_box EditBox

const edit_box_width = 30

func redraw_all(symbol string) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	midy := h / 2
	midx := (w - edit_box_width) / 2

	// unicode box drawing chars around the edit box
	if runewidth.EastAsianWidth {
		termbox.SetCell(midx-1, midy, '|', coldef, coldef)
		termbox.SetCell(midx+edit_box_width, midy, '|', coldef, coldef)
		termbox.SetCell(midx-1, midy-1, '+', coldef, coldef)
		termbox.SetCell(midx-1, midy+1, '+', coldef, coldef)
		termbox.SetCell(midx+edit_box_width, midy-1, '+', coldef, coldef)
		termbox.SetCell(midx+edit_box_width, midy+1, '+', coldef, coldef)
		fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '-'})
		fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '-'})
	} else {
		termbox.SetCell(midx-1, midy, '│', coldef, coldef)
		termbox.SetCell(midx+edit_box_width, midy, '│', coldef, coldef)
		termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
		termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
		termbox.SetCell(midx+edit_box_width, midy-1, '┐', coldef, coldef)
		termbox.SetCell(midx+edit_box_width, midy+1, '┘', coldef, coldef)
		fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '─'})
		fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '─'})
	}

	edit_box.Draw(midx, midy, edit_box_width, 1)
	termbox.SetCursor(midx+edit_box.CursorX(), midy)

	title := fmt.Sprintf(" Enter Amount in %s ", symbol)
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

	redraw_all(symbol)
	defer termbox.HideCursor()
	for {
		for e := range ev {
			switch e.ID {
			case "<Escape>":
				return ""
			case "<Enter>":
				return string(edit_box.text)
			case "<Left>":
				edit_box.MoveCursorOneRuneBackward()
			case "<Right>":
				edit_box.MoveCursorOneRuneForward()
			case "<C-<Backspace>>", "<Backspace>":
				edit_box.DeleteRuneBackward()
			case "<Delete>", "<C-d>":
				edit_box.DeleteRuneForward()
			case "<Tab>":
				edit_box.InsertRune('\t')
			case "<Space>":
				edit_box.InsertRune(' ')
			case "<C-k>":
				edit_box.DeleteTheRestOfTheLine()
			case "<Home>":
				edit_box.MoveCursorToBeginningOfTheLine()
			case "<End>":
				edit_box.MoveCursorToEndOfTheLine()
			default:
				if len(e.ID) == 1 && []rune(e.ID)[0] != 0 {
					edit_box.InsertRune([]rune(e.ID)[0])
				}
			}
			redraw_all(symbol)
		}
	}
}
