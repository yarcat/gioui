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

	slider := newSlider()

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.W.Layout(gtx, slider.Layout)
			if slider.Range.Changed() {
				log.Println(slider.Range)
			}
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

func newSlider() xw.SliderStyle {
	values := []float32{-0.5, 0, 0.5, 1}
	red := color.NRGBA{R: 0xff, A: 0xff}
	yellow := color.NRGBA{R: 0xfd, G: 0xa5, B: 0x0f, A: 0xff}
	green := color.NRGBA{G: 0xff, A: 0xff}
	trackColors := []color.NRGBA{
		red,
		yellow,
		green,
		yellow,
		red,
	}
	thumbColors := []color.NRGBA{
		red, yellow, yellow, red,
	}
	return xw.SliderStyle{
		ThumbRadius: unit.Dp(12),
		TrackWidth:  unit.Dp(5),
		FingerSize:  unit.Dp(20),
		TrackColor: func(i int) color.NRGBA {
			return trackColors[i]
		},
		ThumbColor: func(i int) color.NRGBA {
			return thumbColors[i]
		},
		Min:   -1,
		Max:   2,
		Range: &xw.Range{Values: values},
	}
}
