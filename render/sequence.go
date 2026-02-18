package render

import (
	"image"

	"github.com/tidbyt/gg"
)

// Sequence renders a list of child widgets in sequence.
//
// Each child widget is rendered for the duration of its
// frame count, then the next child wiget in the list will
// be rendered and so on.
//
// It comes in quite useful when chaining animations.
// If you want to know more about that, go check
// out the [animation](animation.md) documentation.
//
// When Duration is set to a value greater than zero, the
// Sequence will report that as its FrameCount and loop its
// children to fill the requested number of frames.
//
// DOC(Children): List of child widgets
// DOC(Duration): Total duration in frames (0 = play once, >0 = loop to fill)
//
// EXAMPLE BEGIN
// render.Sequence(
//
//	children = [
//	  animation.Transformation(...),
//	  animation.Transformation(...),
//	  ...
//	],
//
// ),
// EXAMPLE END
type Sequence struct {
	Widget

	Children []Widget `starlark:"children,required"`
	Duration int      `starlark:"duration"`
}

// baseCycleFrameCount returns the natural frame count of one
// complete pass through all children.
func (s Sequence) baseCycleFrameCount() int {
	fc := 0
	for _, c := range s.Children {
		fc += c.FrameCount()
	}
	return fc
}

func (s Sequence) FrameCount() int {
	if s.Duration > 0 {
		return s.Duration
	}
	return s.baseCycleFrameCount()
}

// resolveFrameIdx maps a frame index into the base cycle range.
// When Duration is not set, returns frameIdx unchanged (original behavior).
// When Duration > 0, wraps via modulo so children loop.
func (s Sequence) resolveFrameIdx(frameIdx int) int {
	if s.Duration <= 0 {
		return frameIdx
	}
	base := s.baseCycleFrameCount()
	if base <= 0 {
		return 0
	}
	return frameIdx % base
}

func (s Sequence) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	frameIdx = s.resolveFrameIdx(frameIdx)

	fc := 0
	for _, c := range s.Children {
		if frameIdx < fc+c.FrameCount() {
			return c.PaintBounds(bounds, frameIdx-fc)
		}
		fc += c.FrameCount()
	}

	return image.Rect(0, 0, 0, 0)
}

func (s Sequence) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	frameIdx = s.resolveFrameIdx(frameIdx)

	fc := 0
	for _, c := range s.Children {
		if frameIdx < fc+c.FrameCount() {
			dc.Push()
			c.Paint(dc, bounds, frameIdx-fc)
			dc.Pop()
			break
		}
		fc += c.FrameCount()
	}
}
