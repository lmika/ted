package main

import (
	"bitbucket.org/lmika/ted-v2/ui"
)

type Mode int

const (
	NilMode Mode = iota

	// The grid is selectable
	GridMode

	// EntryMode is when the text entry is selected
	EntryMode
)

// A frame is a UI instance.
type Frame struct {
	Session *Session

	mode Mode

	uiManager       *ui.Ui
	clientArea      *ui.RelativeLayout
	grid            *ui.Grid
	messageView     *ui.TextView
	textEntry       *ui.TextEntry
	statusBar       *ui.StatusBar
	textEntrySwitch *ui.ProxyLayout
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

	frame.clientArea = &ui.RelativeLayout{Client: frame.grid, South: statusLayout}
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

// Enter the specific mode.
func (frame *Frame) enterMode(mode Mode) {
	switch mode {
	case GridMode:
		frame.uiManager.SetFocusedComponent(frame)
	case EntryMode:
		frame.textEntrySwitch.Component = frame.textEntry
		frame.uiManager.SetFocusedComponent(frame.textEntry)
	}
}

// Exit the specific mode.
func (frame *Frame) exitMode(mode Mode) {
	switch mode {
	case EntryMode:
		frame.textEntrySwitch.Component = frame.messageView
	}
}

func (frame *Frame) setMode(mode Mode) {
	frame.exitMode(frame.mode)
	frame.mode = mode
	frame.enterMode(frame.mode)
}

// Message sets the message view's message
func (frame *Frame) Message(s string) {
	frame.messageView.Text = s
}

// Prompt the user for input.  This switches the mode to entry mode.
func (frame *Frame) Prompt(prompt string, callback func(res string)) {
	frame.textEntry.Prompt = prompt
	frame.textEntry.SetValue("")

	frame.textEntry.OnCancel = frame.exitEntryMode
	frame.textEntry.OnEntry = func(res string) {
		frame.exitEntryMode()
		callback(res)
	}

	frame.setMode(EntryMode)
}

func (frame *Frame) exitEntryMode() {
	frame.textEntry.OnEntry = nil
	frame.setMode(GridMode)
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
