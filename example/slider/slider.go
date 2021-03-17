package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

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
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	slider := xw.SliderStyle{
		Shaper:      th.Shaper,
		ThumbRadius: unit.Dp(12),
		TrackWidth:  unit.Dp(5),
		InColor:     color.NRGBA{R: 0xff, A: 0xff},
		OutColor:    color.NRGBA{R: 0xff, A: 0x7f},
		Min:         0,
		Max:         1,
	}
	slider.Float.Min = new(widget.Float)
	slider.Float.Max = new(widget.Float)

	slider.Float.Min.Value = 0.25
	slider.Float.Max.Value = 0.75

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.W.Layout(gtx, slider.Layout)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}
