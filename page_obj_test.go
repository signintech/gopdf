package gopdf

import (
	"bytes"
	"strings"
	"testing"
)

// TestPageObjWriteOmitsContentsWhenEmpty verifies that PageObj.write does not
// emit a "/Contents" entry when PageObj.Contents is empty.
//
// Previously the write method always wrote "  /Contents %s\n", which produced
// a malformed "  /Contents \n" line in the PDF when no content stream was
// attached (e.g. when importing pages from an existing source).
func TestPageObjWriteOmitsContentsWhenEmpty(t *testing.T) {
	p := &PageObj{
		ResourcesRelate: "5 0 R",
	}

	var buf bytes.Buffer
	if err := p.write(&buf, 1); err != nil {
		t.Fatalf("PageObj.write error = %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "/Contents") {
		t.Fatalf("PageObj.write emitted /Contents for empty Contents; output:\n%s", out)
	}
	if !strings.HasPrefix(out, "<<\n") || !strings.HasSuffix(out, ">>\n") {
		t.Fatalf("PageObj.write produced malformed dictionary:\n%s", out)
	}
}

// TestPageObjWriteIncludesContentsWhenSet verifies that PageObj.write emits the
// "/Contents" entry with the correct indirect reference when PageObj.Contents
// is populated.
func TestPageObjWriteIncludesContentsWhenSet(t *testing.T) {
	p := &PageObj{
		Contents:        "8 0 R",
		ResourcesRelate: "5 0 R",
	}

	var buf bytes.Buffer
	if err := p.write(&buf, 1); err != nil {
		t.Fatalf("PageObj.write error = %v", err)
	}

	want := "  /Contents 8 0 R\n"
	if !strings.Contains(buf.String(), want) {
		t.Fatalf("PageObj.write missing %q in:\n%s", want, buf.String())
	}
}
