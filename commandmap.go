package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/lmika/ted/ui"
)

const (
	ModAlt rune = 1<<31 - 1
)

// A command
type Command struct {
	Name string
	Doc  string

	// TODO: Add argument mapping which will fetch properties from the environment
	Action func(ctx *CommandContext) error
}

// Execute the command
func (cmd *Command) Do(ctx *CommandContext) error {
	return cmd.Action(ctx)
}

// A command mapping
type CommandMapping struct {
	Commands    map[string]*Command
	KeyMappings map[rune]*Command
}

// Creates a new, empty command mapping
func NewCommandMapping() *CommandMapping {
	return &CommandMapping{make(map[string]*Command), make(map[rune]*Command)}
}

// Adds a new command
func (cm *CommandMapping) Define(name string, doc string, opts string, fn func(ctx *CommandContext) error) {
	cm.Commands[name] = &Command{name, doc, fn}
}

// Adds a key mapping
func (cm *CommandMapping) MapKey(key rune, cmd *Command) {
	cm.KeyMappings[key] = cmd
}

// Searches for a command by name.  Returns the command or null
func (cm *CommandMapping) Command(name string) *Command {
	return cm.Commands[name]
}

// Searches for a command by key mapping
func (cm *CommandMapping) KeyMapping(key rune) *Command {
	return cm.KeyMappings[key]
}

// Evaluate a command
func (cm *CommandMapping) Eval(ctx *CommandContext, expr string) error {
	// TODO: Use propper expression language here
	cmd := cm.Commands[expr]
	if cmd != nil {
		return cmd.Do(ctx)
	}

	return fmt.Errorf("no such command: %v", expr)
}

//func (cm *CommandMapping) DoEval(ctx *CommandContext, expr string) {
//	if err := cm.Eval(ctx, expr); err != nil {
//		ctx.ShowError(err)
//	}
//}

