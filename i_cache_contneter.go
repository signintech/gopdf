package gopdf

import "bytes"

type iCacheContent interface {
	toStream() (*bytes.Buffer, error)
}
