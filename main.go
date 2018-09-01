package main

import (
	"github.com/lmika/ted/ui"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: ted FILENAME")
		os.Exit(1)
	}

	uiManager, err := ui.NewUI()
	if err != nil {
		panic(err)
	}
	defer uiManager.Close()

	frame := NewFrame(uiManager)
	session := NewSession(uiManager, frame, CsvFileModelSource{flag.Arg(0)})
	session.LoadFromSource()

	uiManager.SetRootComponent(frame.RootComponent())
	frame.enterMode(GridMode)

	uiManager.Loop()
}