// Registers the standard view navigation commands.  These commands require the frame
func (cm *CommandMapping) RegisterViewCommands() {
	cm.Define("move-down", "Moves the cursor down one row", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(0, 1) }))
	cm.Define("move-up", "Moves the cursor up one row", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(0, -1) }))
	cm.Define("move-left", "Moves the cursor left one column", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(-1, 0) }))
	cm.Define("move-right", "Moves the cursor right one column", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(1, 0) }))

	// TODO: Pages are just 25 rows and 15 columns at the moment
	cm.Define("page-down", "Moves the cursor down one page", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(0, 25) }))
	cm.Define("page-up", "Moves the cursor up one page", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(0, -25) }))
	cm.Define("page-left", "Moves the cursor left one page", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(-15, 0) }))
	cm.Define("page-right", "Moves the cursor right one page", "", gridNavOperation(func(grid *ui.Grid) { grid.MoveBy(15, 0) }))

	cm.Define("row-top", "Moves the cursor to the top of the row", "", gridNavOperation(func(grid *ui.Grid) {
		cellX, _ := grid.CellPosition()
		grid.MoveTo(cellX, 0)
	}))
	cm.Define("row-bottom", "Moves the cursor to the bottom of the row", "", gridNavOperation(func(grid *ui.Grid) {
		cellX, _ := grid.CellPosition()
		_, dimY := grid.Model().Dimensions()
		grid.MoveTo(cellX, dimY-1)
	}))
	cm.Define("col-left", "Moves the cursor to the left-most column", "", gridNavOperation(func(grid *ui.Grid) {
		_, cellY := grid.CellPosition()
		grid.MoveTo(0, cellY)
	}))
	cm.Define("col-right", "Moves the cursor to the right-most column", "", gridNavOperation(func(grid *ui.Grid) {
		_, cellY := grid.CellPosition()
		dimX, _ := grid.Model().Dimensions()
		grid.MoveTo(dimX-1, cellY)
	}))

	cm.Define("delete-row", "Removes the currently selected row", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		_, cellY := grid.CellPosition()

		return ctx.ModelVC().DeleteRow(cellY)
	})
	cm.Define("delete-col", "Removes the currently selected column", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, _ := grid.CellPosition()

		return ctx.ModelVC().DeleteCol(cellX)
	})
	cm.Define("search", "Search for a cell", "", func(ctx *CommandContext) error {
		ctx.Frame().Prompt(PromptOptions{Prompt: "/"}, func(res string) error {
			re, err := regexp.Compile(res)
			if err != nil {
				return fmt.Errorf("invalid regexp: %v", err)
			}

			ctx.session.LastSearch = re
			return ctx.Session().Commands.Eval(ctx, "search-next")
		})
		return nil
	})
	cm.Define("search-next", "Goto the next cell", "", func(ctx *CommandContext) error {
		if ctx.session.LastSearch == nil {
			ctx.Session().Commands.Eval(ctx, "search")
		}

		height, width := ctx.ModelVC().Model().Dimensions()
		startX, startY := ctx.Frame().Grid().CellPosition()
		cellX, cellY := startX, startY

		for {
			cellX++
			if cellX >= width {
				cellX = 0
				cellY = (cellY + 1) % height
			}
			if ctx.session.LastSearch.MatchString(ctx.ModelVC().Model().CellValue(cellY, cellX)) {
				ctx.Frame().Grid().MoveTo(cellX, cellY)
				return nil
			} else if (cellX == startX) && (cellY == startY) {
				return errors.New("No match found")
			}
		}
	})

	cm.Define("open-right", "Inserts a column to the right of the curser", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, _ := grid.CellPosition()

		height, width := ctx.ModelVC().Model().Dimensions()
		if cellX == width-1 {
			return ctx.ModelVC().Resize(height, width+1)
		}
		return nil
	})

	cm.Define("open-down", "Inserts a row below the curser", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		_, cellY := grid.CellPosition()

		height, width := ctx.ModelVC().Model().Dimensions()
		if cellY == height-1 {
			return ctx.ModelVC().Resize(height+1, width)
		}
		return nil
	})

	cm.Define("append", "Inserts a row below the curser", "", func(ctx *CommandContext) error {
		if err := ctx.Session().Commands.Eval(ctx, "open-down"); err != nil {
			return err
		}
		if err := ctx.Session().Commands.Eval(ctx, "move-down"); err != nil {
			return err
		}

		ctx.Session().UIManager.Redraw()

		return ctx.Session().Commands.Eval(ctx, "edit-cell")
	})

	cm.Define("inc-col-width", "Increase the width of the current column", "", func(ctx *CommandContext) error {
		cellX, _ := ctx.Frame().Grid().CellPosition()

		attrs := ctx.ModelVC().ColAttrs(cellX)
		attrs.Size += 2
		ctx.ModelVC().SetColAttrs(cellX, attrs)
		return nil
	})

	cm.Define("dec-col-width", "Decrease the width of the current column", "", func(ctx *CommandContext) error {
		cellX, _ := ctx.Frame().Grid().CellPosition()

		attrs := ctx.ModelVC().ColAttrs(cellX)
		attrs.Size -= 2
		if attrs.Size < 4 {
			attrs.Size = 4
		}
		ctx.ModelVC().SetColAttrs(cellX, attrs)
		return nil
	})

	cm.Define("clear-row-marker", "Clears any row markers", "", func(ctx *CommandContext) error {
		_, cellY := ctx.Frame().Grid().CellPosition()

		attrs := ctx.ModelVC().RowAttrs(cellY)
		attrs.Marker = MarkerNone
		ctx.ModelVC().SetRowAttrs(cellY, attrs)
		return nil
	})

	cm.Define("mark-row-red", "Set row marker to red", "", func(ctx *CommandContext) error {
		_, cellY := ctx.Frame().Grid().CellPosition()

		attrs := ctx.ModelVC().RowAttrs(cellY)
		attrs.Marker = MarkerRed
		ctx.ModelVC().SetRowAttrs(cellY, attrs)
		return nil
	})

	cm.Define("mark-row-green", "Set row marker to green", "", func(ctx *CommandContext) error {
		_, cellY := ctx.Frame().Grid().CellPosition()

		attrs := ctx.ModelVC().RowAttrs(cellY)
		attrs.Marker = MarkerGreen
		ctx.ModelVC().SetRowAttrs(cellY, attrs)
		return nil
	})

	cm.Define("mark-row-blue", "Set row marker to blue", "", func(ctx *CommandContext) error {
		_, cellY := ctx.Frame().Grid().CellPosition()

		attrs := ctx.ModelVC().RowAttrs(cellY)
		attrs.Marker = MarkerBlue
		ctx.ModelVC().SetRowAttrs(cellY, attrs)
		return nil
	})

	cm.Define("enter-command", "Enter command", "", func(ctx *CommandContext) error {
		ctx.Frame().Prompt(PromptOptions{Prompt: ":"}, func(res string) error {
			return cm.Eval(ctx, res)
		})
		return nil
	})

	cm.Define("replace-cell", "Replace the value of the selected cell", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, cellY := grid.CellPosition()
		if _, isRwModel := ctx.ModelVC().Model().(RWModel); isRwModel {
			ctx.Frame().Prompt(PromptOptions{Prompt: "> "}, func(res string) error {
				if err := ctx.ModelVC().SetCellValue(cellY, cellX, res); err != nil {
					return err
				}
				ctx.Frame().ShowCellValue()
				return nil
			})
		}
		return nil
	})
	cm.Define("edit-cell", "Modify the value of the selected cell", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, cellY := grid.CellPosition()

		if _, isRwModel := ctx.ModelVC().Model().(RWModel); isRwModel {
			ctx.Frame().Prompt(PromptOptions{
				Prompt:       "> ",
				InitialValue: grid.Model().CellValue(cellX, cellY),
			}, func(res string) error {
				if err := ctx.ModelVC().SetCellValue(cellY, cellX, res); err != nil {
					return err
				}
				ctx.Frame().ShowCellValue()
				return nil
			})
		}
		return nil
	})

	cm.Define("save", "Save current file", "", func(ctx *CommandContext) error {
		wSource, isWSource := ctx.Session().Source.(WritableModelSource)
		if !isWSource {
			return fmt.Errorf("model is not writable")
		}

		if err := wSource.Write(ctx.ModelVC().Model()); err != nil {
			return err
		}

		ctx.Frame().Message("Wrote " + wSource.String())
		return nil
	})

	cm.Define("quit", "Quit TED", "", func(ctx *CommandContext) error {
		ctx.Session().UIManager.Shutdown()
		return nil
	})

	cm.Define("save-and-quit", "Save current file, then quit", "", func(ctx *CommandContext) error {
		if err := cm.Eval(ctx, "save"); err != nil {
			return nil
		}

		return cm.Eval(ctx, "quit")
	})

	// Aliases
	cm.Commands["w"] = cm.Command("save")
	cm.Commands["q"] = cm.Command("quit")
	cm.Commands["wq"] = cm.Command("save-and-quit")
}

