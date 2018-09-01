package main

import "bitbucket.org/lmika/ted-v2/ui"

// The session is responsible for managing the UI and the model and handling
// the interaction between the two and the user.
type Session struct {
	Model     Model
	Source    ModelSource
	Frame     *Frame
	Commands  *CommandMapping
	UIManager *ui.Ui
}

func NewSession(uiManager *ui.Ui, frame *Frame, source ModelSource) *Session {
	session := &Session{
		Model:     nil,
		Source:	   source,
		Frame:     frame,
		Commands:  NewCommandMapping(),
		UIManager: uiManager,
	}

	frame.SetModel(&SessionGridModel{session})

	session.Commands.RegisterViewCommands()
	session.Commands.RegisterViewKeyBindings()

	// Also assign this session with the frame
	frame.Session = session

	return session
}

// LoadFromSource loads the model from the source, replacing the existing model
func (session *Session) LoadFromSource() error {
	newModel, err := session.Source.Read()
	if err != nil {
		return err
	}

	session.Model = newModel
	return nil
}

// Input from the frame
func (session *Session) KeyPressed(key rune, mod int) {
	// Add the mod key modifier
	if mod&ui.ModKeyAlt != 0 {
		key |= ModAlt
	}

	cmd := session.Commands.KeyMapping(key)
	if cmd != nil {
		err := cmd.Do(&CommandContext{session})
		if err != nil {
			session.Frame.ShowMessage(err.Error())
		}
	}
}

// The command context used by the session
type CommandContext struct {
	session *Session
}

func (scc *CommandContext) Session() *Session {
	return scc.session
}

func (scc *CommandContext) Frame() *Frame {
	return scc.session.Frame
}

// Error displays an error if err is not nil
func (scc *CommandContext) ShowError(err error) {
	if err != nil {
		scc.Frame().Message(err.Error())
	}
}

// Session grid model
type SessionGridModel struct {
	Session *Session
}

// Returns the size of the grid model (width x height)
func (sgm *SessionGridModel) Dimensions() (int, int) {
	rs, cs := sgm.Session.Model.Dimensions()
	return cs, rs
}

// Returns the size of the particular column.  If the size is 0, this indicates that the column is hidden.
func (sgm *SessionGridModel) ColWidth(int) int {
	return 24
}

// Returns the size of the particular row.  If the size is 0, this indicates that the row is hidden.
func (sgm *SessionGridModel) RowHeight(int) int {
	return 1
}

// Returns the value of the cell a position X, Y
func (sgm *SessionGridModel) CellValue(x int, y int) string {
	return sgm.Session.Model.CellValue(y, x)
}
