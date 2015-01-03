// The terminal-box driver

package ui

import (
    "github.com/nsf/termbox-go"
)


type TermboxDriver struct {
}


// Initializes the driver.  Returns an error if there was an error
func (td *TermboxDriver) Init() error {
    return termbox.Init()
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

// Synchronizes the internal buffer with the real buffer
func (td *TermboxDriver) Sync() {
    termbox.Sync()
}

// Wait for an event
func (td *TermboxDriver) WaitForEvent() Event {
    tev := termbox.PollEvent()

    switch tev.Type {
    case termbox.EventResize:
        return Event{EventResize, 0}
    case termbox.EventKey:
        return Event{EventKeyPress, 0}
    default:
        return Event{EventNone, 0}
    }
}