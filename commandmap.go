package main

import (
	"errors"
	"fmt"

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

func (cm *CommandMapping) DoEval(ctx *CommandContext, expr string) {
	if err := cm.Eval(ctx, expr); err != nil {
		ctx.ShowError(err)
	}
}

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

	cm.Define("open-right", "Inserts a column to the right of the curser", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, _ := grid.CellPosition()

		if rwModel, isRwModel := ctx.Session().Model.(RWModel); isRwModel {
			height, width := rwModel.Dimensions()
			if cellX == width-1 {
				rwModel.Resize(height, width+1)
			}
			return nil
		}

		return errors.New("model is read-only")
	})

	cm.Define("open-down", "Inserts a row below the curser", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		_, cellY := grid.CellPosition()

		if rwModel, isRwModel := ctx.Session().Model.(RWModel); isRwModel {
			height, width := rwModel.Dimensions()
			if cellY == height-1 {
				rwModel.Resize(height+1, width)
			}
			return nil
		}

		return errors.New("model is read-only")
	})

	cm.Define("append", "Inserts a row below the curser", "", func(ctx *CommandContext) error {
		if err := ctx.Session().Commands.Eval(ctx, "open-down"); err != nil {
			return err
		}
		if err := ctx.Session().Commands.Eval(ctx, "move-down"); err != nil {
			return err
		}

		ctx.Session().UIManager.Redraw()

		return ctx.Session().Commands.Eval(ctx, "set-cell")
	})

	cm.Define("enter-command", "Enter command", "", func(ctx *CommandContext) error {
		ctx.Frame().Prompt(":", func(res string) {
			cm.DoEval(ctx, res)
		})
		return nil
	})

	cm.Define("set-cell", "Change the value of the selected cell", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, cellY := grid.CellPosition()
		if rwModel, isRwModel := ctx.Session().Model.(RWModel); isRwModel {
			ctx.Frame().Prompt("> ", func(res string) {
				rwModel.SetCellValue(cellY, cellX, res)
			})
		}
		return nil
	})

	cm.Define("save", "Save current file", "", func(ctx *CommandContext) error {
		wSource, isWSource := ctx.Session().Source.(WritableModelSource)
		if !isWSource {
			return fmt.Errorf("model is not writable")
		}

		if err := wSource.Write(ctx.Session().Model); err != nil {
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

	cm.MapKey('e', cm.Command("set-cell"))

	cm.MapKey('a', cm.Command("append"))

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
