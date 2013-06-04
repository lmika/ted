/**
 * The grid component.  This is used for displaying the model.
 */
package main

import "github.com/nsf/termbox-go"
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
    cellX           int             // Left most cell
    cellY           int             // Top most cell
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
    return &Grid{model, 0, 0}
}

/**
 * Returns the requested dimensions of a grid (as required by UiComponent)
 */
func (grid *Grid) RequestDims() (int, int) {
    return -1, -1
}

/**
 * Shifts the viewport of the grid.
 */
func (grid *Grid) ShiftBy(x int, y int) {
    grid.cellX += x
    grid.cellY += y
}


// Gets the cell value and attributes of a particular cell
func (grid *Grid) getCellData(cellX, cellY int) (text string, fg, bg termbox.Attribute) {
    // The fixed cells
    modelCellX := cellX - 1 + grid.cellX
    modelCellY := cellY - 1 + grid.cellY
    modelMaxX, modelMaxY := grid.model.GetDimensions()
        
    if (cellX == 0) && (cellY == 0) {
        return "", termbox.AttrBold, termbox.AttrBold
    } else if (cellX == 0) {
        return strconv.Itoa(modelCellY), termbox.AttrBold, termbox.AttrBold
    } else if (cellY == 0) {
        return strconv.Itoa(modelCellX), termbox.AttrBold, termbox.AttrBold
    } else {
        // The data from the model
        if (modelCellX >= 0) && (modelCellY >= 0) && (modelCellX < modelMaxX) && (modelCellY < modelMaxY) {        
            return grid.model.GetCellValue(modelCellX, modelCellY), 0, 0
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
    
    modelCellX := cellX - 1 + grid.cellX
    modelCellY := cellY - 1 + grid.cellY
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
func (grid *Grid) renderCell(cellClipRect gridRect, sx int, sy int, text string, fg, bg termbox.Attribute) {
    for x := cellClipRect.x1; x <= cellClipRect.x2; x++ {
        for y := cellClipRect.y1; y <= cellClipRect.y2; y++ {
            currRune := ' '
            if (y == 0) {
                textPos := int(x)
                if textPos < len(text) {
                    currRune = rune(text[textPos])
                }
            }
                
            termbox.SetCell(int(x - cellClipRect.x1) + sx, int(y - cellClipRect.y1) + sy, currRune, fg, bg)
        }
    }
}


// Renders a column.  The viewport determines the maximum position of the rendered cell.  CellX and CellY are the
// cell indicies to render, cellOffset are the LOCAL offset of the cell.
// This function will return the new X position (gridRect.x1 + colWidth)
func (grid *Grid) renderColumn(screenViewPort gridRect, cellX int, cellY int, cellOffsetX int, cellOffsetY int) (int) {

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

    for screenY < maxScreenY {

        // Cap the row height if it will go beyond the edge of the viewport.
        _, rowHeight := grid.getCellDimensions(cellX, cellY)
        if screenY + rowHeight > maxScreenY {
            rowHeight = maxScreenY - screenY
        }
        
        cellText, cellFg, cellBg := grid.getCellData(cellX, cellY)

        grid.renderCell(newGridRect(cellOffsetX, cellOffsetY, colWidth - cellOffsetX, rowHeight),
                screenX, screenY, cellText, cellFg, cellBg)  // termbox.AttrReverse, termbox.AttrReverse

        cellY = cellY + 1
        screenY = screenY + rowHeight - cellOffsetY
        cellOffsetY = 0
    }

    return screenX + colWidth
}


// Renders the grid.
func (grid *Grid) renderGrid(screenViewPort gridRect, cellX int, cellY int, cellOffsetX int, cellOffsetY int) {
    
    for screenViewPort.x1 < screenViewPort.x2 {
        screenViewPort.x1 = gridPoint(grid.renderColumn(screenViewPort, cellX, cellY, cellOffsetX, cellOffsetY))
        cellX = cellX + 1
        cellOffsetX = 0
    }
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
func (grid *Grid) Redraw(x int, y int, w int, h int) {
    viewportRect := newGridRect(x, y, x + w, y + h)
    grid.renderGrid(viewportRect, 0, 0, 0, 0)
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
