package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarqueeNoScrollHorizontal(t *testing.T) {
	m := Marquee{
		Width: 6,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
	}

	mv := Marquee{
		Height: 3,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		ScrollDirection: "vertical",
	}

	// Child fits so there's just 1 single frame
	assert.Equal(t, 1, m.FrameCount())
	assert.Equal(t, 1, mv.FrameCount())
	im := PaintWidget(m, image.Rect(0, 0, 100, 100), 0)
	imv := PaintWidget(mv, image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrggb",
		"rrrgg.",
		"rrr...",
	}, im))
	assert.Equal(t, nil, checkImage([]string{
		"rrrggb",
		"rrrgg.",
		"rrr...",
	}, imv))
}

func TestMarqueeNoScrollAlignCenter(t *testing.T) {
	m := Marquee{
		Width: 8,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Align: "center",
	}

	mv := Marquee{
		Height: 5,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		ScrollDirection: "vertical",
		Align:           "center",
	}

	// Child fits so there's just 1 single frame
	assert.Equal(t, 1, m.FrameCount())
	assert.Equal(t, 1, mv.FrameCount())
	im := PaintWidget(m, image.Rect(0, 0, 100, 100), 0)
	imv := PaintWidget(mv, image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rrrggb.",
		".rrrgg..",
		".rrr....",
	}, im))
	assert.Equal(t, nil, checkImage([]string{
		"......",
		"rrrggb",
		"rrrgg.",
		"rrr...",
		"......",
	}, imv))
}

func TestMarqueeNoScrollAlignEnd(t *testing.T) {
	m := Marquee{
		Width: 8,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Align: "end",
	}

	mv := Marquee{
		Height: 5,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		ScrollDirection: "vertical",
		Align:           "end",
	}

	// Child fits so there's just 1 single frame
	assert.Equal(t, 1, m.FrameCount())
	assert.Equal(t, 1, mv.FrameCount())
	im := PaintWidget(m, image.Rect(0, 0, 100, 100), 0)
	imv := PaintWidget(mv, image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrrggb",
		"..rrrgg.",
		"..rrr...",
	}, im))
	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"rrrggb",
		"rrrgg.",
		"rrr...",
	}, imv))
}

// The addition of OffsetStart and OffsetEnd changes the default
// behaviour of Marquee. Passing start==width and end==0 mimics the
// old default.
func TestMarqueeOldBehavior(t *testing.T) {
	m := Marquee{
		Width:       6,
		OffsetStart: 6,
		OffsetEnd:   0,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 3, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 3, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
	}

	// The child's 9 pixels will be scrolled into view (7 frames),
	// scrolled out of view (9 frames) and then finally scrolled
	// back into view again (6 frames). 22 frames in total.
	assert.Equal(t, 22, m.FrameCount())

	// Scrolling into view
	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"......",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 0)))

	assert.Equal(t, nil, checkImage([]string{
		"....rr",
		"....rr",
		"....rr",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 2)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 6)))

	// Scrolling out of view
	assert.Equal(t, nil, checkImage([]string{
		"rgggbb",
		"rggg..",
		"r.....",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 8)))

	assert.Equal(t, nil, checkImage([]string{
		"b.....",
		"......",
		"......",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 14)))

	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"......",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 15)))

	// Scrolling back into view
	assert.Equal(t, nil, checkImage([]string{
		"...rrr",
		"...rrr",
		"...rrr",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 18)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 21)))

	// Later frames keep it fixed in the last frame. This makes
	// multiple simultaneous marquees look nice when they've
	// different length.

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 22)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 26)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, PaintWidget(m, image.Rect(0, 0, 100, 100), 100000)))
}

func TestMarqueeOffsetIdentical(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Width: 6,
		Child: child,
	}
	im := image.Rect(0, 0, 100, 100)

	// Check that identical frames are not repeated after
	// another, if start and end offset are identical.
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 12)))
	assert.Equal(t, 13, m.FrameCount())

	m.OffsetStart = 3
	m.OffsetEnd = 3
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 12)))
	assert.Equal(t, 13, m.FrameCount())
}

func TestMarqueeOffsetStart(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Width: 6,
		Child: child,
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetStart affects the initial position of the child
	m.OffsetStart = 2
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 9)))

	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 15)))
	assert.Equal(t, 16, m.FrameCount())

	// Negative OffsetStart
	m.OffsetStart = -2
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 11)))
	assert.Equal(t, 12, m.FrameCount())

	// Overly negative OffsetStart is truncated to child width
	m.OffsetStart = -1000
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 2)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -7
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -8
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -6
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 2)))
	assert.Equal(t, 8, m.FrameCount())
}

