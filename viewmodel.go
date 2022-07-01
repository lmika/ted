package main

import (
	"errors"
)

type ModelViewCtrl struct {
	model    Model
	rowAttrs []SliceAttr
	colAttrs []SliceAttr
}

func NewGridViewModel(model Model) *ModelViewCtrl {
	gvm := &ModelViewCtrl{}
	gvm.SetModel(model)
	return gvm
}

func (gvm *ModelViewCtrl) Model() Model {
	return gvm.model
}

func (gvm *ModelViewCtrl) SetModel(m Model) {
	gvm.model = m
	gvm.modelWasResized()
}

func (gvm *ModelViewCtrl) RowAttrs(row int) SliceAttr {
	if row < len(gvm.rowAttrs) {
		return gvm.rowAttrs[row]
	}
	return DefaultRowAttrs
}

func (gvm *ModelViewCtrl) ColAttrs(col int) SliceAttr {
	if col < len(gvm.colAttrs) {
		return gvm.colAttrs[col]
	}
	return DefaultColAttrs
}

func (gvm *ModelViewCtrl) SetRowAttrs(row int, newAttrs SliceAttr) {
	gvm.rowAttrs[row] = newAttrs
}

func (gvm *ModelViewCtrl) SetColAttrs(col int, newAttrs SliceAttr) {
	if col >= 0 && col < len(gvm.colAttrs) {
		gvm.colAttrs[col] = newAttrs
	}
}

func (gvm *ModelViewCtrl) SetCellValue(r, c int, newValue string) error {
	rwModel, isRWModel := gvm.model.(RWModel)
	if !isRWModel {
		return ErrModelReadOnly
	}

	rwModel.SetCellValue(r, c, newValue)
	return nil
}

func (gvm *ModelViewCtrl) Resize(newRow, newCol int) error {
	rwModel, isRWModel := gvm.model.(RWModel)
	if !isRWModel {
		return ErrModelReadOnly
	}

	rwModel.Resize(newRow, newCol)
	gvm.modelWasResized()

	return nil
}

func (gvm *ModelViewCtrl) OpenRight(col int) error {
	if col < 0 {
		return errors.New("col out of bound")
	}
	return gvm.insertColumn(col + 1)
}

func (gvm *ModelViewCtrl) insertColumn(col int) error {
	rwModel, isRWModel := gvm.model.(RWModel)
	if !isRWModel {
		return ErrModelReadOnly
	}

	dr, dc := rwModel.Dimensions()
	if col < 0 || col > dc {
		return errors.New("col out of bound")
	}

	rwModel.Resize(dr, dc+1)

	for c := dc; c >= col; c-- {
		for r := 0; r < dr; r++ {
			if c == col {
				rwModel.SetCellValue(r, c, "")
			} else {
				rwModel.SetCellValue(r, c, rwModel.CellValue(r, c-1))
			}
		}
	}

	return nil
}

// Deletes a row of a model
func (gvm *ModelViewCtrl) DeleteRow(row int) error {
	rwModel, isRWModel := gvm.model.(RWModel)
	if !isRWModel {
		return ErrModelReadOnly
	}

	h, w := rwModel.Dimensions()
	for r := row; r < h-1; r++ {
		for c := 0; c < w; c++ {
			rwModel.SetCellValue(r, c, rwModel.CellValue(r+1, c))
			gvm.rowAttrs[r] = gvm.rowAttrs[r+1]
		}
	}

	rwModel.Resize(h-1, w)
	gvm.modelWasResized()
	return nil
}

// Deletes a column of a model
func (gvm *ModelViewCtrl) DeleteCol(col int) error {
	rwModel, isRWModel := gvm.model.(RWModel)
	if !isRWModel {
		return ErrModelReadOnly
	}

	h, w := rwModel.Dimensions()
	for c := col; c < w-1; c++ {
		for r := 0; r < h; r++ {
			rwModel.SetCellValue(r, c, rwModel.CellValue(r, c+1))
			gvm.colAttrs[c] = gvm.colAttrs[c+1]
		}
	}

	rwModel.Resize(h, w-1)
	gvm.modelWasResized()
	return nil
}

func (gvm *ModelViewCtrl) modelWasResized() {
	rows, cols := gvm.model.Dimensions()
	gvm.rowAttrs = gvm.resizeAttrSlice(gvm.rowAttrs, rows, DefaultRowAttrs)
	gvm.colAttrs = gvm.resizeAttrSlice(gvm.colAttrs, cols, DefaultColAttrs)
}

func (gvm *ModelViewCtrl) resizeAttrSlice(oldSlice []SliceAttr, newSize int, defaultAttrs SliceAttr) []SliceAttr {
	oldLen := len(oldSlice)
	newSlice := oldSlice

	if newSize > oldLen {
		newSlice = make([]SliceAttr, newSize)
		for i := 0; i < newSize; i++ {
			if i < oldLen {
				newSlice[i] = oldSlice[i]
			} else {
				newSlice[i] = defaultAttrs
			}
		}
	} else {
		newSlice = newSlice[:newSize]
	}
	return newSlice
}

type SliceAttr struct {
	Size   int
	Marker Marker
}

type Marker int

const (
	MarkerNone Marker = iota
	MarkerRed
	MarkerGreen
	MarkerBlue
)

var DefaultRowAttrs = SliceAttr{Size: 1}
var DefaultColAttrs = SliceAttr{Size: 24}

var ErrModelReadOnly = errors.New("ModelVC is read-only")
