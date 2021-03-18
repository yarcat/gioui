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

	tw := gtx.Px(ss.TrackWidth)

	gtx.Constraints.Min = gtx.Constraints.Max
	ss.Range.Layout(gtx, thumbRadius, ss.Min, ss.Max)

	// Both values are now always in [0..1] range.
	v1 := (ss.Range.Min - ss.Min) / (ss.Max - ss.Min)
	v2 := (ss.Range.Max - ss.Min) / (ss.Max - ss.Min)

	tr := trackRect(thumbRadius, tw, gtx.Constraints.Max.X)

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

	return layout.Dimensions{Size: image.Pt(thumbRadius*2, thumbRadius*2)}
}

func trackRect(rad, width, maxx int) image.Rectangle {
	var mid int
	if rad*2 > width {
		mid = rad
	} else {
		mid = width / 2
	}
	return image.Rect(0, mid-width/2, maxx, mid+width/2)
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
