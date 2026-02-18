package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceOnlyOneFrameAtATime(t *testing.T) {
	seq := Sequence{
		Children: []Widget{
			Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 6, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 9, Height: 3, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// Frame 0
	im := PaintWidget(seq, image.Rect(0, 0, 10, 3), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr",
		"rrr",
		"rrr",
	}, im))

	// Frame 1
	im = PaintWidget(seq, image.Rect(0, 0, 10, 3), 1)
	assert.Equal(t, nil, checkImage([]string{
		"gggggg",
		"gggggg",
		"gggggg",
	}, im))

	// Frame 2
	im = PaintWidget(seq, image.Rect(0, 0, 10, 3), 2)
	assert.Equal(t, nil, checkImage([]string{
		"bbbbbbbbb",
		"bbbbbbbbb",
		"bbbbbbbbb",
	}, im))
}

func TestSequenceWithAnimatedChildren(t *testing.T) {
	// Returns a 2x2 grid with background color, and a single
	// pixel of foreground color at x,y.
	frame := func(x, y int, fg color.Color, bg color.Color) Widget {
		row0 := Row{
			Children: []Widget{
				Box{Width: 1, Height: 1, Color: bg},
				Box{Width: 1, Height: 1, Color: bg},
			},
		}
		row1 := Row{
			Children: []Widget{
				Box{Width: 1, Height: 1, Color: bg},
				Box{Width: 1, Height: 1, Color: bg},
			},
		}
		if y == 0 {
			row0.Children[x] = Box{Width: 1, Height: 1, Color: fg}
		} else {
			row1.Children[x] = Box{Width: 1, Height: 1, Color: fg}
		}
		return Column{
			Children: []Widget{row0, row1},
		}
	}

	black := color.RGBA{0, 0, 0, 0}
	red := color.RGBA{0xff, 0, 0, 0xff}
	green := color.RGBA{0, 0xff, 0, 0xff}
	blue := color.RGBA{0, 0, 0xff, 0xff}

	anim0 := Animation{
		Children: []Widget{
			frame(0, 0, red, black),
			frame(1, 0, red, black),
			frame(1, 1, red, black),
			frame(0, 1, red, black),
		},
	}
	anim1 := Animation{
		Children: []Widget{
			frame(0, 0, green, black),
			frame(1, 0, green, black),
			frame(1, 1, green, black),
			frame(0, 1, green, black),
		},
	}
	anim2 := Animation{
		Children: []Widget{
			frame(0, 0, blue, black),
			frame(1, 0, blue, black),
			frame(1, 1, blue, black),
			frame(0, 1, blue, black),
		},
	}

	seq := Sequence{
		Children: []Widget{
			anim0,
			anim1,
			anim2,
		},
	}

	assert.Equal(t, 12, seq.FrameCount())

	expected := [][]string{
		{
			"r.",
			"..",
		},
		{
			".r",
			"..",
		},
		{
			"..",
			".r",
		},
		{
			"..",
			"r.",
		},
		{
			"g.",
			"..",
		},
		{
			".g",
			"..",
		},
		{
			"..",
			".g",
		},
		{
			"..",
			"g.",
		},
		{
			"b.",
			"..",
		},
		{
			".b",
			"..",
		},
		{
			"..",
			".b",
		},
		{
			"..",
			"b.",
		},
	}

	for i := 0; i < seq.FrameCount(); i++ {
		im := PaintWidget(seq, image.Rect(0, 0, 2, 2), i)
		assert.Equal(t, nil, checkImage(expected[i], im))
	}
}

func TestSequenceDurationZeroIsDefault(t *testing.T) {
	seq := Sequence{
		Children: []Widget{
			Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 6, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
		},
		Duration: 0,
	}

	assert.Equal(t, 2, seq.FrameCount())

	im := PaintWidget(seq, image.Rect(0, 0, 10, 3), 0)
	assert.Equal(t, nil, checkImage([]string{"rrr", "rrr", "rrr"}, im))

	im = PaintWidget(seq, image.Rect(0, 0, 10, 3), 1)
	assert.Equal(t, nil, checkImage([]string{"gggggg", "gggggg", "gggggg"}, im))
}

