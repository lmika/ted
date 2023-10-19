package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lmika/shellwords"

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
	toks := shellwords.Split(expr)
	if len(toks) == 0 {
		return nil
	}

	return cm.Invoke(ctx, toks[0], toks[1:])
}

func (cm *CommandMapping) Invoke(ctx *CommandContext, name string, args []string) error {
	cmd := cm.Commands[name]
	if cmd != nil {
		return cmd.Do(ctx.WithArgs(args))
	}

	return fmt.Errorf("no such command: %v", name)
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
	cm.Define("x-replace", "Performs a search and replace", "", func(ctx *CommandContext) error {
		if len(ctx.Args()) != 2 {
			return errors.New("Usage: x-replace MATCH REPLACEMENT")
		}

		match := ctx.Args()[0]
		repl := ctx.Args()[1]

		re, err := regexp.Compile(match)
		if err != nil {
			return fmt.Errorf("invalid regexp: %v", err)
		}

		matchCount := 0
		height, width := ctx.ModelVC().Model().Dimensions()
		for r := 0; r < height; r++ {
			for c := 0; c < width; c++ {
				cell := ctx.ModelVC().Model().CellValue(r, c)
				if re.FindStringIndex(cell) != nil {
					ctx.ModelVC().SetCellValue(r, c, re.ReplaceAllString(cell, repl))
					matchCount++
				}
			}
		}

		ctx.Frame().ShowMessage(fmt.Sprintf("Replaced %d matches", matchCount))
		return nil
	})

	cm.Define("open-right", "Inserts a column to the right of the curser", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, _ := grid.CellPosition()

		return ctx.ModelVC().OpenRight(cellX)
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
		ctx.Frame().Prompt(PromptOptions{
			Prompt:                 ":",
			CancelOnEmptyBackspace: true,
		}, func(res string) error {
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

		if _, isRwModel := ctx.ModelVC().Model().(RWModel); !isRwModel {
			return errors.New("Model is read-only")
		}

		if len(ctx.Args()) == 1 {
			if err := ctx.ModelVC().SetCellValue(cellY, cellX, ctx.Args()[0]); err != nil {
				return err
			}
		} else {
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
	cm.Define("yank", "Yank cell value", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, cellY := grid.CellPosition()

		// TODO: allow ranges
		ctx.Session().pasteBoard.SetCellValue(0, 0, grid.Model().CellValue(cellX, cellY))

		return nil
	})
	cm.Define("paste", "Paste cell value", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, cellY := grid.CellPosition()

		// TODO: allow ranges
		if _, isRwModel := ctx.ModelVC().Model().(RWModel); !isRwModel {
			return errors.New("Model is read-only")
		}
		if err := ctx.ModelVC().SetCellValue(cellY, cellX, ctx.Session().pasteBoard.CellValue(0, 0)); err != nil {
			return err
		}

		return nil
	})

	cm.Define("to-upper", "Convert cell value to uppercase", "", func(ctx *CommandContext) error {
		grid := ctx.Frame().Grid()
		cellX, cellY := grid.CellPosition()

		// TODO: allow ranges

		if _, isRwModel := ctx.ModelVC().Model().(RWModel); !isRwModel {
			return errors.New("Model is read-only")
		}

		currentValue := ctx.ModelVC().Model().CellValue(cellY, cellX)
		newValue := strings.ToUpper(currentValue)
		if err := ctx.ModelVC().SetCellValue(cellY, cellX, newValue); err != nil {
			return err
		}

		return nil
	})

	cm.Define("each-row", "Executes the command for each row in the column", "", func(ctx *CommandContext) error {
		if len(ctx.args) != 1 {
			return errors.New("Sub-command required")
		}

		grid := ctx.Frame().Grid()
		rows, _ := ctx.ModelVC().Model().Dimensions()

		cellX, cellY := grid.CellPosition()
		defer grid.MoveTo(cellX, cellY)

		subCommand := ctx.args

		for r := 0; r < rows; r++ {
			grid.MoveTo(cellX, r)

			if err := ctx.Session().Commands.Invoke(ctx, subCommand[0], subCommand[1:]); err != nil {
				return fmt.Errorf("at [%d, %d]: %v", cellX, r, err)
			}
		}

		return nil
	})

	cm.Define("save", "Save current file", "", func(ctx *CommandContext) error {
		var source ModelSource
		if len(ctx.args) >= 2 {
			targetCodecName := ctx.args[0]
			codecBuilder, hasCodec := codecModelSourceBuilders[targetCodecName]
			if !hasCodec {
				return fmt.Errorf("unrecognsed codec: %v", targetCodecName)
			}

			targetFilename := ctx.args[1]
			source = codecBuilder(targetFilename)
		} else {
			source = ctx.Session().Source
		}

		wSource, isWSource := source.(WritableModelSource)
		if !isWSource {
			return fmt.Errorf("model is not writable")
		}

		if err := wSource.Write(ctx.ModelVC().Model()); err != nil {
			return err
		}

		ctx.Frame().Message("Wrote " + wSource.String())
		ctx.Session().Source = wSource
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

	cm.Define("enter-zygo", "Evaluate a Zygo Lisp expression", "", func(ctx *CommandContext) error {
		ctx.Frame().Prompt(PromptOptions{Prompt: "% "}, func(expr string) error {
			res, err := ctx.session.extContext.evalCmd(expr)
			if err != nil {
				return err
			}
			ctx.Frame().ShowMessage(res)
			return nil
		})
		return nil
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

	cm.MapKey('O', cm.Command("open-right"))
	cm.MapKey('D', cm.Command("delete-row"))

	cm.MapKey('/', cm.Command("search"))
	cm.MapKey('n', cm.Command("search-next"))

	cm.MapKey('y', cm.Command("yank"))
	cm.MapKey('p', cm.Command("paste"))

	cm.MapKey('0', cm.Command("clear-row-marker"))
	cm.MapKey('1', cm.Command("mark-row-red"))
	cm.MapKey('2', cm.Command("mark-row-green"))
	cm.MapKey('3', cm.Command("mark-row-blue"))

	cm.MapKey('{', cm.Command("dec-col-width"))
	cm.MapKey('}', cm.Command("inc-col-width"))

	cm.MapKey(':', cm.Command("enter-command"))
	cm.MapKey('%', cm.Command("enter-zygo"))
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
