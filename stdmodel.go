package main

// Cell
type Cell struct {
    Value   string
}

// Standard model
type StdModel struct {
    Cells       [][]Cell
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
}

// Sets the cell value
func (sm *StdModel) SetCellValue(r, c int, value string) {
    rs, cs := sm.Dimensions()
    if (r >= 0) && (c >= 0) && (r < rs) && (c < cs) {
        sm.Cells[r][c].Value = value
    }
}