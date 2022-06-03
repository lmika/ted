package main

import (
	"flag"
	"fmt"
	"github.com/lmika/ted/ui"
	"os"
)

func main() {
	var flagCodec = flag.String("c", "csv", "file codec to use")
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

	codecBuilder, hasCodec := codecModelSourceBuilders[*flagCodec]
	if !hasCodec {
		fmt.Fprintf(os.Stderr, "unrecognised codec: %v", *flagCodec)
		os.Exit(1)
	}

	frame := NewFrame(uiManager)
	session := NewSession(uiManager, frame, codecBuilder(flag.Arg(0)))
	session.LoadFromSource()

	uiManager.SetRootComponent(frame.RootComponent())
	frame.enterMode(GridMode)

	uiManager.Loop()
}

type codecModelSourceBuilder func(filename string) ModelSource

var codecModelSourceBuilders = map[string]codecModelSourceBuilder{
	"csv": func(filename string) ModelSource {
		return NewCsvFileModelSource(filename, CsvFileModelSourceOptions{Comma: ','})
	},
	"tsv": func(filename string) ModelSource {
		return NewCsvFileModelSource(filename, CsvFileModelSourceOptions{Comma: '\t'})
	},
}
