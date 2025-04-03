package gopdf

import (
	"fmt"
	"io"
)

type cacheColorSpace struct {
	countOfSpaceColor int
}

func (c *cacheColorSpace) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "/CS%d CS 1.0000 SCN\n", c.countOfSpaceColor+1)
	return nil
}
