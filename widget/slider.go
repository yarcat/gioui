package widget

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

// SliderStyle is a multi-slider widget.
type SliderStyle struct {
	Shaper text.Shaper
	Font   text.Font

	ThumbRadius unit.Value
	TrackWidth  unit.Value
	FingerSize  unit.Value

	TrackColor func(int) color.NRGBA
	ThumbColor func(int) color.NRGBA

	Min, Max float32

	Range *Range
}

func (ss SliderStyle) Layout(gtx layout.Context) layout.Dimensions {
	thumbRadius := gtx.Px(ss.ThumbRadius)
	trackWidth := gtx.Px(ss.TrackWidth)
	fingerSize := gtx.Px(ss.FingerSize)

	width := gtx.Constraints.Max.X
	height := max(thumbRadius*2, max(trackWidth, fingerSize*2))
	size := image.Pt(width, height)
	br := image.Rectangle{Max: size}

	defer op.Save(gtx.Ops).Load()
	clip.Rect(br).Add(gtx.Ops)

	gtx.Constraints.Min = image.Pt(width, br.Dy())
	ss.Range.Layout(gtx, thumbRadius, fingerSize, ss.Min, ss.Max)

	op.Offset(f32.Pt(float32(thumbRadius), float32(height)/2)).Add(gtx.Ops)

	tr := image.Rect(0, -trackWidth/2, width-2*thumbRadius, trackWidth/2)
	var prevV float32
	for i, v := range ss.Range.Values {
		v = (v - ss.Min) / (ss.Max - ss.Min)
		if v != prevV {
			// Draw track before the first thumb.
			drawTrack(gtx.Ops, ss.TrackColor(i), tr, prevV, v)
		}
		prevV = v
	}
	if prevV != 1 {
		drawTrack(gtx.Ops, ss.TrackColor(len(ss.Range.Values)), tr, prevV, 1)
	}

	for i, v := range ss.Range.Values {
		v = (v - ss.Min) / (ss.Max - ss.Min)
		// Draw the first thumb.
		drawThumb(gtx.Ops, ss.ThumbColor(i), tr, float32(thumbRadius), v)
	}

	return layout.Dimensions{Size: size}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// drawTrack draws track segment betwee a and b. Both a and b values are
// normalized to and must be in range [0, 1].
func drawTrack(ops *op.Ops, c color.NRGBA, tr image.Rectangle, a, b float32) {
	paint.FillShape(ops, c, clip.Rect{
		Min: image.Pt(int(float32(tr.Max.X)*a), tr.Min.Y),
		Max: image.Pt(int(float32(tr.Max.X)*b), tr.Max.Y),
	}.Op())
}

func drawThumb(ops *op.Ops, c color.NRGBA, tr image.Rectangle, rad, a float32) {
	paint.FillShape(ops, c,
		clip.Circle{
			Center: f32.Pt(float32(tr.Dx())*a, float32(tr.Min.Y+tr.Dy()/2)),
			Radius: rad,
		}.Op(ops))

}
