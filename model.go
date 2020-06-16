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

// Deletes a row of a model
func DeleteRow(model RWModel, row int) {
	h, w := model.Dimensions()
	for r := row; r < h-1; r++ {
		for c := 0; c < w; c++ {
			model.SetCellValue(r, c, model.CellValue(r+1, c))
		}
	}

	model.Resize(h-1, w)
}

// Deletes a column of a model
func DeleteCol(model RWModel, col int) {
	h, w := model.Dimensions()
	for c := col; c < w-1; c++ {
		for r := 0; r < h; r++ {
			model.SetCellValue(r, c, model.CellValue(r, c+1))
		}
	}
	model.Resize(h, w-1)
}
