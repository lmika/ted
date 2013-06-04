GO = go

GOFILES = main.go ui.go grid.go

ted: $(GOFILES)
	$(GO) build -o ted $(GOFILES)

clean:
	-rm ted

dependencies:
	$(GO) get "github.com/nsf/termbox-go"
