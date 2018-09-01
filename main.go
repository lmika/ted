package main

import (
	"bitbucket.org/lmika/ted-v2/ui"
)

func main() {
	uiManager, err := ui.NewUI()
	if err != nil {
		panic(err)
	}
	defer uiManager.Close()

	frame := NewFrame(uiManager)
	session := NewSession(uiManager, frame, CsvFileModelSource{"test.csv"})
	session.LoadFromSource()

	uiManager.SetRootComponent(frame.RootComponent())
	frame.enterMode(GridMode)

	uiManager.Loop()
}
