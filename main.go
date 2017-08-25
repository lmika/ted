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

	model := &StdModel{}
	model.Resize(5, 5)

	frame := NewFrame(uiManager)
	NewSession(uiManager, frame, model)

	uiManager.SetRootComponent(frame.RootComponent())
	frame.enterMode(GridMode)

	uiManager.Loop()
}
