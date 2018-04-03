package main

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"strings"
)

// Renderer ...
type Renderer struct {
	FontSize float64
	DPI      float64
	StartX   float64
	StartY   float64
	Spacing  float64
	Theme    *chroma.Style
	Font     *truetype.Font
	Context  *freetype.Context

	fontFace font.Face
}

// NewRenderer ...
func NewRenderer() Renderer {
	var (
		fontSize = 16.0
		dpi      = 120.0
		spacing  = 1.5
		startX   = fontSize
		startY   = fontSize
	)

	// init font
	fontBytes, err := Asset("assets/FiraCode-Regular.ttf")
	checkError(err)
	customFont, err := freetype.ParseFont(fontBytes)
	checkError(err)

	style := styles.Get("monokai")

	// init fontface
	face := truetype.NewFace(customFont, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	context := freetype.NewContext()
	context.SetFont(customFont)
	context.SetDPI(dpi)
	context.SetFontSize(fontSize)

	return Renderer{
		FontSize: fontSize,
		Font:     customFont,
		fontFace: face,
		Spacing:  spacing,
		DPI:      dpi,
		Theme:    style,
		StartX:   startX,
		StartY:   startY,
		Context:  context,
	}
}

// ChangeStyle to highlight
func (s *Renderer) ChangeStyle(styleName string) {
	s.Theme = styles.Get(styleName)
}

func (s *Renderer) textWidth(text string) float64 {
	return float64(font.MeasureString(s.fontFace, text).Ceil())
}

func (s *Renderer) bound(text string) (width float64, height float64) {
	lines := strings.Split(text, "\n")
	maxWidth := 0.0
	tmpWidth := 0.9
	for _, line := range lines {
		tmpWidth = s.textWidth(line)
		if maxWidth < tmpWidth {
			maxWidth = tmpWidth
		}
	}

	height = float64(len(lines)) * s.Spacing * s.FontSize

	return maxWidth, height
}

func (s *Renderer) createCanvas(w, h float64) *image.RGBA {
	styleBg := s.Theme.Get(chroma.Background).Background

	rec := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	bg := image.NewUniform(color.RGBA{
		styleBg.Red(),
		styleBg.Green(),
		styleBg.Blue(),
		255,
	})
	draw.Draw(rec, rec.Bounds(), bg, image.ZP, draw.Src)
	return rec
}

func (s *Renderer) updateContext(dst *image.RGBA) {
	s.Context.SetDst(dst)
	s.Context.SetClip(dst.Bounds())
}

func (s *Renderer) fixPoint(x float64) fixed.Int26_6 {
	return s.Context.PointToFixed(x)
}

// Render ...
func (s *Renderer) Render(contents, lang string) *image.RGBA {
	contents = strings.Replace(contents, "\t", "    ", -1)

	w, h := s.bound(contents)

	w += s.StartX * 2
	h += s.StartY * 2
	canvas := s.createCanvas(
		float64(s.fixPoint(w).Ceil()),
		float64(s.fixPoint(h).Ceil()),
	)

	s.updateContext(canvas)

	pt := freetype.Pt(
		s.fixPoint(s.StartX).Ceil(),
		s.fixPoint(s.StartY+(s.FontSize*s.Spacing)).Ceil(),
	)

	var lexer chroma.Lexer
	if strings.HasPrefix(lang, "match:") {
		lang := strings.Replace(lang, "match:", "", 1)
		lexer = lexers.Match(lang)
	} else {
		lexer = lexers.Get(lang)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, contents)
	checkError(err)
	tokens := iterator.Tokens()

	appendX := 0.0
	for _, token := range tokens {
		value := token.Value
		if value == "" {
			continue
		}

		tColor := s.Theme.Get(token.Type).Colour
		s.Context.SetSrc(image.NewUniform(color.RGBA{
			tColor.Red(),
			tColor.Green(),
			tColor.Blue(),
			255,
		}))

		lines := strings.Split(value, "\n")

		for i, line := range lines {
			if i > 0 {
				appendX = 0
				y := (s.FontSize * s.Spacing)
				pt.Y += s.fixPoint(y)
			}

			pt.X = s.fixPoint(s.StartX + appendX)

			if line != " " {
				s.Context.DrawString(line, pt)
			}

			appendX += s.textWidth(line)
		}
	}

	return canvas
}
