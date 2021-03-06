package main

import (
	"flag"
	keyboard "github.com/jteeuwen/keyboard/termbox"
	"github.com/nirasan/go-solitaire/klondike"
	"github.com/nirasan/go-solitaire/klondike/renderer"
	"github.com/nsf/termbox-go"
)

type Renderer interface {
	Render()
	SetError(error)
}

var (
	k              *klondike.Klondike
	rendererString = flag.String("mode", "basic", "execute mode \"basic\" or \"ls\"")
	colorFlag      = flag.Bool("color", true, "render color")
	debugFlag      = flag.Bool("debug", false, "show debug")
	r              Renderer
)

func main() {
	flag.Parse()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	k = klondike.NewKlondike()
	k.Init()

	switch *rendererString {
	case "basic":
		r = renderer.NewBasicRenderer(k, *colorFlag, *debugFlag)
	case "ls":
		r = renderer.NewLsRenderer(k, *colorFlag, *debugFlag)
	case "simple":
		r = renderer.NewSimpleRenderer(k, *colorFlag)
	}

	draw()

	pollEvent()
}

func draw() {
	r.Render()
}

func pollEvent() {
	running := true

	kb := keyboard.New()
	kb.Bind(func() { running = false }, "escape", "q")
	kb.Bind(func() { k.CursorUp(); draw() }, "up", "k")
	kb.Bind(func() { k.CursorDown(); draw() }, "down", "j")
	kb.Bind(func() { k.CursorLeft(); draw() }, "left", "h")
	kb.Bind(func() { k.CursorRight(); draw() }, "right", "l")
	kb.Bind(func() { k.CursorJump(); draw() }, "tab")
	kb.Bind(func() { r.SetError(k.Select()); draw() }, "space")

	for running {
		kb.Poll(termbox.PollEvent())
	}
}
