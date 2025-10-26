package shaping

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	gotextfont "github.com/go-text/typesetting/font"
	hb "github.com/go-text/typesetting/harfbuzz"
)

func TestMyanmarMetrics_Debug(t *testing.T) {
	fontPath := "./NotoSansMyanmar-Regular.ttf"
	data, err := os.ReadFile(fontPath)
	if err != nil {
		t.Fatalf("read font: %v", err)
	}
	face, err := gotextfont.ParseTTF(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("parse font: %v", err)
	}
	f := hb.NewFont(face)
	buf := hb.NewBuffer()
	// Test the exact problem word
	text := []rune("နေ့")
	buf.AddRunes(text, 0, -1)
	buf.GuessSegmentProperties()
	buf.Shape(f, nil)

	info := buf.Info
	pos := buf.Pos
	for i := range info {
		fmt.Printf("idx=%d gid=%d adv=%d xoff=%d yoff=%d\n", i, info[i].Glyph, pos[i].XAdvance, pos[i].XOffset, pos[i].YOffset)
	}
}
