package main

import "github.com/glycerine/zygomys/v6/zygo"

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

	env.AddFunction("cursor_value", ec.builtinCursorValue)

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

func (ec *extContext) builtinCursorValue(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	grid := ec.session.Frame.Grid()
	cx, cy := grid.CellPosition()

	val := ec.session.model.CellValue(cy, cx)

	return &zygo.SexpStr{S: val}, nil
}
