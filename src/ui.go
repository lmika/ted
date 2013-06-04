/**
 * UI package.
 */
package main

import "github.com/nsf/termbox-go"


// ==========================================================================
// UI event.

/**
 * The types of events.
 */
type    EventType   int
const (
    EventKeyPress   EventType   =   iota
)

/**
 * An event callback
 */
type UiEvent struct {
    eventType       EventType
    eventPar        int   
}

// ==========================================================================
// UI component.

type UiComponent interface {
    /**
     * Request to redraw this component.
     */
    Redraw(x int, y int, w int, h int)

    /**
     * Request the minimum dimensions of the component (width, height).  If
     * either dimension is -1, no minimum is imposed.
     */
    RequestDims() (int, int)
}


// ==========================================================================
// UI context.

/**
 * Ui context type.
 */
type Ui struct {
    statusBar       UiComponent
}


/**
 * Creates a new UI context.  This also initializes the UI state.
 * Returns the context and an error.
 */
func NewUI() (*Ui, error) {
    termboxError := termbox.Init()

    if termboxError != nil {
        return nil, termboxError
    } else {
        uiCtx := &Ui{&UiStatusBar{"Hello", "World"}}
        return uiCtx, nil
    }
}


/**
 * Closes the UI context.
 */
func (ui *Ui) Close() {
    termbox.Close()
}


/**
 * Redraws the UI.
 */
func (ui *Ui) Redraw() {
    var width, height int = termbox.Size()
    ui.redrawInternal(width, height)
}

/**
 * Internal redraw function which does not query the terminal size.
 */
func (ui *Ui) redrawInternal(width, height int) {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    // TODO: This will eventually offload to UI "components"

    // Draws the status bar
    ui.statusBar.Redraw(0, height - 2, width, 2)

    termbox.Flush()
}


/**
 * Waits for a UI event.  Returns the event (if it's relevant to the user).
 */
func (ui *Ui) NextEvent() UiEvent {
    for {
        event := termbox.PollEvent()
        if event.Type == termbox.EventResize {
            ui.redrawInternal(event.Width, event.Height)
        } else {
            return UiEvent{EventKeyPress, 0}
        }
    }
}


// ==========================================================================
// Status bar component

type UiStatusBar struct {
    left    string          // Left aligned string
    right   string          // Right aligned string
}

// Minimum dimensions
func (sbar *UiStatusBar) RequestDims() (int, int) {
    return -1, 2
}

// Status bar redraw
func (sbar *UiStatusBar) Redraw(x int, y int, w int, h int) {
    leftLen := len(sbar.left)
    rightLen := len(sbar.right)
    rightPos := w - rightLen
    
    for x1 := 0; x1 < w; x1++ {
        var runeToPrint rune = ' '
        if x1 < leftLen {
            runeToPrint = rune(sbar.left[x1])
        } else if x1 >= rightPos {
            runeToPrint = rune(sbar.right[x1 - rightPos])
        }
        termbox.SetCell(x1, y, runeToPrint, termbox.AttrReverse, termbox.AttrReverse)
    }
}
