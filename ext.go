package main

import (
	"errors"

	"github.com/glycerine/zygomys/v6/zygo"
)

type extContext struct {
	session *Session
	env     *zygo.Zlisp
}

func newExtContext(session *Session) *extContext {
	env := zygo.NewZlisp()
	ec := &extContext{
		env:     env,
		session: session,
	}

	env.AddFunction("cellValue", ec.builtinCellValue)
	env.AddFunction("cellRow", ec.builtinCellRow)
	env.AddFunction("cellCol", ec.builtinCellCol)

	env.AddFunction("modelResize", ec.builtinModelResize)

	return ec
}

func (ec *extContext) evalCmd(cmd string) (string, error) {
	if err := ec.env.LoadString(cmd); err != nil {
		return "", err
	}
	res, err := ec.env.Run()
	if err != nil {
		return "", err
	}
	return res.SexpString(zygo.NewPrintState()), nil
}

func (ec *extContext) builtinCellRow(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	grid := ec.session.Frame.Grid()
	cx, cy := grid.CellPosition()

	if len(args) == 0 {
		return &zygo.SexpInt{Val: int64(cy)}, nil
	}

	newInt, ok := args[0].(*zygo.SexpInt)
	if !ok {
		return nil, errors.New("expected arg 0 to be an int")
	}

	grid.MoveTo(cx, int(newInt.Val))
	return newInt, nil
}

func (ec *extContext) builtinCellCol(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	grid := ec.session.Frame.Grid()
	cx, cy := grid.CellPosition()

	if len(args) == 0 {
		return &zygo.SexpInt{Val: int64(cx)}, nil
	}

	newInt, ok := args[0].(*zygo.SexpInt)
	if !ok {
		return nil, errors.New("expected arg 0 to be an int")
	}

	grid.MoveTo(int(newInt.Val), cy)
	return newInt, nil
}

func (ec *extContext) builtinCellValue(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	grid := ec.session.Frame.Grid()
	cx, cy := grid.CellPosition()

	if len(args) == 0 {
		val := ec.session.model.CellValue(cy, cx)
		return &zygo.SexpStr{S: val}, nil
	}

	newStr, ok := args[0].(*zygo.SexpStr)
	if !ok {
		return nil, errors.New("expected arg 0 to be a string")
	}

	ec.session.modelController.SetCellValue(cy, cx, newStr.S)
	return newStr, nil
}

func (ec *extContext) builtinModelResize(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	if len(args) != 2 {
		return nil, errors.New("expected two args")
	}

	newRow, ok := args[0].(*zygo.SexpInt)
	if !ok {
		return nil, errors.New("expected arg 0 to be an int")
	}

	newCol, ok := args[1].(*zygo.SexpInt)
	if !ok {
		return nil, errors.New("expected arg 1 to be an int")
	}

	if err := ec.session.modelController.Resize(int(newRow.Val), int(newCol.Val)); err != nil {
		return nil, err
	}
	return zygo.NullRT, nil
}
