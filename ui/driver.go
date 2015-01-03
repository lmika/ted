// The UI driver.  This is used to interact with the terminal drawing routines.

package ui

// The set of attributes a specific cell can have
type Attribute  uint16

const (
    // Can have only one of these
    ColorDefault Attribute = iota
    ColorBlack
    ColorRed
    ColorGreen
    ColorYellow
    ColorBlue
    ColorMagenta
    ColorCyan
    ColorWhite
)

// and zero or more of these (combined using OR '|')
const (
    AttrBold Attribute = 1 << (iota + 9)
    AttrUnderline
    AttrReverse
)

// The type of events supported by the driver
type    EventType   int

const (
    EventNone       EventType   =   iota

    // Event when the window is resized
    EventResize
    
    // Event indicating a key press.  The event parameter is the key scancode?
    EventKeyPress
)

// Data from an event callback.
type Event struct {
    Type            EventType
    Par             int
}


// The terminal driver interface.
type Driver interface {

    // Initializes the driver.  Returns an error if there was an error
    Init() error

    // Closes the driver
    Close()

    // Returns the size of the window.
    Size() (int, int)

    // Sets the value of a specific cell
    SetCell(x, y int, ch rune, fg, bg Attribute)

    // Synchronizes the internal buffer with the real buffer
    Sync()

    // Wait for an event
    WaitForEvent() Event
}

// 