package gopdf

import "testing"

func TestSetProtection(t *testing.T) {

	var pp PDFProtection
	pp.setProtection(PermissionsPrint|PermissionsCopy|PermissionsModify, []byte("5555"), []byte("1234"))
	var realOValue = []byte{
		0xbb, 0xb8, 0x04, 0x6d, 0x96, 0xa9, 0x9a, 0x23, 0x46, 0xa9, 0x41, 0x21, 0x06, 0x8c, 0xad, 0x4f, 0x83, 0x5e, 0x5d, 0x0e, 0xcb, 0xb6, 0x20, 0xa8, 0xb7, 0xa3, 0x16, 0x13, 0x3c, 0x8f, 0x02, 0x91,
	}

	if !isSliceEq(pp.oValue, realOValue) {
		t.Errorf("wrong oValue")
		return
	}

}

func isSliceEq(a, b []byte) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
