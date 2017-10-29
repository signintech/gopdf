package gopdf

import (
	"io"
)

type iCacheContent interface {
	write(w io.Writer, protection *PDFProtection) error
}
