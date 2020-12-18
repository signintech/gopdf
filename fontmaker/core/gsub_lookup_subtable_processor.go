package core

type gsubLookupSubtableProcessor interface {
	process(glyphindexs []uint) ([]uint, error)
}
