package gopdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type failedPDFWriter struct{}

func (failedPDFWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestRepeatedPDFWrites(t *testing.T) {
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
					UserPass:      []byte("reader"),
					OwnerPass:     []byte("owner"),
				},
			})
			pdf.AddPage()
			if err := pdf.AddTTFFont("fixture", "test/res/LiberationSerif-Regular.ttf"); err != nil {
				t.Fatal(err)
			}
			if err := pdf.SetFont("fixture", "", 12); err != nil {
				t.Fatal(err)
			}
			if err := pdf.Cell(nil, "First page"); err != nil {
				t.Fatal(err)
			}
			pdf.Line(10, 10, 50, 50)
			pdf.SetInfo(PdfInfo{
				Title:        "Test Title",
				Author:       "Test Author",
				Subject:      "Test Subject",
				Creator:      "Test Creator",
				Producer:     "Test Producer",
				CreationDate: time.Date(2026, 10, 31, 12, 0, 0, 0, time.FixedZone("IST", 5*3600+30*60)),
			})
			first, err := pdf.GetBytesPdfReturnErr()
			if err != nil {
				t.Fatal(err)
			}
			first = append([]byte(nil), first...)

			expectedInfoRef := fmt.Sprintf("/Info %d 0 R\n", pdf.infoObjID)
			if !strings.Contains(string(first), expectedInfoRef) {
				t.Fatalf("trailer missing indirect Info reference %q", expectedInfoRef)
			}
			plainInfoObj := fmt.Sprintf("%d 0 obj\n<<\n/Author <FEFF005400650073007400200041007500740068006F0072>\n/Title <FEFF00540065007300740020005400690074006C0065>\n/Subject <FEFF00540065007300740020005300750062006A006500630074>\n/Creator <FEFF0054006500730074002000430072006500610074006F0072>\n/Producer <FEFF0054006500730074002000500072006F00640075006300650072>\n/CreationDate(D:20261031120000+05'30')\n>>\nendobj\n", pdf.infoObjID)
			if !protected {
				if !strings.Contains(string(first), plainInfoObj) {
					t.Fatalf("PDF body missing expected Info object:\n%s", plainInfoObj)
				}
			} else {
				if strings.Contains(string(first), "20261031120000") {
					t.Fatalf("protected PDF leaked unencrypted CreationDate in Info object")
				}
				encDate, err := rc4Cip(pdf.protection().objectkey(pdf.infoObjID), []byte("D:20261031120000+05'30'"))
				if err != nil {
					t.Fatal(err)
				}
				expectedEncDate := fmt.Sprintf("/CreationDate <%X>\n", encDate)
				if !strings.Contains(string(first), expectedEncDate) {
					t.Fatalf("protected PDF missing encrypted CreationDate %q", expectedEncDate)
				}
			}

			var output bytes.Buffer
			if _, err := pdf.WriteTo(&output); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(first, output.Bytes()) {
				t.Error("WriteTo changed an unchanged document")
			}

			path := filepath.Join(t.TempDir(), "document.pdf")
			if err := pdf.WritePdf(path); err != nil {
				t.Fatal(err)
			}
			written, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(first, written) {
				t.Error("WritePdf changed an unchanged document")
			}

			// WriteTo currently discards writer errors. Even after a failed
			// destination write, preparing the document again must be safe.
			if n, _ := pdf.WriteTo(failedPDFWriter{}); n != 0 {
				t.Fatalf("failing writer accepted %d bytes", n)
			}
			retried, err := pdf.GetBytesPdfReturnErr()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(first, retried) {
				t.Error("retry after a failed write changed the document")
			}

			pdf.AddPage()
			pdf.Line(20, 20, 60, 60)
			updated, err := pdf.GetBytesPdfReturnErr()
			if err != nil {
				t.Fatal(err)
			}
			updated = append([]byte(nil), updated...)
			if !bytes.Contains(updated, []byte("/Count 2\n")) {
				t.Error("adding a page did not produce a two-page document")
			}
			repeated, err := pdf.GetBytesPdfReturnErr()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(updated, repeated) {
				t.Error("repeated write changed the updated document")
			}
		})
	}
}
