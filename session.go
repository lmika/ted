package main

import "./ui"

// The session is responsible for managing the UI and the model and handling
// the interaction between the two and the user.
type Session struct {
    Frame           *Frame
    Commands        *CommandMapping
}

func NewSession(frame *Frame) *Session {
    session := &Session{
        Frame: frame,
        Commands: NewCommandMapping(),
    }

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