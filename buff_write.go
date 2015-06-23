package gopdf

import "io"

func WriteUInt32(w io.Writer, v uint32) error {
	a := byte(v >> 24)
	b := byte(v >> 16)
	c := byte(v >> 8)
	d := byte(v)
	_, err := w.Write([]byte{a, b, c, d})
	if err != nil {
		return err
	}
	return nil
}

func WriteUInt16(w io.Writer, v uint16) error {

	a := byte(v >> 8)
	b := byte(v)
	_, err := w.Write([]byte{a, b})
	if err != nil {
		return err
	}
	return nil
}
