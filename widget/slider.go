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

	InColor  color.NRGBA
	OutColor color.NRGBA

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
	op.Offset(f32.Pt(float32(thumbRadius), 0)).Add(gtx.Ops)

	gtx.Constraints.Min = image.Pt(width-2*thumbRadius, br.Dy())
	ss.Range.Layout(gtx, thumbRadius, ss.Min, ss.Max)

	// Both values are now always in [0..1] range.
	v1 := (ss.Range.Min - ss.Min) / (ss.Max - ss.Min)
	v2 := (ss.Range.Max - ss.Min) / (ss.Max - ss.Min)

	op.Offset(f32.Pt(0, float32(height)/2)).Add(gtx.Ops)

	tr := image.Rect(0, -trackWidth/2, width-2*thumbRadius, trackWidth/2)
	// Draw track before the first thumb.
	drawTrack(gtx.Ops, ss.OutColor, tr, 0, v1)
	// Draw track before the second thumb.
	drawTrack(gtx.Ops, ss.InColor, tr, v1, v2)
	// Draw track after the second thumb.
	drawTrack(gtx.Ops, ss.OutColor, tr, v2, 1)
	// Draw the first thumb.
	drawThumb(gtx.Ops, ss.InColor, tr, float32(thumbRadius), v1)
	// Draw the second thumb.
	drawThumb(gtx.Ops, ss.InColor, tr, float32(thumbRadius), v2)

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
