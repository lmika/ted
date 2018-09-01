/**
 * The grid component.  This is used for displaying the model.
 */
package ui

import "strconv"

/**
 * An abstract display model.
 */
type GridModel interface {
	/**
	 * Returns the size of the grid model (width x height)
	 */
	Dimensions() (int, int)

	/**
	 * Returns the size of the particular column.  If the size is 0, this indicates that the column is hidden.
	 */
	ColWidth(int) int

	/**
	 * Returns the size of the particular row.  If the size is 0, this indicates that the row is hidden.
	 */
	RowHeight(int) int

	/**
	 * Returns the value of the cell a position X, Y
	 */
	CellValue(int, int) string
}

type gridPoint int

/**
 * The grid component.
 */
type Grid struct {
	model GridModel // The grid model

	viewCellX int // Left most cell
	viewCellY int // Top most cell
	selCellX  int // The currently selected cell
	selCellY  int
	cellsWide int // Measured number of cells.  Recalculated on redraw.
	cellsHigh int
}

/**
 * Clipping rectangle
 */
type gridRect struct {
	x1 gridPoint
	y1 gridPoint
	x2 gridPoint
	y2 gridPoint
}

/**
 * Creates a new gridRect from integers.
 */
func newGridRect(x1, y1, x2, y2 int) gridRect {
	return gridRect{gridPoint(x1), gridPoint(y1), gridPoint(x2), gridPoint(y2)}
}

/**
 * Creates a new grid.
 */
func NewGrid(model GridModel) *Grid {
	return &Grid{model, 0, 0, 0, 0, -1, -1}
}

// Returns the model
func (grid *Grid) Model() GridModel {
	return grid.model
}

// Sets the model
func (grid *Grid) SetModel(model GridModel) {
	grid.model = model
}

/**
 * Shifts the viewport of the grid.
 */
func (grid *Grid) ShiftBy(x int, y int) {
	grid.viewCellX += x
	grid.viewCellY += y
}

// Returns the display value of the currently selected cell.
func (grid *Grid) CurrentCellDisplayValue() string {
	if grid.isCellValid(grid.selCellX, grid.selCellY) {
		return grid.model.CellValue(grid.selCellX, grid.selCellY)
	} else {
		return ""
	}
}

// Moves the currently selected cell by a delta.  This will be implemented as single stepped
// moveTo calls to handle invalid cells.
func (grid *Grid) MoveBy(x int, y int) {
	grid.MoveTo(grid.selCellX+x, grid.selCellY+y)
}

// Moves the currently selected cell to a specific row.  The row must be valid, otherwise the
// currently selected cell will not be changed.  Returns true if the move was successful
func (grid *Grid) MoveTo(newX, newY int) {
	maxX, maxY := grid.model.Dimensions()
	newX = intMinMax(newX, 0, maxX-1)
	newY = intMinMax(newY, 0, maxY-1)

	if grid.isCellValid(newX, newY) {
		grid.selCellX = newX
		grid.selCellY = newY
		grid.reposition()
	}
}

// Returns the currently selected cell position.
func (grid *Grid) CellPosition() (int, int) {
	return grid.selCellX, grid.selCellY
}

// Returns true if the user can enter the specific cell
func (grid *Grid) isCellValid(x int, y int) bool {
	maxX, maxY := grid.model.Dimensions()
	return (x >= 0) && (y >= 0) && (x < maxX) && (y < maxY)
}

// Determine the topmost cell based on the location of the currently selected cell
func (grid *Grid) reposition() {

	// If we have no measurement information, forget it.
	if (grid.cellsWide == -1) || (grid.cellsHigh == -1) {
		return
	}

	if grid.selCellX < grid.viewCellX {
		grid.viewCellX = grid.selCellX
	} else if grid.selCellX >= (grid.viewCellX + grid.cellsWide - 3) {
		grid.viewCellX = grid.selCellX - (grid.cellsWide - 3)
	}

	if grid.selCellY < grid.viewCellY {
		grid.viewCellY = grid.selCellY
	} else if grid.selCellY >= (grid.viewCellY + grid.cellsHigh - 3) {
		grid.viewCellY = grid.selCellY - (grid.cellsHigh - 3)
	}
}

