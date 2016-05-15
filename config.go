package gopdf

//Config static config
type Config struct {
	//pt , mm , cm , in
	Unit     string
	PageSize Rect
	K        float64
}
