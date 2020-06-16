/**
 * The model.
 */
package main

// An abstract model interface.  At a minimum, models must be read only.
type Model interface {

	// The dimensions of the model (height, width).
	Dimensions() (int, int)

	// Returns the value of a cell
	CellValue(r, c int) string
}

// A read/write model.
type RWModel interface {
	Model

	// Resize the model.
	Resize(newRow, newCol int)

	// Sets the cell value
	SetCellValue(r, c int, value string)

	// Returns true if the model has been modified in some way
	IsDirty() bool
}


