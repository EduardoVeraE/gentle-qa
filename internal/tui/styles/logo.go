package styles

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"strings"

	"golang.org/x/image/draw"
)

//go:embed cosmonaut.png
var cosmonautPNG []byte

// logoWidth is the output width in terminal columns (one pixel per column).
const logoWidth = 60

// RenderLogo returns the Gentle-QA cosmonaut rendered with Unicode half-block
// characters (▀) and ANSI TrueColor escape codes.
//
// Each character cell covers two pixel rows: the upper pixel sets the
// foreground color and the lower pixel sets the background color.
// A reset sequence (\x1b[0m) is emitted at the end of every line.
func RenderLogo() string {
	src, _, err := image.Decode(bytes.NewReader(cosmonautPNG))
	if err != nil {
		return ""
	}

	srcBounds := src.Bounds()
	targetW := logoWidth
	targetH := srcBounds.Dy() * targetW / srcBounds.Dx()
	if targetH == 0 {
		return ""
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	var b strings.Builder

	for y := 0; y < targetH; y += 2 {
		for x := 0; x < targetW; x++ {
			top := dst.RGBAAt(x, y)
			var br, bg, bb uint8
			if y+1 < targetH {
				bot := dst.RGBAAt(x, y+1)
				br, bg, bb = bot.R, bot.G, bot.B
			}
			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				top.R, top.G, top.B,
				br, bg, bb)
		}
		b.WriteString("\x1b[0m\n")
	}

	// Brand line: visor face (Teal #9ccfd8) + product name (Lavender #c4a7e7)
	b.WriteString("\x1b[38;2;156;207;216m\x1b[1m  [*_*]\x1b[0m")
	b.WriteString("\x1b[38;2;196;167;231m\x1b[1m  Gentle-QA\x1b[0m")

	return b.String()
}
