package gopdf

import "bytes"

type iCacheContent interface {
	toStream(protection *PDFProtection) (*bytes.Buffer, error)
}
