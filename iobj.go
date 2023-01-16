package gopdf

import (
	"io"
)

// IObj inteface for all pdf object
type IObj interface {
	init(func() *GoPdf)
	getType() string
	write(w io.Writer, objID int) error
}
