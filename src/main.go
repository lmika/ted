package main

import (
    "fmt"
)

func main() {
    uiCtx, _ := NewUI()

    uiCtx.Redraw()
    uiCtx.NextEvent()    

    uiCtx.Close()
    fmt.Printf("OK!")
}
