package gopdf

func parseStyle(style string) string {
	op := "S"
	if style == "F" {
		op = "f"
	} else if style == "FD" || style == "DF" {
		op = "B"
	}
	return op
}
