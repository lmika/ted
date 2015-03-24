package main

import "./ui"

// The session is responsible for managing the UI and the model and handling
// the interaction between the two and the user.
type Session struct {
    Model           Model
    Frame           *Frame
    Commands        *CommandMapping
}

func NewSession(frame *Frame, model Model) *Session {
    session := &Session{
        Model: model,
        Frame: frame,
        Commands: NewCommandMapping(),
    }

    frame.SetModel(&SessionGridModel{session})

    session.Commands.RegisterViewCommands()
    session.Commands.RegisterViewKeyBindings()

    // Also assign this session with the frame
    frame.Session = session

    return session
}

// Input from the frame
func (session *Session) KeyPressed(key rune, mod int) {
    // Add the mod key modifier
    if (mod & ui.ModKeyAlt != 0) {
        key |= ModAlt
    }

    cmd := session.Commands.KeyMapping(key)
    if cmd != nil {
        err := cmd.Do(SessionCommandContext{session})
        if err != nil {
            session.Frame.ShowMessage(err.Error())
        }
    }
}


// The command context used by the session
type SessionCommandContext struct {
    Session     *Session
}

func (scc SessionCommandContext) Frame() *Frame {
    return scc.Session.Frame
}



// Session grid model
type SessionGridModel struct {
    Session     *Session
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
