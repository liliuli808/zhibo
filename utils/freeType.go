package utils

import (
	"bufio"
	"flag"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofrs/uuid"
)

var (
	dpi      = flag.Float64("dpi", 300, "screen resolution in Dots Per Inch")
	fontFile = flag.String("fontFile", "./ttl/hua.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	white    = flag.Bool("white", false, "white text on a black background")
)

func TextToImage(text []string) string {
	flag.Parse()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontFile)
	if err != nil {
		log.Println(err)
		return ""
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{R: 0xdd, G: 0xdd, B: 0xdd, A: 0xff}
	if *white {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 2400, 1080))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return ""
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	u2, err := uuid.NewV4()
	imagePath := "./images/" + u2.String() + ".png"
	// Save that RGBA image to disk.
	outFile, err := os.Create(imagePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return imagePath
}
