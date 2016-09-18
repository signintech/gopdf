package gopdf

import "testing"

func TestSetProtection(t *testing.T) {
	var pp PDFProtection
	pp.setProtection(PermissionsPrint|PermissionsCopy|PermissionsModify, []byte("5555"), []byte("1234"))
}