// Gets the cell value and attributes of a particular cell
func (grid *Grid) getCellData(cellX, cellY int) (text string, fg, bg Attribute) {
	// The fixed cells
	modelCellX := cellX - 1 + grid.viewCellX
	modelCellY := cellY - 1 + grid.viewCellY
	modelMaxX, modelMaxY := grid.model.Dimensions()

	if (cellX == 0) && (cellY == 0) {
		return "", AttrBold, AttrBold
	} else if cellX == 0 {
		if modelCellY == grid.selCellY {
			return strconv.Itoa(modelCellY), AttrBold | AttrReverse, AttrReverse
		} else {
			return strconv.Itoa(modelCellY), AttrBold, 0
		}
	} else if cellY == 0 {
		if modelCellX == grid.selCellX {
			return strconv.Itoa(modelCellX), AttrBold | AttrReverse, AttrReverse
		} else {
			return strconv.Itoa(modelCellX), AttrBold, 0
		}
	} else {
		// The data from the model
		if (modelCellX >= 0) && (modelCellY >= 0) && (modelCellX < modelMaxX) && (modelCellY < modelMaxY) {
			if (modelCellX == grid.selCellX) && (modelCellY == grid.selCellY) {
				return grid.model.CellValue(modelCellX, modelCellY), AttrReverse, AttrReverse
			} else {
				return grid.model.CellValue(modelCellX, modelCellY), 0, 0
			}
		} else {
			return "~", ColorBlue, 0
		}
	}
}

// Gets the cell dimensions
func (grid *Grid) getCellDimensions(cellX, cellY int) (width, height int) {

	var cellWidth, cellHeight int

	modelCellX := cellX - 1 + grid.viewCellX
	modelCellY := cellY - 1 + grid.viewCellY
	modelMaxX, modelMaxY := grid.model.Dimensions()

	// Get the cell width & height from model (if within range)
	if (modelCellX >= 0) && (modelCellX < modelMaxX) {
		cellWidth = grid.model.ColWidth(modelCellX)
	} else {
		cellWidth = 8
	}

	if (modelCellY >= 0) && (modelCellY < modelMaxY) {
		cellHeight = grid.model.RowHeight(modelCellY)
	} else {
		cellHeight = 1
	}

	if (cellX == 0) && (cellY == 0) {
		return 8, 1
	} else if cellX == 0 {
		return 8, cellHeight
	} else if cellY == 0 {
		return cellWidth, 1
	} else {
		return cellWidth, cellHeight
	}

	// XXX: Workaround for bug in compiler
	panic("Unreachable code")
	return 0, 0
}

/**
 * Renders a cell which contains text.  The clip rectangle defines the size of the cell, as well as the left offset
 * of the cell.  The sx and sy determine the screen position of the cell top-left.
 */
func (grid *Grid) renderCell(ctx *DrawContext, cellClipRect gridRect, sx int, sy int, text string, fg, bg Attribute) {
	for x := cellClipRect.x1; x <= cellClipRect.x2; x++ {
		for y := cellClipRect.y1; y <= cellClipRect.y2; y++ {
			currRune := ' '
			if y == 0 {
				textPos := int(x)
				if textPos < len(text) {
					currRune = rune(text[textPos])
				}
			}

			// TODO: This might be better if this wasn't so low-level
			ctx.DrawRuneWithAttrs(int(x-cellClipRect.x1)+sx, int(y-cellClipRect.y1)+sy, currRune, fg, bg)
		}
	}
}

// Renders a column.  The viewport determines the maximum position of the rendered cell.  CellX and CellY are the
// cell indicies to render, cellOffset are the LOCAL offset of the cell.
// This function will return the new X position (gridRect.x1 + colWidth)
func (grid *Grid) renderColumn(ctx *DrawContext, screenViewPort gridRect, cellX int, cellY int, cellOffsetX int, cellOffsetY int) (gridPoint, int) {

	// The top-left position of the column
	screenX := int(screenViewPort.x1)
	screenY := int(screenViewPort.y1)
	screenWidth := int(screenViewPort.x2 - screenViewPort.x1)
	screenHeight := int(screenViewPort.y2 - screenViewPort.y1)

	// Work out the column width and cap it if it will spill over the edge of the viewport
	colWidth, _ := grid.getCellDimensions(cellX, cellY)
	colWidth -= cellOffsetX
	if colWidth > screenWidth {
		colWidth = screenHeight
	}

	// The maximum
	maxScreenY := screenY + screenHeight
	cellsHigh := 0

	for screenY < maxScreenY {

		// Cap the row height if it will go beyond the edge of the viewport.
		_, rowHeight := grid.getCellDimensions(cellX, cellY)
		if screenY+rowHeight > maxScreenY {
			rowHeight = maxScreenY - screenY
		}

		cellText, cellFg, cellBg := grid.getCellData(cellX, cellY)

		grid.renderCell(ctx, newGridRect(cellOffsetX, cellOffsetY, colWidth-cellOffsetX, rowHeight),
			screenX, screenY, cellText, cellFg, cellBg) // termbox.AttrReverse, termbox.AttrReverse

		cellY++
		cellsHigh++
		screenY = screenY + rowHeight - cellOffsetY
		cellOffsetY = 0
	}

	return gridPoint(screenX + colWidth), cellsHigh
}

