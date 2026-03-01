package main

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row, Col, Xpixel, Ypixel uint16
}

func getTermWidth() int {
	ws := &winsize{}
	syscall.Syscall(syscall.SYS_IOCTL, 1, 0x5413, uintptr(unsafe.Pointer(ws)))
	if ws.Col > 0 {
		return int(ws.Col)
	}
	return 70
}

// hslToRgb: h=0..360, s=0..1, l=0..1 → r,g,b=0..255
func hslToRgb(h, s, l float64) (int, int, int) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2
	var r1, g1, b1 float64
	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}
	clamp := func(v float64) int {
		if v > 255 {
			return 255
		}
		return int(v + 0.5)
	}
	return clamp((r1+m)*255), clamp((g1+m)*255), clamp((b1+m)*255)
}

var subBlocks = []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}

// renderProgressBar renders a full-color gradient progress bar.
// width is the total visual width including the score label.
func renderProgressBar(score, goal, width int) string {
	label := fmt.Sprintf(" %d/%d", score, goal)
	barWidth := width - len(label)
	if barWidth < 8 {
		barWidth = 8
	}

	pct := math.Min(1.0, float64(score)/float64(goal))
	totalUnits := barWidth * 8
	filledUnits := int(pct * float64(totalUnits))
	fullCells := filledUnits / 8 // number of fully filled cells

	// We use a 2-row bar: top row thin accent, main row blocks
	// Actually single row with gradient + sub-block precision + glow at edge

	buf := make([]byte, 0, barWidth*40)

	for i := 0; i < barWidth; i++ {
		// Position drives color (gradient independent of progress)
		t := float64(i) / float64(barWidth-1)
		hue := 240.0 - t*240.0 // electric blue → cyan → green → yellow → red

		cellUnits := filledUnits - i*8

		switch {
		case cellUnits >= 8:
			// Fully filled — glow effect at the leading edge
			lightness := 0.48
			if fullCells > 0 {
				switch i {
				case fullCells - 1:
					lightness = 0.78 // brightest: leading edge
				case fullCells - 2:
					lightness = 0.62
				case fullCells - 3:
					lightness = 0.54
				}
			}
			r, g, b := hslToRgb(hue, 1.0, lightness)
			// Background: very dark tint of same hue
			br, bg, bb := r/8, g/8, b/8
			buf = fmt.Appendf(buf, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm█\033[0m",
				r, g, b, br, bg, bb)

		case cellUnits > 0:
			// Partial block at progress edge — bright
			r, g, b := hslToRgb(hue, 1.0, 0.75)
			buf = fmt.Appendf(buf, "\033[38;2;%d;%d;%dm\033[48;2;10;10;16m%c\033[0m",
				r, g, b, subBlocks[cellUnits])

		default:
			// Empty — dim dark dots
			buf = fmt.Appendf(buf, "\033[38;2;28;28;40m░\033[0m")
		}
	}

	// Score label: bold score, dim goal
	buf = fmt.Appendf(buf, "\033[0m \033[1;97m%d\033[0;2m/%d\033[0m", score, goal)

	return string(buf)
}
