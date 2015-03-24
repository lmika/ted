package main

import (
    "./ui"
)

type Mode   int

const (
    // The grid is selectable
    GridMode    Mode    =   iota
)

// A frame is a UI instance.
type Frame struct {
    Session             *Session

    uiManager           *ui.Ui
    clientArea          *ui.RelativeLayout
    grid                *ui.Grid
    messageView         *ui.TextView
    textEntry           *ui.TextEntry
    statusBar           *ui.StatusBar
    textEntrySwitch     *ui.ProxyLayout
}

// Creates the UI and returns a new frame
func NewFrame(uiManager *ui.Ui) *Frame {
    frame := &Frame{
        uiManager: uiManager,
    }

    frame.grid = ui.NewGrid(nil)
    frame.messageView = &ui.TextView{"Hello"}
    frame.statusBar = &ui.StatusBar{"Test", "Status"}
    frame.textEntrySwitch = &ui.ProxyLayout{frame.messageView}
    frame.textEntry = &ui.TextEntry{}

    // Build the UI frame
    statusLayout := &ui.VertLinearLayout{}
    statusLayout.Append(frame.statusBar)
    statusLayout.Append(frame.textEntrySwitch)

    frame.clientArea = &ui.RelativeLayout{ Client: frame.grid, South: statusLayout }
    return frame
}

// Returns the root component of the frame
func (frame *Frame) RootComponent() ui.UiComponent {
    return frame.clientArea
}

// Sets the current model of the frame
func (frame *Frame) SetModel(model ui.GridModel) {
    frame.grid.SetModel(model)
}

// Returns the grid component
func (frame *Frame) Grid() *ui.Grid {
    return frame.grid
}

// Sets the specific mode.
func (frame *Frame) EnterMode(mode Mode) {
    switch mode {
    case GridMode:
        frame.uiManager.SetFocusedComponent(frame)
    }
}

// Show a message.  This will switch the bottom to the messageView and select the frame
func (frame *Frame) ShowMessage(msg string) {
    frame.messageView.Text = msg
    frame.textEntrySwitch.Component = frame.messageView
    //frame.EnterMode(GridMode)
}

// Shows the value of the currently select grid cell
func (frame *Frame) ShowCellValue() {
    displayValue := frame.grid.CurrentCellDisplayValue()
    frame.ShowMessage(displayValue)
}

// Handle the main grid input as this is the "component" that handles command input.
func (frame *Frame) KeyPressed(key rune, mod int) {
    if frame.Session != nil {
        frame.Session.KeyPressed(key, mod)
    }
}