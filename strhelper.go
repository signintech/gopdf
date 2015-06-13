package gopdf

func StrHelperGetStringWidth(str string, fontSize int, ifont IFont) float64 {
	w := 0
	bs := []byte(str)
	i := 0
	max := len(bs)
	for i < max {
		w += ifont.GetCw()[bs[i]]
		i++
	}
	return float64(w) * (float64(fontSize) / 1000.0)
}

func CreateEmbeddedFontSubsetName(name string) string {
	//TODO ทำด้วย  :-)
	return name
}
