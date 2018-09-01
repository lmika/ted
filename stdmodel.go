package main

// Cell
type Cell struct {
	Value string
}

// Standard model
type StdModel struct {
	Cells [][]Cell
	dirty bool
}

//
func NewSingleCellStdModel() *StdModel {
	sm := new(StdModel)
	sm.appendStr([]string { "" })
	sm.dirty = false

	return sm
}

/**
 * The dimensions of the model (height, width).
 */
func (sm *StdModel) Dimensions() (int, int) {
	if len(sm.Cells) == 0 {
		return 0, 0
	} else {
		return len(sm.Cells), len(sm.Cells[0])
	}
}

/**
 * Returns the value of a cell.
 */
func (sm *StdModel) CellValue(r, c int) string {
	rs, cs := sm.Dimensions()
	if (r >= 0) && (c >= 0) && (r < rs) && (c < cs) {
		return sm.Cells[r][c].Value
	} else {
		return ""
	}
}

// Resize the model.
func (sm *StdModel) Resize(rs, cs int) {
	oldRowCount := len(sm.Cells)

	newRows := make([][]Cell, rs)
	for r := range newRows {
		newCols := make([]Cell, cs)
		if r < oldRowCount {
			copy(newCols, sm.Cells[r])
		}
		newRows[r] = newCols
	}

	sm.Cells = newRows
	sm.dirty = true
}

// Sets the cell value
func (sm *StdModel) SetCellValue(r, c int, value string) {
	rs, cs := sm.Dimensions()
	if (r >= 0) && (c >= 0) && (r < rs) && (c < cs) {
		sm.Cells[r][c].Value = value
	}
	sm.dirty = true
}

// appendStr appends the model with the given row
func (sm *StdModel) appendStr(row []string) {
	if len(sm.Cells) == 0 {
		cells := sm.strSliceToCell(row, len(row))
		sm.Cells = [][]Cell{ cells }
		return
	}

	cols := len(sm.Cells[0])
	if len(row) > cols {
		sm.Resize(len(sm.Cells), len(row))
		cols = len(sm.Cells[0])
	}

	cells := sm.strSliceToCell(row, cols)
	sm.Cells = append(sm.Cells, cells)
}

func (sm *StdModel) strSliceToCell(row []string, targetRowLen int) []Cell {
	cs := make([]Cell, targetRowLen)
	for i := 0; i < targetRowLen; i++ {
		if i < len(row) {
			cs[i].Value = row[i]
		}
	}
	return cs
}

func (sm *StdModel) IsDirty() bool {
	return sm.dirty
}
