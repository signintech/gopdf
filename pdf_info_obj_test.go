package gopdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode/utf16"
)

func utf16BEWithBOM(s string) []byte {
	encoded := utf16.Encode([]rune(s))
	out := make([]byte, 2+len(encoded)*2)
	out[0] = 0xFE
	out[1] = 0xFF
	for i, v := range encoded {
		binary.BigEndian.PutUint16(out[2+i*2:4+i*2], v)
	}
	return out
}

// TestSetInfoNonBMP checks that SetInfo writes valid UTF-16BE metadata
// for emoji and CJK characters, both with and without PDF protection.
func TestSetInfoNonBMP(t *testing.T) {
	const (
		title  = "Report 😀"
		author = "𠀀"
	)

	for _, protected := range []bool{false, true} {
		name := "plain"
		if protected {
			name = "protected"
		}
		t.Run(name, func(t *testing.T) {
			var pdf GoPdf
			pdf.Start(Config{
				PageSize: *PageSizeA4,
				Protection: PDFProtectionConfig{
					UseProtection: protected,
					UserPass:      []byte("user"),
					OwnerPass:     []byte("owner"),
					Permissions:   PermissionsPrint,
				},
			})
			pdf.SetInfo(PdfInfo{
				Title:  title,
				Author: author,
			})
			pdf.AddPage()

			var buf bytes.Buffer
			if _, err := pdf.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo failed: %v", err)
			}

			body := buf.String()
			objHeader := fmt.Sprintf("%d 0 obj\n", pdf.infoObjID)
			idx := strings.Index(body, objHeader)
			if idx < 0 {
				t.Fatalf("missing Info object %q", objHeader)
			}
			rest := body[idx:]
			endIdx := strings.Index(rest, "endobj\n")
			if endIdx < 0 {
				t.Fatalf("Info object missing endobj")
			}
			infoObjStr := rest[:endIdx+len("endobj\n")]
			if !strings.Contains(infoObjStr, "<<\n") || !strings.Contains(infoObjStr, ">>\nendobj\n") {
				t.Fatalf("Info object dictionary is malformed/unterminated:\n%s", infoObjStr)
			}

			if !protected {
				wantTitle := "/Title <FEFF005200650070006F007200740020D83DDE00>\n"
				wantAuthor := "/Author <FEFFD840DC00>\n"
				if !strings.Contains(infoObjStr, wantTitle) {
					t.Fatalf("Info object missing %q in:\n%s", wantTitle, infoObjStr)
				}
				if !strings.Contains(infoObjStr, wantAuthor) {
					t.Fatalf("Info object missing %q in:\n%s", wantAuthor, infoObjStr)
				}
				return
			}

			key := pdf.protection().objectkey(pdf.infoObjID)
			encTitle, err := rc4Cip(key, utf16BEWithBOM(title))
			if err != nil {
				t.Fatal(err)
			}
			encAuthor, err := rc4Cip(key, utf16BEWithBOM(author))
			if err != nil {
				t.Fatal(err)
			}
			wantTitle := fmt.Sprintf("/Title <%X>\n", encTitle)
			wantAuthor := fmt.Sprintf("/Author <%X>\n", encAuthor)
			if !strings.Contains(infoObjStr, wantTitle) {
				t.Fatalf("protected Info object missing encrypted title %q in:\n%s", wantTitle, infoObjStr)
			}
			if !strings.Contains(infoObjStr, wantAuthor) {
				t.Fatalf("protected Info object missing encrypted author %q in:\n%s", wantAuthor, infoObjStr)
			}
		})
	}
}

type failingIObj struct {
	err error
}

func (f failingIObj) init(func() *GoPdf) {}

func (f failingIObj) getType() string { return "FailingMock" }

func (f failingIObj) write(io.Writer, int) error { return f.err }

// TestCompilePdfPropagatesObjWriteError checks that WriteTo returns an error
// when writing a PDF object fails, instead of hiding it.
func TestCompilePdfPropagatesObjWriteError(t *testing.T) {
	var pdf GoPdf
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	sentinel := errors.New("obj write failed")
	pdf.addObj(failingIObj{err: sentinel})

	var buf bytes.Buffer
	_, err := pdf.WriteTo(&buf)
	if !errors.Is(err, sentinel) {
		t.Fatalf("WriteTo error = %v, want %v", err, sentinel)
	}
}

// TestCompilePdfPropagatesWriterError checks that WriteTo returns an error
// when the output writer fails.
func TestCompilePdfPropagatesWriterError(t *testing.T) {
	var pdf GoPdf
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	_, err := pdf.WriteTo(failedPDFWriter{})
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("WriteTo error = %v, want %v", err, io.ErrClosedPipe)
	}
}
