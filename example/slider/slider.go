package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	xw "github.com/yarcat/gioui/widget"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops

	slider := xw.SliderStyle{
		ThumbRadius: unit.Dp(12),
		TrackWidth:  unit.Dp(5),
		InColor:     color.NRGBA{R: 0xff, A: 0xff},
		OutColor:    color.NRGBA{R: 0xff, A: 0x7f},
		Min:         -1,
		Max:         2,
		Range:       &xw.Range{Min: 1.5, Max: 1.75},
	}

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.W.Layout(gtx, slider.Layout)
			if slider.Range.Changed() {
				log.Println(slider.Range.Min, slider.Range.Max)
			}
			e.Frame(gtx.Ops)
		}
	}
	return nil
}
