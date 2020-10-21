package main

import (
	"regexp"

	"github.com/lmika/ted/ui"
)

// The session is responsible for managing the UI and the model and handling
// the interaction between the two and the user.
type Session struct {
	model           Model
	Source          ModelSource
	Frame           *Frame
	Commands        *CommandMapping
	UIManager       *ui.Ui
	modelController *ModelViewCtrl
	pasteBoard		RWModel

	LastSearch *regexp.Regexp
}

func NewSession(uiManager *ui.Ui, frame *Frame, source ModelSource) *Session {
	model := NewSingleCellStdModel()
	session := &Session{
		model:           model,
		Source:          source,
		Frame:           frame,
		Commands:        NewCommandMapping(),
		UIManager:       uiManager,
		modelController: NewGridViewModel(model),
		pasteBoard: 	 NewSingleCellStdModel(),
	}

	frame.SetModel(&SessionGridModel{session.modelController})

	session.Commands.RegisterViewCommands()
	session.Commands.RegisterViewKeyBindings()

	// Also assign this session with the frame
	frame.Session = session

	return session
}

// LoadFromSource loads the model from the source, replacing the existing model
func (session *Session) LoadFromSource() {
	newModel, err := session.Source.Read()
	if err != nil {
		session.Frame.Message(err.Error())
		return
	}

	session.model = newModel
	session.modelController.SetModel(newModel)
}

// Input from the frame
func (session *Session) KeyPressed(key rune, mod int) {
	// Add the mod key modifier
	if mod&ui.ModKeyAlt != 0 {
		key |= ModAlt
	}

	cmd := session.Commands.KeyMapping(key)
	if cmd != nil {
		err := cmd.Do(&CommandContext{session, nil})
		if err != nil {
			session.Frame.ShowMessage(err.Error())
		}
	}
}

// The command context used by the session
type CommandContext struct {
	session *Session
	args    []string
}

func (scc *CommandContext) WithArgs(args []string) *CommandContext {
	return &CommandContext{
		session: scc.session,
		args:    args,
	}
}

func (scc *CommandContext) Args() []string {
	return scc.args
}

func (scc *CommandContext) ModelVC() *ModelViewCtrl {
	return scc.session.modelController
}

func (scc *CommandContext) Session() *Session {
	return scc.session
}

func (scc *CommandContext) Frame() *Frame {
	return scc.session.Frame
}

// Error displays an error if err is not nil
func (scc *CommandContext) Error(err error) {
	scc.Frame().Error(err)
}

// Session grid model
type SessionGridModel struct {
	GridViewModel *ModelViewCtrl
}

// Returns the size of the grid model (width x height)
func (sgm *SessionGridModel) Dimensions() (int, int) {
	rs, cs := sgm.GridViewModel.Model().Dimensions()
	return cs, rs
}

// Returns the size of the particular column.  If the size is 0, this indicates that the column is hidden.
func (sgm *SessionGridModel) ColWidth(col int) int {
	return sgm.GridViewModel.ColAttrs(col).Size
}

// Returns the size of the particular row.  If the size is 0, this indicates that the row is hidden.
func (sgm *SessionGridModel) RowHeight(row int) int {
	return sgm.GridViewModel.RowAttrs(row).Size
}

// Returns the value of the cell a position X, Y
func (sgm *SessionGridModel) CellValue(x int, y int) string {
	return sgm.GridViewModel.Model().CellValue(y, x)
}

func (sgm *SessionGridModel) CellAttributes(x int, y int) (fg, bg ui.Attribute) {
	rowAttrs := sgm.GridViewModel.RowAttrs(y)
	colAttrs := sgm.GridViewModel.ColAttrs(y)

	if rowAttrs.Marker != MarkerNone {
		return markerAttributes[rowAttrs.Marker], 0
	} else if colAttrs.Marker != MarkerNone {
		return markerAttributes[colAttrs.Marker], 0
	}
	return 0, 0
}

var markerAttributes = map[Marker]ui.Attribute{
	MarkerRed:   ui.ColorRed,
	MarkerGreen: ui.ColorGreen,
	MarkerBlue:  ui.ColorBlue,
}
