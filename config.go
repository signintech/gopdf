package gopdf

//Config static config
type Config struct {
	//pt , mm , cm , in
	Unit       string
	PageSize   Rect
	K          float64
	Protection PDFProtectionConfig
}

//PDFProtectionConfig config of pdf protection
type PDFProtectionConfig struct {
	UseProtection bool
	Permissions   int
	UserPass      []byte
	OwnerPass     []byte
}
