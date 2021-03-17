package widget

import (
	"image"
	"log"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

// Range is for selecting a range.
type Range struct {
	Min, Max float32

	drag gesture.Drag
}

// Layout updates the range accordingly to gestures.
func (f *Range) Layout(gtx layout.Context, thumbRadius int, min, max float32) layout.Dimensions {
	size := gtx.Constraints.Min
	length := float32(size.X)

	var de *pointer.Event
	for _, e := range f.drag.Events(gtx.Metric, gtx, gesture.Horizontal) {
		if e.Type == pointer.Press || e.Type == pointer.Drag {
			de = &e
		}
	}

	if de != nil {
		pos := de.Position.X / length
		if pos < f.Min {
			f.Min = pos
		} else if pos > f.Max {
			f.Max = pos
		}
	}

	defer op.Save(gtx.Ops).Load()
	log.Println(size)
	pointer.Rect(image.Rectangle{Max: size}).Add(gtx.Ops)
	f.drag.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
