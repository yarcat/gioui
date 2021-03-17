package widget

import (
	"image"
	"image/color"
	"log"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
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

	Float struct {
		Min *widget.Float
		Max *widget.Float
	}
}

func (ss SliderStyle) Layout(gtx layout.Context) layout.Dimensions {
	thumbRadius := gtx.Px(ss.ThumbRadius)

	tw := gtx.Px(ss.TrackWidth)

	// size := image.Pt(gtx.Constraints.Max.X, r*2)

	fs := sliderStyle{
		Min:      ss.Min,
		Max:      ss.Max,
		MidValue: (ss.Float.Min.Value + ss.Float.Max.Value) / 2,
		Radius:   thumbRadius,
	}
	fs.Layout(gtx, ss.Float.Min)
	fs.Layout(gtx, ss.Float.Max)

	v1, v2 := ss.Float.Min.Value, ss.Float.Max.Value
	log.Println(v1, v2)

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

type sliderStyle struct {
	Min, Max float32
	MidValue float32
	Radius   int
}

func (ss sliderStyle) Layout(gtx layout.Context, f *widget.Float) {
	st := op.Save(gtx.Ops)
	mid := float32(gtx.Constraints.Max.X) * ss.MidValue
	gtx.Constraints.Min.Y = ss.Radius
	var min, max float32
	if f.Value >= ss.MidValue {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X - int(mid)
		op.Offset(f32.Pt(mid, 0)).Add(gtx.Ops)
		paint.FillShape(gtx.Ops, color.NRGBA{G: 0xff, A: 0xff},
			clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 20)}.Op())
		min, max = ss.MidValue, ss.Max
	} else {
		gtx.Constraints.Min.X = int(mid)
		paint.FillShape(gtx.Ops, color.NRGBA{B: 0xff, A: 0xff},
			clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 20)}.Op())
		min, max = 0, ss.MidValue
	}
	f.Layout(gtx, ss.Radius, min, max)
	st.Load()
}
