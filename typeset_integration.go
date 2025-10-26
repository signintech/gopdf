package gopdf

import (
	"bytes"

	gotextfont "github.com/go-text/typesetting/font"
	hb "github.com/go-text/typesetting/harfbuzz"
)

// ensureHB initializes cached go-text/font Face and harfbuzz Font for this subset font.
func (s *SubsetFontObj) ensureHB() error {
	if s.hbFont != nil && s.hbFace != nil {
		return nil
	}
	face, err := gotextfont.ParseTTF(bytes.NewReader(s.ttfp.FontData()))
	if err != nil {
		return err
	}
	s.hbFace = face
	s.hbFont = hb.NewFont(face)
	// Use default HarfBuzz scaling (matches face units)
	return nil
}

// AddShapedGlyph records an extra shaped glyph so it is embedded and mapped.
func (s *SubsetFontObj) AddShapedGlyph(gid uint, r rune) {
	if s.extraGlyphs == nil {
		s.extraGlyphs = make(map[uint]rune)
	}
	// Keep the first representative rune we see for this glyph id.
	if _, ok := s.extraGlyphs[gid]; !ok {
		s.extraGlyphs[gid] = r
	}
}

// ExtraGlyphs returns the additional shaped glyphs used by the document.
func (s *SubsetFontObj) ExtraGlyphs() map[uint]rune { return s.extraGlyphs }

// shapeTextMetrics shapes the provided text and returns:
// - glyph ids
// - per-glyph X advances in HarfBuzz scaled font units (same units as XOffset/YOffset)
// - per-glyph X offsets in HarfBuzz 26.6 fixed-point units (NOT shifted)
// - per-glyph Y offsets in HarfBuzz 26.6 fixed-point units (NOT shifted)
func (s *SubsetFontObj) shapeTextMetrics(text string, fontSize float64, charSpacing float64) ([]uint, []int, []int, []int, error) {
	if err := s.ensureHB(); err != nil {
		return nil, nil, nil, nil, err
	}
	runes := []rune(text)
	buf := hb.NewBuffer()
	buf.AddRunes(runes, 0, -1)
	buf.GuessSegmentProperties()

	// Let HarfBuzz choose all default features for the script/language.
	// This preserves correct mark positioning and reordering for complex scripts.
	buf.Shape(s.hbFont, nil)

	info := buf.Info
	pos := buf.Pos
	glyphs := make([]uint, len(info))
	adv := make([]int, len(info))
	xoffs := make([]int, len(info))
	yoffs := make([]int, len(info))

	for i := range info {
		gid := uint(info[i].Glyph)
		glyphs[i] = gid
		// record extra glyphs for subsetting and ToUnicode mapping
		cl := info[i].Cluster
		var rr rune
		if cl >= 0 && cl < len(runes) {
			rr = runes[cl]
		}
		s.AddShapedGlyph(gid, rr)

		// Preserve XAdvance as returned by HarfBuzz (same units as XOffset/YOffset)
		adv[i] = int(pos[i].XAdvance)
		// Preserve offsets in raw 26.6 fixed-point for higher precision downstream
		xoffs[i] = int(pos[i].XOffset)
		yoffs[i] = int(pos[i].YOffset)
	}
	return glyphs, adv, xoffs, yoffs, nil
}
