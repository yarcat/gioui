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

func (f *Range) updateFromEvent(
	evt *pointer.Event, thumbRadius, fingerSize int, length float32,
	min, max float32,
) {
	if evt == nil {
		if !f.drag.Dragging() {
			f.action = rangeActionNone
		}
		return
	}
	pos := (evt.Position.X-float32(thumbRadius))/length*(max-min) + min
	if f.action == rangeActionNone {
		d := float32(fingerSize) / length

		for i, v := range []float32{} {
			if v-d < pos && pos < v+d || i == 0 {
				log.Println("dragging", i)
				f.action = rangeActionDragging
				f.dragIndex = i
				break
			} else if pos < d {
				f.action = rangeActionDraggingBoth
				f.dragIndex = i
				f.pos = pos
				break
			}
		}
	}
	switch f.action {
	case rangeActionDragging:
		f.setRange(f.dragIndex, pos, min, max)
	case rangeActionDraggingBoth:
		// dpos := pos - f.pos
		// f.pos = pos
		// f.setRange(f.dragIndex, min, max)
	}
}

func (r *Range) setRange(index int, v, rangeMin, rangeMax float32) {
	switch index {
	case 0:
		if v < rangeMin {
			v = rangeMin
		}
	case len(r.Values) - 1:
		if v > rangeMax {
			v = rangeMax
		}
	default:
		if v < r.Values[index-1] {
			v = r.Values[index-1]
		}
		if v > r.Values[index+1] {
			v = r.Values[index+1]
		}
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
