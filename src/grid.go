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
    offsetX         int             // Offset of the viewport (REAL characters, not cells)
    offsetY         int             // Offset of the viewport
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
    grid.offsetX += x
    grid.offsetY += y
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
    colWidth := grid.model.GetColWidth(cellX) - cellOffsetX
    if colWidth > screenWidth {
        colWidth = screenHeight
    }

    // The maximum
    maxScreenY := screenY + screenHeight

    for screenY < maxScreenY {

        // Cap the row height if it will go beyond the edge of the viewport.
        rowHeight := grid.model.GetRowHeight(cellY)
        if screenY + rowHeight > maxScreenY {
            rowHeight = maxScreenY - screenY
        }

        grid.renderCell(newGridRect(cellOffsetX, cellOffsetY, colWidth - cellOffsetX, rowHeight),
                screenX, screenY, grid.model.GetCellValue(cellX, cellY), 0, 0)  // termbox.AttrReverse, termbox.AttrReverse

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

    cellX, cellY, posX, posY := grid.pointToCell(grid.offsetX, grid.offsetY)
/*
    grid.renderCell(gridRect{0, 0, 14, 0}, 0, 0, "Hello", termbox.AttrReverse, termbox.AttrReverse)
    grid.renderCell(gridRect{1, 0, 14, 0}, 0, 1, "Hello", termbox.AttrReverse, termbox.AttrReverse)
    grid.renderCell(gridRect{2, 0, 14, 0}, 0, 2, "Hello", termbox.AttrReverse, termbox.AttrReverse)
    grid.renderCell(gridRect{3, 0, 14, 0}, 0, 3, "Hello", termbox.AttrReverse, termbox.AttrReverse)
    grid.renderCell(gridRect{4, 0, 14, 0}, 0, 4, "Hello", termbox.AttrReverse, termbox.AttrReverse)
    grid.renderCell(gridRect{5, 0, 14, 0}, 0, 5, "Hello", termbox.AttrReverse, termbox.AttrReverse)
*/
    grid.renderGrid(viewportRect, cellX, cellY, grid.offsetX - posX, grid.offsetY - posY)
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
