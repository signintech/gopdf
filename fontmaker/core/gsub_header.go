package core

type GSUBHeader struct {
	majorVersion      uint
	minorVersion      int
	gsubOffset        int64
	scriptListOffset  int64
	featureListOffset int64
	lookupListOffset  int64
}
