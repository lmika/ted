package main

import (
    "./ui"
)

func main() {
    uiManager, err := ui.NewUI()
    if err != nil {
        panic(err)
    }
    defer uiManager.Close()

    model := &StdModel{}

    frame := NewFrame(uiManager)
    NewSession(frame, model)

    uiManager.SetRootComponent(frame.RootComponent())
    frame.EnterMode(GridMode)

    uiManager.Loop()
}
