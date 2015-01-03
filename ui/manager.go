// The UI manager.  This manages the components that make up the UI and dispatches
// events.

package ui


// The UI manager
type Ui struct {
    //grid            *Grid
    //statusBar       *UiStatusBar

    // The root component
    rootComponent       UiComponent

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

    /*
    termboxError := termbox.Init()

    if termboxError != nil {
        return nil, termboxError
    } else {
        uiCtx := new(Ui)  // &Ui{&UiStatusBar{"Hello", "World"}}
        uiCtx.grid = NewGrid(&TestModel{})
        uiCtx.statusBar = &UiStatusBar{"Hello", "World"}
        return uiCtx, nil
    }
    
    // XXX: Workaround for bug in compiler
    panic("Unreachable code")
    return nil, nil
    */
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

/**
 * Internal redraw function which does not query the terminal size.
 */
 /*
func (ui *Ui) redrawInternal(width, height int) {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    // TODO: This will eventually offload to UI "components"
    ui.grid.Redraw(0, 0, width, height - 2)

    // Draws the status bar
    ui.statusBar.Redraw(0, height - 2, width, 2)

    termbox.Flush()
}
*/


// Enter the UI loop
func (ui *Ui) Loop() {
    //for {
        ui.Redraw()
        ui.driver.WaitForEvent()

        /*
        if event.Type == termbox.EventResize {
            ui.redrawInternal(event.Width, event.Height)
        } else {

            // !!TEMP!!
            if (event.Ch == 'i') {
                ui.grid.MoveBy(0, -1)
            } else if (event.Ch == 'k') {
                ui.grid.MoveBy(0, 1)
            } else if (event.Ch == 'j') {
                ui.grid.MoveBy(-1, 0)
            } else if (event.Ch == 'l') {
                ui.grid.MoveBy(1, 0)
            } else {
                return UiEvent{EventKeyPress, 0}
            }
            // !!END TEMP!!

            ui.Redraw()
            //return UiEvent{EventKeyPress, 0}
        }
        */
    //}
    
    // XXX: Workaround for bug in compiler
}