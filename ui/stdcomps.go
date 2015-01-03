// Standard components

package ui


// Status bar component.  This component displays text on the left and right of it's
// allocated space.
type StatusBar struct {
    Left    string          // Left aligned string
    Right   string          // Right aligned string
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
    /*
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
    */
}