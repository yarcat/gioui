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

func (f *Range) updateFromEvent(
	evt *pointer.Event, thumbRadius int, length float32,
	min, max float32,
) {
	if evt == nil {
		if !f.drag.Dragging() {
			f.action = rangeActionNone
		}
		return
	}
	pos := (evt.Position.X/length)*(max-min) + min
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
		f.setRange(pos, f.Max, min, max)
	case rangeActionDraggingMax:
		f.setRange(f.Min, pos, min, max)
	case rangeActionDraggingBoth:
		dpos := pos - f.pos
		f.pos = pos
		f.setRange(f.Min+dpos, f.Max+dpos, min, max)
	}
}

func (f *Range) setRange(valMin, valMax, rangeMin, rangeMax float32) {
	if valMin < rangeMin {
		valMin = rangeMin
	}
	if valMin > f.Max {
		valMin = f.Max
	}
	if valMax > rangeMax {
		valMax = rangeMax
	}
	if valMax < f.Min {
		valMax = f.Min
	}
	if valMin != f.Min || valMax != f.Max {
		f.Min, f.Max = valMin, valMax
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

	f.updateFromEvent(de, thumbRadius, length, min, max)

	defer op.Save(gtx.Ops).Load()
	pointer.Rect(image.Rectangle{Max: size}).Add(gtx.Ops)
	f.drag.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
