package gopdf

import (
	"fmt"
	"io"
)

type cacheContentSaveGraphicsState struct{}

func (c *cacheContentSaveGraphicsState) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprint(w, "q\n")
	return nil
}

type cacheContentRestoreGraphicsState struct{}

func (c *cacheContentRestoreGraphicsState) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprint(w, "Q\n")
	return nil
}
