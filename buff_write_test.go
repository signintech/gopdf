package gopdf

import (
	"bytes"
	"testing"
)

func TestWriteUInt32(t *testing.T) {
	var buff bytes.Buffer
	err := WriteUInt32(&buff, 65536)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	b := buff.Bytes()
	if b[0] != 0 || b[1] != 1 || b[2] != 0 || b[3] != 0 {
		t.Errorf("WriteUInt64 fail")
	}
}