func TestSequenceDurationLoops(t *testing.T) {
	red := color.RGBA{0xff, 0, 0, 0xff}
	green := color.RGBA{0, 0xff, 0, 0xff}
	blue := color.RGBA{0, 0, 0xff, 0xff}

	seq := Sequence{
		Children: []Widget{
			Box{Width: 3, Height: 3, Color: red},
			Box{Width: 3, Height: 3, Color: green},
			Box{Width: 3, Height: 3, Color: blue},
		},
		Duration: 7, // base cycle = 3, so: r g b r g b r
	}

	assert.Equal(t, 7, seq.FrameCount())
	assert.Equal(t, 3, seq.baseCycleFrameCount())

	// First cycle
	im := PaintWidget(seq, image.Rect(0, 0, 3, 3), 0)
	assert.Nil(t, checkImage([]string{"rrr", "rrr", "rrr"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 1)
	assert.Nil(t, checkImage([]string{"ggg", "ggg", "ggg"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 2)
	assert.Nil(t, checkImage([]string{"bbb", "bbb", "bbb"}, im))

	// Second cycle (looped)
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 3)
	assert.Nil(t, checkImage([]string{"rrr", "rrr", "rrr"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 4)
	assert.Nil(t, checkImage([]string{"ggg", "ggg", "ggg"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 5)
	assert.Nil(t, checkImage([]string{"bbb", "bbb", "bbb"}, im))

	// Partial third cycle
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 6)
	assert.Nil(t, checkImage([]string{"rrr", "rrr", "rrr"}, im))
}

func TestSequenceDurationWithAnimatedChildren(t *testing.T) {
	red := color.RGBA{0xff, 0, 0, 0xff}
	green := color.RGBA{0, 0xff, 0, 0xff}

	anim0 := Animation{
		Children: []Widget{
			Box{Width: 2, Height: 2, Color: red},
			Box{Width: 2, Height: 2, Color: red},
		},
	}
	anim1 := Animation{
		Children: []Widget{
			Box{Width: 2, Height: 2, Color: green},
			Box{Width: 2, Height: 2, Color: green},
		},
	}

	seq := Sequence{
		Children: []Widget{anim0, anim1},
		Duration: 6, // base cycle = 4 (2+2), loops into partial second cycle
	}

	assert.Equal(t, 6, seq.FrameCount())
	assert.Equal(t, 4, seq.baseCycleFrameCount())

	// Frame 0,1 => anim0 (red)
	im := PaintWidget(seq, image.Rect(0, 0, 2, 2), 0)
	assert.Nil(t, checkImage([]string{"rr", "rr"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 2, 2), 1)
	assert.Nil(t, checkImage([]string{"rr", "rr"}, im))

	// Frame 2,3 => anim1 (green)
	im = PaintWidget(seq, image.Rect(0, 0, 2, 2), 2)
	assert.Nil(t, checkImage([]string{"gg", "gg"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 2, 2), 3)
	assert.Nil(t, checkImage([]string{"gg", "gg"}, im))

	// Frame 4,5 => looped back: anim0 (red)
	im = PaintWidget(seq, image.Rect(0, 0, 2, 2), 4)
	assert.Nil(t, checkImage([]string{"rr", "rr"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 2, 2), 5)
	assert.Nil(t, checkImage([]string{"rr", "rr"}, im))
}

func TestSequenceDurationShorterThanBase(t *testing.T) {
	red := color.RGBA{0xff, 0, 0, 0xff}
	green := color.RGBA{0, 0xff, 0, 0xff}

	seq := Sequence{
		Children: []Widget{
			Box{Width: 3, Height: 3, Color: red},
			Box{Width: 3, Height: 3, Color: green},
			Box{Width: 3, Height: 3, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
		Duration: 2, // truncated: only red and green
	}

	assert.Equal(t, 2, seq.FrameCount())

	im := PaintWidget(seq, image.Rect(0, 0, 3, 3), 0)
	assert.Nil(t, checkImage([]string{"rrr", "rrr", "rrr"}, im))
	im = PaintWidget(seq, image.Rect(0, 0, 3, 3), 1)
	assert.Nil(t, checkImage([]string{"ggg", "ggg", "ggg"}, im))
}
