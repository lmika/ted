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


    statusLayout := &ui.VertLinearLayout{}
    statusLayout.Append(&ui.StatusBar{"Test", "Component"})
    statusLayout.Append(&ui.StatusBar{"Another", "Component"})
    statusLayout.Append(&ui.StatusBar{"Third", "Test"})

    clientArea := &ui.RelativeLayout{ South: statusLayout }

    uiManager.SetRootComponent(clientArea)

    uiManager.Loop()
    /*
    uiCtx, _ := NewUI()

    uiCtx.Redraw()
    uiCtx.NextEvent()    

    uiCtx.Close()
    fmt.Printf("OK!")
    */
}