// Renders the grid.  Returns the number of cells in the X and Y direction were rendered.
//
func (grid *Grid) renderGrid(ctx *DrawContext, screenViewPort gridRect, cellX int, cellY int, cellOffsetX int, cellOffsetY int) (int, int) {

	var cellsHigh = 0
	var cellsWide = 0

	for screenViewPort.x1 < screenViewPort.x2 {
		screenViewPort.x1, cellsHigh = grid.renderColumn(ctx, screenViewPort, cellX, cellY, cellOffsetX, cellOffsetY)
		cellX = cellX + 1
		cellsWide++
		cellOffsetX = 0
	}

	return cellsWide, cellsHigh
}

/**
 * Returns the cell of the particular point, along with the top-left position of the cell.
 */
func (grid *Grid) pointToCell(x int, y int) (cellX int, cellY int, posX int, posY int) {
	var wid, hei int = grid.model.Dimensions()
	posX = 0
	posY = 0

	cellX = -1
	cellY = -1

	// Go through columns to locate the particular cellX
	for cx := 0; cx < wid; cx++ {
		if (x >= posX) && (x < posX+grid.model.ColWidth(cx)) {
			// We found the X position
			cellX = int(cx)
			break
		}
	}

	for cy := 0; cy < hei; cy++ {
		if (y >= posY) && (y < posY+grid.model.RowHeight(cy)) {
			// And the Y position
			cellY = int(cy)
			break
		}
	}

	return
}

/**
 * Returns the requested dimensions of a grid (as required by UiComponent)
 */
func (grid *Grid) Remeasure(w, h int) (int, int) {
	return w, h
}

/**
 * Redraws the grid.
 */
func (grid *Grid) Redraw(ctx *DrawContext) {
	viewportRect := newGridRect(0, 0, ctx.W, ctx.H)
	grid.cellsWide, grid.cellsHigh = grid.renderGrid(ctx, viewportRect, 0, 0, 0, 0)
}

// Called when the component has focus and a key has been pressed.
// This is the default behaviour of the grid, but it is not used by the main grid.
func (grid *Grid) KeyPressed(key rune, mod int) {
	// TODO: Not sure if this would be better handled using commands
	if (key == 'i') || (key == KeyArrowUp) {
		grid.MoveBy(0, -1)
	} else if (key == 'k') || (key == KeyArrowDown) {
		grid.MoveBy(0, 1)
	} else if (key == 'j') || (key == KeyArrowLeft) {
		grid.MoveBy(-1, 0)
	} else if (key == 'l') || (key == KeyArrowRight) {
		grid.MoveBy(1, 0)
	}
}

// --------------------------------------------------------------------------------------------
// Test Model

type TestModel struct {
	thing int
}

/**
 * Returns the size of the grid model (width x height)
 */
func (model *TestModel) Dimensions() (int, int) {
	return 100, 100
}

/**
 * Returns the size of the particular column.  If the size is 0, this indicates that the column is hidden.
 */
func (model *TestModel) ColWidth(int) int {
	return 16
}

/**
 * Returns the size of the particular row.  If the size is 0, this indicates that the row is hidden.
 */
func (model *TestModel) RowHeight(int) int {
	return 1
}

/**
 * Returns the value of the cell a position X, Y
 */
func (model *TestModel) CellValue(x int, y int) string {
	return strconv.Itoa(x) + "," + strconv.Itoa(y)
}
