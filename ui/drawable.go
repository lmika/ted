// Provides access to the drawable primitives.

package ui

// A drawable context.  Each context is allocated part of the screen and provides methods for
// drawing within it's context.
type DrawContext struct {

    // The left and top position of the context.
    X, Y        int

    // The width and height position of the context.
    W, H        int

    // The current driver
    driver      Driver

    // The current foregound and background attributes
    fa, ba      Attribute
}


// Returns a new subcontext.  The sub-context must be an area within the current context.
func (dc *DrawContext) NewSubContext(offsetX, offsetY, width, height int) *DrawContext {
    return &DrawContext{
        X: intMax(dc.X + offsetX, dc.X),
        Y: intMax(dc.Y + offsetY, dc.Y),
        W: intMax(width, 0),
        H: intMax(height, 0),

        driver: dc.driver,
    }
}

// Sets the foreground attribute
func (dc *DrawContext) SetFgAttr(attr Attribute) {
    dc.fa = attr
}

// Sets the background attribute
func (dc *DrawContext) SetBgAttr(attr Attribute) {
    dc.ba = attr
}


// Draws a horizontal rule with a specific rune.
func (dc *DrawContext) HorizRule(y int, ch rune) {
    for x := 0; x < dc.W; x++ {
        dc.DrawRune(x, y, ch)
    }
}

// Prints a string at a specific offset.  This will be bounded by the size of the drawing context.
func (dc *DrawContext) Print(x, y int, str string) {
    for _, ch := range str {
        dc.DrawRune(x, y, ch)
        x++
    }
}

// Prints a right-justified string at a specific offset.  This will be bounded by the size of the drawing context.
func (dc *DrawContext) PrintRight(x, y int, str string) {
    l := len(str)
    dc.Print(x - l, y, str)
}

// Draws a rune at a local point X, Y with the current foreground and background attributes
func (dc *DrawContext) DrawRune(x, y int, ch rune) {
    if rx, ry, isWithinContext := dc.localPointToRealPoint(x, y); isWithinContext {
        dc.driver.SetCell(rx, ry, ch, dc.fa, dc.ba)
    }
}

// Draws a rune at a local point X, Y with specific foreground and background attributes
func (dc *DrawContext) DrawRuneWithAttrs(x, y int, ch rune, fa, ba Attribute) {
    if rx, ry, isWithinContext := dc.localPointToRealPoint(x, y); isWithinContext {
        dc.driver.SetCell(rx, ry, ch, fa, ba)
    }
}

// Converts a local point to a real point.  If the point is within the context, also returns true
func (dc *DrawContext) localPointToRealPoint(x, y int) (int, int, bool) {
    rx, ry := x + dc.X, y + dc.Y
    return rx, ry, (rx >= dc.X) && (ry >= dc.Y) && (rx < dc.X + dc.W) && (ry < dc.Y + dc.H)
}