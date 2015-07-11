package core

type TableDirectoryEntry struct {
	CheckSum uint64
	Offset   uint64
	Length   uint64
}

func (t TableDirectoryEntry) PaddedLength() int {
	l := int(t.Length)
	return (l + 3) & ^3
}
