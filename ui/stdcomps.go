// Standard components

package ui

import (
	"unicode"
	)

// A text component.  This simply renders a text string.
type TextView struct {

	// The string to render
	Text string
}

// Minimum dimensions
func (tv *TextView) Remeasure(w, h int) (int, int) {
	return w, 1
}

// Status bar redraw
func (tv *TextView) Redraw(context *DrawContext) {
	context.SetFgAttr(0)
	context.SetBgAttr(0)

	context.HorizRule(0, ' ')
	context.Print(0, 0, tv.Text)
	context.HideCursor()
}

// Status bar component.  This component displays text on the left and right of it's
// allocated space.
type StatusBar struct {
	Left  string // Left aligned string
	Right string // Right aligned string
}

// Minimum dimensions
func (sbar *StatusBar) Remeasure(w, h int) (int, int) {
	return w, 1
}

// Status bar redraw
func (sbar *StatusBar) Redraw(context *DrawContext) {
	context.SetFgAttr(AttrReverse)
	context.SetBgAttr(AttrReverse)

	context.HorizRule(0, ' ')
	context.Print(0, 0, sbar.Left)
	context.PrintRight(context.W, 0, sbar.Right)
}

// A single-text entry component.
type TextEntry struct {
	Prompt string

	value         string
	cursorOffset  int
	displayOffset int

	// Called when the user presses Enter
	OnEntry func(val string)

	// Called when the user presses Esc or CtrlC
	OnCancel func()
}

func (te *TextEntry) Remeasure(w, h int) (int, int) {
	return w, 1
}

func (te *TextEntry) Redraw(context *DrawContext) {
	context.HorizRule(0, ' ')
	valueOffsetX := 0
	displayOffsetX := te.calculateDisplayOffset(context.W)

	if te.Prompt != "" {
		context.SetFgAttr(ColorDefault | AttrBold)
		context.Print(0, 0, te.Prompt)
		context.SetFgAttr(ColorDefault)

		valueOffsetX = len(te.Prompt)
	}

	context.Print(valueOffsetX, 0, te.value[displayOffsetX:intMin(displayOffsetX+context.W, len(te.value))])
	context.SetCursorPosition(te.cursorOffset+valueOffsetX-displayOffsetX, 0)

	//context.Print(0, 0, fmt.Sprintf("%d,%d", te.cursorOffset, displayOffsetX))
}

func (te *TextEntry) calculateDisplayOffset(displayWidth int) int {
	if te.Prompt != "" {
		displayWidth -= len(te.Prompt)
	}
	virtualCursorOffset := te.cursorOffset - te.displayOffset

	if virtualCursorOffset >= displayWidth {
		te.displayOffset = te.cursorOffset - displayWidth + 10
	} else if virtualCursorOffset < 0 {
		te.displayOffset = intMax(te.cursorOffset-displayWidth+1, 0)
	}

	return te.displayOffset
}

// SetValue sets the value of the text entry
func (te *TextEntry) SetValue(val string) {
	te.value = val
	te.cursorOffset = len(val)
}

func (te *TextEntry) KeyPressed(key rune, mod int) {
	if (key >= ' ') && (key <= '~') {
		te.insertRune(key)
	} else if key == KeyArrowLeft {
		te.moveCursorBy(-1)
	} else if key == KeyArrowRight {
		te.moveCursorBy(1)
	} else if key == KeyHome {
		te.moveCursorTo(0)
	} else if key == KeyEnd {
		te.moveCursorTo(len(te.value))
	} else if (key == KeyBackspace) || (key == KeyBackspace2) {
		if mod&ModKeyAlt != 0 {
			te.backspaceWhile(unicode.IsSpace)
			te.backspaceWhile(func(r rune) bool { return !unicode.IsSpace(r) })
		} else {
			te.backspace()
		}
	} else if key == KeyCtrlK {
		te.killLine()
	} else if key == KeyDelete {
		te.removeCharAtPos(te.cursorOffset)
	} else if key == KeyEnter {
		if te.OnEntry != nil {
			te.OnEntry(te.value)
		}
	} else if key == KeyCtrlC {
		if te.OnCancel != nil {
			te.OnCancel()
		}
	}

	//panic(fmt.Sprintf("Entered key: '%x', mod: '%x'", key, mod))
}

// Backspace
func (te *TextEntry) backspace() {
	te.removeCharAtPos(te.cursorOffset - 1)
	te.moveCursorBy(-1)
}

// Backspace while the character underneith the cursor matches the guard
func (te *TextEntry) backspaceWhile(guard func(r rune) bool) {
	for te.cursorOffset > 0 {
		ch := rune(te.value[te.cursorOffset-1])
		if guard(ch) {
			te.backspace()
		} else {
			break
		}
	}
}

// Kill the line.  If the cursor is at the end of the line, kill to the start.
// Otherwise, trim the line.
func (te *TextEntry) killLine() {
	if te.cursorOffset < len(te.value) {
		te.value = te.value[:te.cursorOffset]
	} else {
		te.value = ""
		te.cursorOffset = 0
	}
}

// Inserts a rune at the cursor position
func (te *TextEntry) insertRune(key rune) {
	if te.cursorOffset >= len(te.value) {
		te.value += string(key)
	} else {
		te.value = te.value[:te.cursorOffset] + string(key) + te.value[te.cursorOffset:]
	}
	te.moveCursorBy(1)
}

// Remove the character at a specific position
func (te *TextEntry) removeCharAtPos(pos int) {
	if (pos >= 0) && (pos < len(te.value)) {
		te.value = te.value[:pos] + te.value[pos+1:]
	}
}

// Move the cursor
func (te *TextEntry) moveCursorBy(byX int) {
	te.moveCursorTo(te.cursorOffset + byX)
}

func (te *TextEntry) moveCursorTo(toX int) {
	te.cursorOffset = intMinMax(toX, 0, len(te.value))
}
