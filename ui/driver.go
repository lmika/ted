// The UI driver.  This is used to interact with the terminal drawing routines.

package ui

// The set of attributes a specific cell can have
type Attribute uint16

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

// Special keys
const (
	KeyCtrlSpace rune = 0x8000 + iota
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
	KeyCtrl3
	KeyCtrl4
	KeyCtrl5
	KeyCtrl6
	KeyCtrl7
	KeyCtrl8

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight

	KeyBackspace = KeyCtrlH
	KeyBackspace2
	KeyEnter = KeyCtrlM
	KeyEsc   = KeyCtrl3
)

// The type of events supported by the driver
type EventType int

const (
	EventNone EventType = iota

	// Event when the window is resized
	EventResize

	// Event indicating a key press.  The key is set in Ch and modifications
	// are set in Or
	EventKeyPress
)

const (
	ModKeyAlt int = (1 << iota)
)

// Data from an event callback.
type Event struct {
	Type EventType
	Par  int
	Ch   rune
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

	// Move the position of the cursor
	SetCursor(x, y int)

	// Hide the cursor
	HideCursor()
}
