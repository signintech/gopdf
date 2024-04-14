package gopdf

type PlaceHolderTextOption struct {
	//Left 8 , Right 2 ,Center  16
	Align int
}

type placeHolderTextInfo struct {
	indexOfContent   int
	indexInContent   int
	fontISubset      *SubsetFontObj
	placeHolderWidth float64
	fontSize         float64
	charSpacing      float64
}