func TestMarqueeOffsetEnd(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Width: 6,
		Child: child,
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetEnd affects the final position of the child
	m.OffsetEnd = 2
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 11)))
	assert.Equal(t, 12, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 1024)))

	// Negative offset places child outside of marquee
	m.OffsetEnd = -4
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 15)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 16)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 17)))
	assert.Equal(t, 18, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 18)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 19)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 1024)))

	// Very negative offset is truncated to width of child
	m.OffsetEnd = -133
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 17)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 18)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 19)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 20)))
	assert.Equal(t, 21, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 21)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 22)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 23)))

	// OffsetEnd >= width means it doesn't scroll back
	m.OffsetEnd = 6
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 7)))
	assert.Equal(t, 8, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 1024)))

}

func TestMarqueeDelayScrollOffsetStart(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Width: 6,
		Child: child,
		Delay: 2,
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetStart affects the initial position of the child
	m.OffsetStart = 2
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 11)))

	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 15)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 16)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 17)))
	assert.Equal(t, 18, m.FrameCount())

	// // Negative OffsetStart
	m.OffsetStart = -2
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 13)))
	assert.Equal(t, 14, m.FrameCount())

	// // Overly negative OffsetStart is truncated to child width
	m.OffsetStart = -1000
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 4)))
	assert.Equal(t, 9, m.FrameCount())
	m.OffsetStart = -7
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 3)))
	assert.Equal(t, 9, m.FrameCount())
	m.OffsetStart = -8
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 3)))
	assert.Equal(t, 9, m.FrameCount())
	m.OffsetStart = -6
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 4)))
	assert.Equal(t, 10, m.FrameCount())
}

func TestMarqueeVerticalScroll(t *testing.T) {
	child := Column{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 1, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 1, Height: 4, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Height:          6,
		Child:           child,
		ScrollDirection: "vertical",
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetEnd affects the final position of the child
	m.OffsetStart = 2
	assert.Equal(t, nil, checkImage([]string{
		".",
		".",
		"r",
		"g",
		"g",
		"b",
	}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{
		".",
		"r",
		"g",
		"g",
		"b",
		"b",
	}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{
		"r",
		"g",
		"g",
		"b",
		"b",
		"b",
	}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{
		"g",
		"g",
		"b",
		"b",
		"b",
		"b",
	}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{
		"g",
		"b",
		"b",
		"b",
		"b",
		".",
	}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		"b",
		"b",
		"b",
		".",
		".",
	}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		"b",
		"b",
		".",
		".",
		".",
	}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		"b",
		".",
		".",
		".",
		".",
	}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		".",
		".",
		".",
		".",
		".",
	}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{
		".",
		".",
		".",
		".",
		".",
		".",
	}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 15)))
	assert.Equal(t, 16, m.FrameCount())

	// Negative OffsetStart
	m.OffsetStart = -2
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", ".", ".", ".", "."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 11)))
	assert.Equal(t, 12, m.FrameCount())

	// Overly negative OffsetStart is truncated to child width
	m.OffsetStart = -1000
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, PaintWidget(m, im, 2)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -7
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -8
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -6
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 2)))
	assert.Equal(t, 8, m.FrameCount())

	// OffsetEnd affects the final position of the child
	m.OffsetStart = 0
	m.OffsetEnd = 2
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", ".", ".", ".", "."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 7)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 11)))
	assert.Equal(t, 12, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 1024)))

	// Negative offset places child outside of marquee
	m.OffsetEnd = -4
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, PaintWidget(m, im, 10)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 11)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, PaintWidget(m, im, 15)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, PaintWidget(m, im, 16)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 17)))
	assert.Equal(t, 18, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 18)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 19)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 1024)))

	// OffsetEnd >= width means it doesn't scroll back
	m.OffsetEnd = 6
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, PaintWidget(m, im, 6)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 7)))
	assert.Equal(t, 8, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 8)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 1024)))
}

func TestMarqueeEndDelay(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	// EndDelay adds extra frames at the end position
	m := Marquee{
		Width:    6,
		Child:    child,
		EndDelay: 3,
	}

	// Without EndDelay, FrameCount would be 13 (offstart==offend==0).
	// With EndDelay=3, FrameCount is 16.
	assert.Equal(t, 16, m.FrameCount())

	// Normal scrolling still works
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 1)))

	// Last scroll frame
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 12)))

	// EndDelay frames freeze at the end position
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 15)))

	// Beyond FrameCount still freezes
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 100)))
}

