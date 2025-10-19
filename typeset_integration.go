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
// - per-glyph X advances in TTF units (units-per-EM)
// - per-glyph X offsets in TTF units (units-per-EM)
func (s *SubsetFontObj) shapeTextMetrics(text string, fontSize float64, charSpacing float64) ([]uint, []int, []int, error) {
	if err := s.ensureHB(); err != nil {
		return nil, nil, nil, err
	}
	runes := []rune(text)
	buf := hb.NewBuffer()
	buf.AddRunes(runes, 0, -1)
	buf.GuessSegmentProperties()

	// Disable discretionary/common ligatures to preserve 1:1 mapping for ToUnicode.
	var feats []hb.Feature
	if f, err := hb.ParseFeature("liga=0"); err == nil {
		feats = append(feats, f)
	}
	if f, err := hb.ParseFeature("clig=0"); err == nil {
		feats = append(feats, f)
	}
	buf.Shape(s.hbFont, feats)

	info := buf.Info
	pos := buf.Pos
	glyphs := make([]uint, len(info))
	adv := make([]int, len(info))
	xoffs := make([]int, len(info))

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

		// Convert from 26.6 fixed-point to TTF units by >> 6
		adv[i] = int(pos[i].XAdvance >> 6)
		xoffs[i] = int(pos[i].XOffset >> 6)
	}
	return glyphs, adv, xoffs, nil
}
