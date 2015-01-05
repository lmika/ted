// The UI manager.  This manages the components that make up the UI and dispatches
// events.

package ui


// The UI manager
type Ui struct {
    // The root component
    rootComponent       UiComponent
    focusedComponent    FocusableComponent

    drawContext         *DrawContext
    driver              Driver
}


// Creates a new UI context.  This also initializes the UI state.
// Returns the context and an error.
func NewUI() (*Ui, error) {
    driver := &TermboxDriver{}
    err := driver.Init()

    if err != nil {
        return nil, err
    }

    drawContext := &DrawContext{ driver: driver }
    ui := &Ui{ drawContext: drawContext,  driver: driver }

    return ui, nil
}


// Closes the UI context.
func (ui *Ui) Close() {
    ui.driver.Close()
}

// Sets the root component
func (ui *Ui) SetRootComponent(comp UiComponent) {
    ui.rootComponent = comp
    ui.Remeasure()
}

// Sets the focused component
func (ui *Ui) SetFocusedComponent(newFocused FocusableComponent) {
    ui.focusedComponent = newFocused
}

// Remeasures the UI
func (ui *Ui) Remeasure() {
    ui.drawContext.X = 0
    ui.drawContext.Y = 0
    ui.drawContext.W, ui.drawContext.H = ui.driver.Size()

    ui.rootComponent.Remeasure(ui.drawContext.W, ui.drawContext.H)
}

// Redraws the UI.
func (ui *Ui) Redraw() {
    ui.Remeasure()

    ui.rootComponent.Redraw(ui.drawContext)
    ui.driver.Sync()
}


// Enter the UI loop
func (ui *Ui) Loop() {
    for {
        ui.Redraw()
        event := ui.driver.WaitForEvent()

        // TODO: If the event is a key-press, do something.
        if event.Type == EventKeyPress {
            if ui.focusedComponent != nil {
                ui.focusedComponent.KeyPressed(event.Ch, event.Par)
            }
        } else if event.Type == EventResize {

            // HACK: Find another way to refresh the size of the screen to prevent a full redraw.
            ui.driver.Sync()
        }
    }
}