package gopdf

// MapOfCharacterToGlyphIndex map of CharacterToGlyphIndex
type MapOfCharacterToGlyphIndex struct {
	keyIndexs map[rune]int //for search index in keys
	Keys      []rune
	Vals      []uint
}

// NewMapOfCharacterToGlyphIndex new CharacterToGlyphIndex
func NewMapOfCharacterToGlyphIndex() *MapOfCharacterToGlyphIndex {
	var m MapOfCharacterToGlyphIndex
	m.keyIndexs = make(map[rune]int)
	return &m
}

// KeyExists key is exists?
func (m *MapOfCharacterToGlyphIndex) KeyExists(k rune) bool {
	/*for _, key := range m.Keys {
		if k == key {
			return true
		}
	}*/
	if _, ok := m.keyIndexs[k]; ok {
		return true
	}
	return false
}

// Set set key and value to map
func (m *MapOfCharacterToGlyphIndex) Set(k rune, v uint) {
	m.keyIndexs[k] = len(m.Keys)
	m.Keys = append(m.Keys, k)
	m.Vals = append(m.Vals, v)
}

// Index get index by key
func (m *MapOfCharacterToGlyphIndex) Index(k rune) (int, bool) {
	/*for i, key := range m.Keys {
		if k == key {
			return i, true
		}
	}*/
	if index, ok := m.keyIndexs[k]; ok {
		return index, true
	}
	return -1, false
}

// Val get value by Key
func (m *MapOfCharacterToGlyphIndex) Val(k rune) (uint, bool) {
	i, ok := m.Index(k)
	if !ok {
		return 0, false
	}
	return m.Vals[i], true
}

// AllKeys get keys
func (m *MapOfCharacterToGlyphIndex) AllKeys() []rune {
	return m.Keys
}

// AllVals get all values
func (m *MapOfCharacterToGlyphIndex) AllVals() []uint {
	return m.Vals
}
