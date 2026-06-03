package gopdf

import "strings"

// FontFallback identifies a registered TTF font that may render missing glyphs.
type FontFallback struct {
	Family string
	Style  string
}

type fontKey struct {
	family string
	style  int
}

type textRun struct {
	text           string
	fontSubset     *SubsetFontObj
	fontCountIndex int
}

type textLayout struct {
	runs []textRun
}

func (l textLayout) plainText() string {
	var b strings.Builder
	for _, run := range l.runs {
		b.WriteString(run.text)
	}
	return b.String()
}

// SetFontFallback configures ordered fallback fonts for a primary font.
func (gp *GoPdf) SetFontFallback(primaryFamily string, primaryStyle string, fallbacks ...FontFallback) error {
	primary := fontKey{family: primaryFamily, style: getConvertedStyle(primaryStyle) &^ Underline}
	if _, ok := gp.subsetFontByKey(primary); !ok {
		return ErrFontNotFound
	}
	if len(fallbacks) == 0 {
		delete(gp.fontFallbacks, primary)
		return nil
	}
	keys := make([]fontKey, len(fallbacks))
	for i, fallback := range fallbacks {
		key := fontKey{family: fallback.Family, style: getConvertedStyle(fallback.Style) &^ Underline}
		if _, ok := gp.subsetFontByKey(key); !ok {
			return ErrFontNotFound
		}
		keys[i] = key
	}
	if gp.fontFallbacks == nil {
		gp.fontFallbacks = make(map[fontKey][]fontKey)
	}
	gp.fontFallbacks[primary] = keys
	return nil
}

func (gp *GoPdf) cloneFontFallbacks() map[fontKey][]fontKey {
	if len(gp.fontFallbacks) == 0 {
		return nil
	}
	cl := make(map[fontKey][]fontKey, len(gp.fontFallbacks))
	for key, fallbacks := range gp.fontFallbacks {
		cl[key] = append([]fontKey(nil), fallbacks...)
	}
	return cl
}

func (gp *GoPdf) subsetFontByKey(key fontKey) (*SubsetFontObj, bool) {
	for _, obj := range gp.pdfObjs {
		if obj.getType() != subsetFont {
			continue
		}
		font, ok := obj.(*SubsetFontObj)
		if !ok {
			continue
		}
		if font.GetFamily() == key.family && font.GetTtfFontOption().Style == key.style {
			return font, true
		}
	}
	return nil, false
}

func (gp *GoPdf) resolveTextLayout(text string) (textLayout, error) {
	current := gp.curr.FontISubset
	currentKey := fontKey{family: current.GetFamily(), style: current.GetTtfFontOption().Style &^ Underline}
	fallbacks, hasFallbacks := gp.fontFallbacks[currentKey]
	if !hasFallbacks || len(fallbacks) == 0 {
		text, err := current.AddChars(text)
		if err != nil {
			return textLayout{}, err
		}
		return textLayout{runs: []textRun{{text: text, fontSubset: current, fontCountIndex: current.CountOfFont + 1}}}, nil
	}

	layout := textLayout{}
	for _, r := range text {
		font := current
		glyphIndex, err := font.AvailableGlyphIndex(r)
		if err == ErrGlyphNotFound {
			font = nil
			for _, fallbackKey := range fallbacks {
				fallbackFont, ok := gp.subsetFontByKey(fallbackKey)
				if !ok {
					continue
				}
				glyphIndex, err = fallbackFont.AvailableGlyphIndex(r)
				if err == nil {
					font = fallbackFont
					break
				}
				if err != ErrGlyphNotFound {
					return textLayout{}, err
				}
			}
			if font == nil {
				replaced, err := current.AddChars(string(r))
				if err != nil {
					return textLayout{}, err
				}
				gp.appendTextRun(&layout, current, replaced)
				continue
			}
		} else if err != nil {
			return textLayout{}, err
		}
		font.AddRune(r, glyphIndex)
		gp.appendTextRun(&layout, font, string(r))
	}
	return layout, nil
}

func (gp *GoPdf) appendTextRun(layout *textLayout, font *SubsetFontObj, text string) {
	if text == "" {
		return
	}
	fontCountIndex := font.CountOfFont + 1
	lastIndex := len(layout.runs) - 1
	if lastIndex >= 0 && layout.runs[lastIndex].fontSubset == font {
		layout.runs[lastIndex].text += text
		return
	}
	layout.runs = append(layout.runs, textRun{text: text, fontSubset: font, fontCountIndex: fontCountIndex})
}
