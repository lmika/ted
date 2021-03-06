// The terminal-box driver

package ui

import (
	"github.com/gdamore/tcell/termbox"
)

type TermboxDriver struct {
}

// Initializes the driver.  Returns an error if there was an error
func (td *TermboxDriver) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	termbox.SetInputMode(termbox.InputAlt)
	return nil
}

// Closes the driver
func (td *TermboxDriver) Close() {
	termbox.Close()
}

// Returns the size of the window.
func (td *TermboxDriver) Size() (int, int) {
	return termbox.Size()
}

// Sets the value of a specific cell
func (td *TermboxDriver) SetCell(x, y int, ch rune, fg, bg Attribute) {
	termbox.SetCell(x, y, ch, termbox.Attribute(fg), termbox.Attribute(bg))
}

// Hide the cursor
func (td *TermboxDriver) HideCursor() {
	termbox.HideCursor()
}

// Synchronizes the internal buffer with the real buffer
func (td *TermboxDriver) Sync() {
	termbox.Flush()
}

// Wait for an event
func (td *TermboxDriver) WaitForEvent() Event {
	tev := termbox.PollEvent()

	switch tev.Type {
	case termbox.EventResize:
		return Event{EventResize, 0, 0}
	case termbox.EventKey:
		mod := 0
		if tev.Mod&termbox.ModAlt != 0 {
			mod = ModKeyAlt
		}
		if tev.Ch != 0 {
			return Event{EventKeyPress, mod, tev.Ch}
		} else if spec, hasSpec := termboxKeysToSpecialKeys[tev.Key]; hasSpec {
			return Event{EventKeyPress, mod, spec}
		} else {
			return Event{EventNone, mod, 0}
		}
	default:
		return Event{EventNone, 0, 0}
	}
}

// Move the position of the cursor
func (td *TermboxDriver) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

// Map from termbox Keys to driver key runes
var termboxKeysToSpecialKeys = map[termbox.Key]rune{
	termbox.KeySpace: ' ',

	termbox.KeyF1:         KeyF1,
	termbox.KeyF2:         KeyF2,
	termbox.KeyF3:         KeyF3,
	termbox.KeyF4:         KeyF4,
	termbox.KeyF5:         KeyF5,
	termbox.KeyF6:         KeyF6,
	termbox.KeyF7:         KeyF7,
	termbox.KeyF8:         KeyF8,
	termbox.KeyF9:         KeyF9,
	termbox.KeyF10:        KeyF10,
	termbox.KeyF11:        KeyF11,
	termbox.KeyF12:        KeyF12,
	termbox.KeyInsert:     KeyInsert,
	termbox.KeyDelete:     KeyDelete,
	termbox.KeyHome:       KeyHome,
	termbox.KeyEnd:        KeyEnd,
	termbox.KeyPgup:       KeyPgup,
	termbox.KeyPgdn:       KeyPgdn,
	termbox.KeyArrowUp:    KeyArrowUp,
	termbox.KeyArrowDown:  KeyArrowDown,
	termbox.KeyArrowLeft:  KeyArrowLeft,
	termbox.KeyArrowRight: KeyArrowRight,
	termbox.KeyBackspace2: KeyBackspace2,
	termbox.KeyCtrlA:      KeyCtrlA,
	termbox.KeyCtrlB:      KeyCtrlB,
	termbox.KeyCtrlC:      KeyCtrlC,
	termbox.KeyCtrlD:      KeyCtrlD,
	termbox.KeyCtrlE:      KeyCtrlE,
	termbox.KeyCtrlF:      KeyCtrlF,
	termbox.KeyCtrlG:      KeyCtrlG,
	termbox.KeyCtrlH:      KeyCtrlH,
	termbox.KeyCtrlI:      KeyCtrlI,
	termbox.KeyCtrlJ:      KeyCtrlJ,
	termbox.KeyCtrlK:      KeyCtrlK,
	termbox.KeyCtrlL:      KeyCtrlL,
	termbox.KeyCtrlM:      KeyCtrlM,
	termbox.KeyCtrlN:      KeyCtrlN,
	termbox.KeyCtrlO:      KeyCtrlO,
	termbox.KeyCtrlP:      KeyCtrlP,
	termbox.KeyCtrlQ:      KeyCtrlQ,
	termbox.KeyCtrlR:      KeyCtrlR,
	termbox.KeyCtrlS:      KeyCtrlS,
	termbox.KeyCtrlT:      KeyCtrlT,
	termbox.KeyCtrlU:      KeyCtrlU,
	termbox.KeyCtrlV:      KeyCtrlV,
	termbox.KeyCtrlW:      KeyCtrlW,
	termbox.KeyCtrlX:      KeyCtrlX,
	termbox.KeyCtrlY:      KeyCtrlY,
	termbox.KeyCtrlZ:      KeyCtrlZ,
	termbox.KeyCtrl3:      KeyCtrl3,
	termbox.KeyCtrl4:      KeyCtrl4,
	termbox.KeyCtrl5:      KeyCtrl5,
	termbox.KeyCtrl6:      KeyCtrl6,
	termbox.KeyCtrl7:      KeyCtrl7,
}
