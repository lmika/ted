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
    GetDimensions() (int, int)

    /**
     * Returns the size of the particular column.  If the size is 0, this indicates that the column is hidden.
     */
    GetColWidth(int) int

    /**
     * Returns the size of the particular row.  If the size is 0, this indicates that the row is hidden.
     */
    GetRowHeight(int) int

    /**
     * Returns the value of the cell a position X, Y
     */
    GetCellValue(int, int) string
}


type gridPoint      int

/**
 * The grid component.
 */
type Grid struct {
    model           GridModel
    viewCellX       int             // Left most cell
    viewCellY       int             // Top most cell
    selCellX        int             // The currently selected cell
    selCellY        int
    cellsWide       int             // Measured number of cells.  Recalculated on redraw.
    cellsHigh       int
}

/**
 * Clipping rectangle
 */
type gridRect struct {
    x1              gridPoint
    y1              gridPoint
    x2              gridPoint
    y2              gridPoint
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

/**
 * Returns the requested dimensions of a grid (as required by UiComponent)
 */
func (grid *Grid) Remeasure(w, h int) (int, int) {
    return w, h
}

/**
 * Shifts the viewport of the grid.
 */
func (grid *Grid) ShiftBy(x int, y int) {
    grid.viewCellX += x
    grid.viewCellY += y
}


// Moves the currently selected cell by a delta.
func (grid *Grid) MoveBy(x int, y int) {
    grid.selCellX += x
    grid.selCellY += y
    grid.reposition()
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
    modelMaxX, modelMaxY := grid.model.GetDimensions()
        
    if (cellX == 0) && (cellY == 0) {
        return strconv.Itoa(grid.cellsWide), AttrBold, AttrBold
    } else if (cellX == 0) {
        if (modelCellY == grid.selCellY) {
            return strconv.Itoa(modelCellY), AttrBold | AttrReverse, AttrBold | AttrReverse
        } else {
            return strconv.Itoa(modelCellY), AttrBold, AttrBold
        }
    } else if (cellY == 0) {
        if (modelCellX == grid.selCellX) {
            return strconv.Itoa(modelCellX), AttrBold | AttrReverse, AttrBold | AttrReverse
        } else {
            return strconv.Itoa(modelCellX), AttrBold, AttrBold
        }
    } else {
        // The data from the model
        if (modelCellX >= 0) && (modelCellY >= 0) && (modelCellX < modelMaxX) && (modelCellY < modelMaxY) {     
            if (modelCellX == grid.selCellX) && (modelCellY == grid.selCellY) {   
                return grid.model.GetCellValue(modelCellX, modelCellY), AttrReverse, AttrReverse
            } else {
                return grid.model.GetCellValue(modelCellX, modelCellY), 0, 0
            }
        } else {
            return "~", 0, 0
        }
    }    
    
    // XXX: Workaround for bug in compiler
    panic("Unreachable code")
    return "", 0, 0
}

// Gets the cell dimensions
func (grid *Grid) getCellDimensions(cellX, cellY int) (width, height int) {

    var cellWidth, cellHeight int
    
    modelCellX := cellX - 1 + grid.viewCellX
    modelCellY := cellY - 1 + grid.viewCellY
    modelMaxX, modelMaxY := grid.model.GetDimensions()
    
    // Get the cell width & height from model (if within range)
    if (modelCellX >= 0) && (modelCellX < modelMaxX) {
        cellWidth = grid.model.GetColWidth(modelCellX)
    } else {
        cellWidth = 8
    }
    
    if (modelCellY >= 0) && (modelCellY < modelMaxY) {
        cellHeight = grid.model.GetRowHeight(modelCellY)
    } else {
        cellHeight = 2
    }        
    
    if (cellX == 0) && (cellY == 0) {
        return 8, 1
    } else if (cellX == 0) {
        return 8, cellHeight
    } else if (cellY == 0) {
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
            if (y == 0) {
                textPos := int(x)
                if textPos < len(text) {
                    currRune = rune(text[textPos])
                }
            }
                
            //termbox.SetCell(int(x - cellClipRect.x1) + sx, int(y - cellClipRect.y1) + sy, currRune, fg, bg)

            // TODO: This might be better if this wasn't so low-level
            ctx.DrawRuneWithAttrs(int(x - cellClipRect.x1) + sx, int(y - cellClipRect.y1) + sy, currRune, fg, bg)
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
        if screenY + rowHeight > maxScreenY {
            rowHeight = maxScreenY - screenY
        }
        
        cellText, cellFg, cellBg := grid.getCellData(cellX, cellY)

        grid.renderCell(ctx, newGridRect(cellOffsetX, cellOffsetY, colWidth - cellOffsetX, rowHeight),
                screenX, screenY, cellText, cellFg, cellBg)  // termbox.AttrReverse, termbox.AttrReverse

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
    var wid, hei int = grid.model.GetDimensions()
    posX = 0
    posY = 0

    cellX = -1
    cellY = -1

    // Go through columns to locate the particular cellX
    for cx := 0; cx < wid; cx++ {
        if (x >= posX) && (x < posX + grid.model.GetColWidth(cx)) {
            // We found the X position
            cellX = int(cx)
            break
        }
    }

    for cy := 0; cy < hei; cy++ {
        if (y >= posY) && (y < posY + grid.model.GetRowHeight(cy)) {
            // And the Y position
            cellY = int(cy)
            break
        }
    }

    return
}

/**
 * Redraws the grid.
 */
func (grid *Grid) Redraw(ctx *DrawContext) {
    viewportRect := newGridRect(0, 0, ctx.W, ctx.H)
    grid.cellsWide, grid.cellsHigh = grid.renderGrid(ctx, viewportRect, 0, 0, 0, 0)
}

// --------------------------------------------------------------------------------------------
// Test Model

type TestModel struct {
    thing       int
}

/**
 * Returns the size of the grid model (width x height)
 */
func (model *TestModel) GetDimensions() (int, int) {
    return 100, 100
}

/**
 * Returns the size of the particular column.  If the size is 0, this indicates that the column is hidden.
 */
func (model *TestModel) GetColWidth(int) int {
    return 16
}

/**
 * Returns the size of the particular row.  If the size is 0, this indicates that the row is hidden.
 */
func (model *TestModel) GetRowHeight(int) int {
    return 1
}

/**
 * Returns the value of the cell a position X, Y
 */
func (model *TestModel) GetCellValue(x int, y int) string {
    return strconv.Itoa(x) + "," + strconv.Itoa(y)
}
