package gopdf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
)

// ErrFontNotFound occurs when the specified font family is not found in the container.
var ErrFontNotFound = errors.New("font not found")

// FontContainer manages a collection of fonts.
type FontContainer struct {
	fonts sync.Map
}

// AddTTFFontWithOption adds a font by the specified TTF file path with options.
func (fc *FontContainer) AddTTFFontWithOption(family string, ttfpath string, option TtfOption) error {
	if _, err := os.Stat(ttfpath); os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(ttfpath)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(data)
	return fc.AddTTFFontByReaderWithOption(family, rd, option)
}

// AddTTFFont adds a font by the specified TTF file path.
func (fc *FontContainer) AddTTFFont(family string, ttfpath string) error {
	return fc.AddTTFFontWithOption(family, ttfpath, defaultTtfFontOption())
}

// AddTTFFontByReader adds font by reader.
func (fc *FontContainer) AddTTFFontByReader(family string, rd io.Reader) error {
	return fc.AddTTFFontByReaderWithOption(family, rd, defaultTtfFontOption())
}

// AddTTFFontByReaderWithOption adds font by reader with option.
func (fc *FontContainer) AddTTFFontByReaderWithOption(family string, rd io.Reader, option TtfOption) error {
	var subsetFont SubsetFontObj
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFByReader(rd)
	if err != nil {
		return err
	}
	fc.fonts.Store(family, subsetFont)
	return nil
}

// AddTTFFontData adds font by data.
func (fc *FontContainer) AddTTFFontData(family string, fontData []byte) error {
	return fc.AddTTFFontDataWithOption(family, fontData, defaultTtfFontOption())
}

// AddTTFFontDataWithOption adds font by data with option.
func (fc *FontContainer) AddTTFFontDataWithOption(family string, fontData []byte, option TtfOption) error {
	var subsetFont SubsetFontObj
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFData(fontData)
	if err != nil {
		return err
	}
	fc.fonts.Store(family, subsetFont)
	return nil
}

// AddTTFFontFromFontContainer adds a font from a FontContainer
func (gp *GoPdf) AddTTFFontFromFontContainer(family string, container *FontContainer) error {
	untypedSubsetFontObj, ok := container.fonts.Load(family)
	if !ok {
		return ErrFontNotFound
	}
	subsetFont := untypedSubsetFontObj.(SubsetFontObj)
	subsetFont.init(func() *GoPdf {
		return gp
	})
	return gp.setSubsetFontObject(&subsetFont, family, subsetFont.ttfFontOption)
}
