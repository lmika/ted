// Standard layout components.

package ui

// Instance of a component.
type componentInstance struct {
    component       UiComponent
    x, y            int         // X and Y offset of this component
    height          int         // Height allocated to the component
    width           int         // Width allocated to the component
}

// A layout component.
type LinearLayout struct {
    // The set of components that this layout component contains
    subcomponents    []*componentInstance
}

// Adds a component to the end of the component list managed by this layout.
func (l *LinearLayout) Append(component UiComponent) {
    l.subcomponents = append(l.subcomponents, &componentInstance{component, 0, 0, 0, 0})
}


// A vertical layout component.
type VertLinearLayout struct {
    LinearLayout
    maxHeight   int
}


// Request from the manager for the component to draw itself.  This is given a drawable context.
func (vl *VertLinearLayout) Redraw(context *DrawContext) {
    for _, ci := range vl.LinearLayout.subcomponents {
        subContext := context.NewSubContext(ci.x, ci.y, ci.width, ci.height)
        ci.component.Redraw(subContext)
    }
}

// Remeasures the components currently managed within this layout.
func (vl *VertLinearLayout) Remeasure(w, h int) (int, int) {
    posy := 0

    // TODO: At the moment, this simply takes the minimum and maximum dimensions requested
    // by the components.  This needs to be extended to allow for "special" components which
    // take up dynamic space.
    for _, ci := range vl.LinearLayout.subcomponents {
        _, rh := ci.component.Remeasure(w, h - posy)

        // All components within this layer are given the full width of the screen.
        ci.x = 0
        ci.width = w
        ci.y = posy
        ci.height = rh

        // Put the next component directly below the previous one
        posy += rh
    }

    return w, posy
}


// A relative layout component.  This has a "client" component, bordered by a north,
// south, east and west component.  The N,S,E,W components will be provided with the full dimensions
// whereas the client component will be provided with the remaining size.  Each one of the components
// can be null.
type RelativeLayout struct {
    North, South, East, West, Client    UiComponent

    // Measured client borders
    ct, cb, cl, cr, ch, cw  int

    // North/south heights and east/west widths
    nh, sh                  int
    ew, ww                  int
}

func (rl *RelativeLayout) Remeasure(w, h int) (int, int) {
    if rl.North != nil {
        _, rl.nh = rl.North.Remeasure(w, h)
        rl.ct = rl.nh
    } else {
        rl.ct = 0
    }

    if rl.South != nil {
        _, rl.sh = rl.South.Remeasure(w, h - rl.nh)
        rl.cb = h - rl.sh
    } else {
        rl.cb = h
    }

    // TODO: East and west
    rl.cl = 0
    rl.cr = w

    rl.ch = h - rl.nh - rl.sh
    rl.cw = w

    return w, h
}

func (vl *RelativeLayout) Redraw(context *DrawContext) {
    if vl.North != nil {
        vl.North.Redraw(context.NewSubContext(0, 0, context.W, vl.nh))
    }
    if vl.South != nil {
        vl.South.Redraw(context.NewSubContext(0, vl.cb, context.W, vl.sh))
    }
    if vl.Client != nil {
        vl.Client.Redraw(context.NewSubContext(vl.cl, vl.ct, vl.cw, vl.ch))
    }
}