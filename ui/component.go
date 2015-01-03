// Components of the UI and various event interfaces.

package ui


// The set of event types supported by the UI package.




// An interface of a UI component.
type UiComponent interface {

    // Request from the manager for the component to draw itself.  This is given a drawable context.
    Redraw(context *DrawContext)

    // Called to remeasure the size of the component.  Provided with the maximum dimensions of the component
    // and expected to provide the minimum component size.  When called to redraw, the component will be
    // provided with AT LEAST the minimum dimensions returned by this method.
    Remeasure(w, h int) (int, int)
}


// ==========================================================================
// UI context.




// ==========================================================================
// Status bar component

/*
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
*/