func TestMarqueeEndDelayWithDelay(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Width:    6,
		Child:    child,
		Delay:    2,
		EndDelay: 2,
	}

	// FrameCount = 13 + 2 (delay) + 2 (end_delay) = 17
	assert.Equal(t, 17, m.FrameCount())

	// Delay phase
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 2)))

	// Scrolling begins
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 3)))
}

func TestMarqueeEndDelayWithOffsetEnd(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Width:     6,
		Child:     child,
		OffsetEnd: 2,
		EndDelay:  3,
	}

	// Without EndDelay: FrameCount = 12. With EndDelay=3: 15.
	assert.Equal(t, 15, m.FrameCount())

	// Final scroll position is at offset 2
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 11)))

	// EndDelay frames freeze at offend position
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 14)))
}

func TestMarqueeLoopHorizontal(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Width: 6,
		Child: child,
		Loop:  true,
	}

	// FrameCount = delay(0) + cw(7) + gap(10) + 1 = 18
	assert.Equal(t, 18, m.FrameCount())

	// Frame 0: child at offset 0
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))

	// Scrolling left
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, PaintWidget(m, im, 6)))

	// First copy off screen, 10px gap
	assert.Equal(t, nil, checkImage([]string{"......"}, PaintWidget(m, im, 7)))

	// Second copy starts appearing from the right (after 10px gap)
	assert.Equal(t, nil, checkImage([]string{".....r"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 15)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 16)))

	// Cycle complete: back to start position
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 17)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 100)))
}

func TestMarqueeLoopNoScrollWhenFits(t *testing.T) {
	m := Marquee{
		Width: 10,
		Child: Box{Width: 5, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
		Loop:  true,
	}

	// Child fits, so only 1 frame even with Loop=true
	assert.Equal(t, 1, m.FrameCount())
}

func TestMarqueeLoopWithDelay(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Width: 6,
		Child: child,
		Loop:  true,
		Delay: 2,
	}

	// FrameCount = delay(2) + cw(7) + gap(10) + 1 = 20
	assert.Equal(t, 20, m.FrameCount())

	// Delay phase
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 2)))

	// Scroll begins
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, PaintWidget(m, im, 4)))
}

func TestMarqueeLoopWithEndDelay(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Width:    6,
		Child:    child,
		Loop:     true,
		EndDelay: 3,
	}

	// FrameCount = delay(0) + cw(7) + gap(10) + 1 + end_delay(3) = 21
	assert.Equal(t, 21, m.FrameCount())

	// Last scroll frame before cycle wraps
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 16)))

	// Cycle complete, then EndDelay freezes at start position
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 17)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 18)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 19)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 20)))
}

func TestMarqueeLoopVertical(t *testing.T) {
	child := Column{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 1, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 1, Height: 4, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Height:          6,
		Child:           child,
		ScrollDirection: "vertical",
		Loop:            true,
	}

	// FrameCount = delay(0) + cw(7) + gap(10) + 1 = 18
	assert.Equal(t, 18, m.FrameCount())

	// Frame 0: child at offset 0
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 0)))

	// Scrolling up
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, PaintWidget(m, im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, PaintWidget(m, im, 4)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", ".", ".", ".", "."}, PaintWidget(m, im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, PaintWidget(m, im, 6)))

	// Gap (10px)
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, PaintWidget(m, im, 7)))

	// Second copy wrapping in from bottom (after 10px gap)
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, PaintWidget(m, im, 12)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, PaintWidget(m, im, 13)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, PaintWidget(m, im, 14)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, PaintWidget(m, im, 15)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, PaintWidget(m, im, 16)))

	// Cycle complete: back to start
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, PaintWidget(m, im, 17)))
}

func TestMarqueeLoopWithOffsetStart(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	im := image.Rect(0, 0, 100, 100)

	m := Marquee{
		Width:       6,
		Child:       child,
		Loop:        true,
		OffsetStart: 2,
	}

	// FrameCount = delay(0) + cw(7) + gap(10) + 1 = 18
	assert.Equal(t, 18, m.FrameCount())

	// Frame 0: child at offset 2 (shifted right)
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, PaintWidget(m, im, 0)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, PaintWidget(m, im, 1)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, PaintWidget(m, im, 2)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, PaintWidget(m, im, 3)))
}
