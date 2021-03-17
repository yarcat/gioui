package widget

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

// Range is for selecting a range.
type Range struct {
	Min, Max float32

	drag   gesture.Drag
	action rangeAction
	pos    float32

	changed bool
}

type rangeAction uint8

const (
	rangeActionNone rangeAction = iota
	rangeActionDraggingMin
	rangeActionDraggingMax
	rangeActionDraggingBoth
)

func (f *Range) updateFromEvent(de *pointer.Event, thumbRadius int, length float32) {
	if de == nil {
		if !f.drag.Dragging() {
			f.action = rangeActionNone
			f.pos = 0
		}
		return
	}
	pos := de.Position.X / length
	if f.action == rangeActionNone {
		d := float32(thumbRadius) / length
		if pos < f.Min+d {
			f.action = rangeActionDraggingMin
		} else if pos > f.Max-d {
			f.action = rangeActionDraggingMax
		} else {
			f.action = rangeActionDraggingBoth
			f.pos = pos
		}
	}
	switch f.action {
	case rangeActionDraggingMin:
		f.setRange(pos, f.Max)
	case rangeActionDraggingMax:
		f.setRange(f.Min, pos)
	case rangeActionDraggingBoth:
		dpos := pos - f.pos
		f.setRange(f.Min+dpos, f.Max+dpos)
		f.pos = pos
	}
}

func (f *Range) setRange(min, max float32) {
	if min != f.Min || max != f.Max {
		f.Min, f.Max = min, max
		f.changed = true
	}
}

// Changed returns whether any of min/max values were changed since the last
// method invocation.
func (f *Range) Changed() (changed bool) {
	changed, f.changed = f.changed, false
	return
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

	f.updateFromEvent(de, thumbRadius, length)

	defer op.Save(gtx.Ops).Load()
	pointer.Rect(image.Rectangle{Max: size}).Add(gtx.Ops)
	f.drag.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
