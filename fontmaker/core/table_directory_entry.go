package core

type TableDirectoryEntry struct {
	CheckSum uint
	Offset   uint
	Length   uint
}

func (t TableDirectoryEntry) PaddedLength() int {
	l := int(t.Length)
	return (l + 3) & ^3
}
