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


    cmdText := &ui.TextEntry{ Prompt: "Enter: " }

    statusLayout := &ui.VertLinearLayout{}
    statusLayout.Append(&ui.StatusBar{"Test", "Component"})
    statusLayout.Append(cmdText)
    //statusLayout.Append(&ui.StatusBar{"Another", "Component"})
    //statusLayout.Append(&ui.StatusBar{"Third", "Test"})

    grid := ui.NewGrid(&ui.TestModel{})

    clientArea := &ui.RelativeLayout{ Client: grid, South: statusLayout }

    uiManager.SetRootComponent(clientArea)
    //uiManager.SetFocusedComponent(grid)
    uiManager.SetFocusedComponent(cmdText)

    uiManager.Loop()
    /*
    uiCtx, _ := NewUI()

    uiCtx.Redraw()
    uiCtx.NextEvent()    

    uiCtx.Close()
    fmt.Printf("OK!")
    */
}
