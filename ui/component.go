// Components of the UI and various event interfaces.

package ui


// An interface of a UI component.
type UiComponent interface {

    // Request from the manager for the component to draw itself.  This is given a drawable context.
    Redraw(context *DrawContext)

    // Called to remeasure the size of the component.  Provided with the maximum dimensions of the component
    // and expected to provide the minimum component size.  When called to redraw, the component will be
    // provided with AT LEAST the minimum dimensions returned by this method.
    Remeasure(w, h int) (int, int)
}

// A component implementing this interface can be given focus and will receive keyboard events.
type FocusableComponent interface {
    
    // Called when the component has focus and a key has been pressed
    KeyPressed(key rune, mod int)
}
