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
	Values []float32

	dragIndex int

	drag   gesture.Drag
	action rangeAction
	pos    float32

	changed bool
}

type rangeAction uint8

const (
	rangeActionNone rangeAction = iota
	rangeActionDragging
	rangeActionDraggingBoth
)

func (r *Range) updateFromEvent(
	evt *pointer.Event, thumbRadius, fingerSize int, length float32,
	min, max float32,
) {
	if evt == nil {
		if !r.drag.Dragging() {
			r.action = rangeActionNone
		}
		return
	}
	pos := (evt.Position.X-float32(thumbRadius))/length*(max-min) + min
	if r.action == rangeActionNone {
		r.setAction(fingerSize, length, pos)
	}
	switch r.action {
	case rangeActionDragging:
		r.setRange(r.dragIndex, pos, min, max)
	case rangeActionDraggingBoth:
		dpos := pos - r.pos
		r.pos = pos
		r.setRange(r.dragIndex-1, r.Values[r.dragIndex-1]+dpos, min, max)
		r.setRange(r.dragIndex, r.Values[r.dragIndex]+dpos, min, max)
	}
}

func (r *Range) setAction(fingerSize int, length, pos float32) {
	d := float32(fingerSize) / length
	if pos < r.Values[0]+d {
		r.dragIndex, r.action = 0, rangeActionDragging
		return
	}
	if pos > r.Values[len(r.Values)-1]-d {
		r.dragIndex, r.action = len(r.Values)-1, rangeActionDragging
		return
	}
	for i, v := range r.Values {
		if v-d < pos && pos < v+d {
			r.dragIndex, r.action = i, rangeActionDragging
			return
		}
		if pos < v {
			r.dragIndex, r.action = i, rangeActionDraggingBoth
			r.pos = pos
			return
		}
	}
}

func (r *Range) setRange(index int, v, rangeMin, rangeMax float32) {
	switch index {
	case 0:
		if len(r.Values) > 1 {
			rangeMax = r.Values[1]
		}
	case len(r.Values) - 1:
		rangeMin = r.Values[index-1]
	default:
		rangeMin = r.Values[index-1]
		rangeMax = r.Values[index+1]
	}
	if v < rangeMin {
		v = rangeMin
	}
	if v > rangeMax {
		v = rangeMax
	}
	if v != r.Values[index] {
		r.Values[index] = v
		r.changed = true
	}
}

// Changed returns whether any of min/max values were changed since the last
// method invocation.
func (f *Range) Changed() (changed bool) {
	changed, f.changed = f.changed, false
	return
}

// Layout updates the range accordingly to gestures.
func (f *Range) Layout(gtx layout.Context, thumbRadius, fingerSize int, min, max float32) layout.Dimensions {
	size := gtx.Constraints.Min
	length := float32(size.X - 2*thumbRadius)

	var de *pointer.Event
	for _, e := range f.drag.Events(gtx.Metric, gtx, gesture.Horizontal) {
		if e.Type == pointer.Press || e.Type == pointer.Drag {
			de = &e
		}
	}

	f.updateFromEvent(de, thumbRadius, fingerSize, length, min, max)

	defer op.Save(gtx.Ops).Load()
	pointer.Rect(image.Rectangle{Max: size}).Add(gtx.Ops)
	f.drag.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