// Registers the standard view key bindings.  These commands require the frame
func (cm *CommandMapping) RegisterViewKeyBindings() {
	cm.MapKey('i', cm.Command("move-up"))
	cm.MapKey('k', cm.Command("move-down"))
	cm.MapKey('j', cm.Command("move-left"))
	cm.MapKey('l', cm.Command("move-right"))
	cm.MapKey('I', cm.Command("page-up"))
	cm.MapKey('K', cm.Command("page-down"))
	cm.MapKey('J', cm.Command("page-left"))
	cm.MapKey('L', cm.Command("page-right"))
	cm.MapKey(ui.KeyCtrlI, cm.Command("row-top"))
	cm.MapKey(ui.KeyCtrlK, cm.Command("row-bottom"))
	cm.MapKey(ui.KeyCtrlJ, cm.Command("col-left"))
	cm.MapKey(ui.KeyCtrlL, cm.Command("col-right"))

	cm.MapKey(ui.KeyArrowUp, cm.Command("move-up"))
	cm.MapKey(ui.KeyArrowDown, cm.Command("move-down"))
	cm.MapKey(ui.KeyArrowLeft, cm.Command("move-left"))
	cm.MapKey(ui.KeyArrowRight, cm.Command("move-right"))

	cm.MapKey('e', cm.Command("edit-cell"))
	cm.MapKey('r', cm.Command("replace-cell"))

	cm.MapKey('a', cm.Command("append"))

	cm.MapKey('D', cm.Command("delete-row"))

	cm.MapKey('/', cm.Command("search"))
	cm.MapKey('n', cm.Command("search-next"))

	cm.MapKey('0', cm.Command("clear-row-marker"))
	cm.MapKey('1', cm.Command("mark-row-red"))
	cm.MapKey('2', cm.Command("mark-row-green"))
	cm.MapKey('3', cm.Command("mark-row-blue"))

	cm.MapKey('{', cm.Command("dec-col-width"))
	cm.MapKey('}', cm.Command("inc-col-width"))

	cm.MapKey(':', cm.Command("enter-command"))
}

// A nativation command factory.  This will perform the passed in operation with the current grid and
// will display the cell value in the message box.
func gridNavOperation(op func(grid *ui.Grid)) func(ctx *CommandContext) error {
	return func(ctx *CommandContext) error {
		op(ctx.Frame().Grid())
		ctx.Frame().ShowCellValue()
		return nil
	}
}
