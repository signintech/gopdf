package gopdf

import "bytes"

type cacheContent struct {
	rectangle      *Rect
	textColor      Rgb
	grayFill       float64
	fontCountIndex int //Curr.Font_FontCount+1
	fontSize       int
	fontStyle      string
	x, y           float64
	text           bytes.Buffer
	fontSubset     *SubsetFontObj
}

func (c *cacheContent) isSame(cache cacheContent) bool {
	if c.rectangle != nil {
		//if rectangle != nil we assumes this is not same content
		return false
	}
	if c.textColor.equal(cache.textColor) &&
		c.grayFill == cache.grayFill &&
		c.fontCountIndex == cache.fontCountIndex &&
		c.fontSize == cache.fontSize &&
		c.fontStyle == cache.fontStyle &&
		//c.x == cache.x &&
		c.y == cache.y {
		return true
	}

	return false
}

type listCacheContent struct {
	caches []cacheContent
}

func (l *listCacheContent) last() *cacheContent {
	max := len(l.caches)
	if max > 0 {
		return &l.caches[max-1]
	}
	return nil
}

func (l *listCacheContent) appendTextToCache(cache cacheContent, text string) error {
	//TODO ควรจัดการ Curr.X กะ Curr.Y ที่นี้
	mustMakeNewCache := true
	cacheFont := l.last()
	if cacheFont != nil {
		if cacheFont.isSame(cache) {
			mustMakeNewCache = false
		}
	}

	if mustMakeNewCache {
		l.caches = append(l.caches, cache)
		cacheFont = l.last()
	}
	_, err := cacheFont.text.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

func (l *listCacheContent) toStream() (*bytes.Buffer, error) {

	return nil, nil
}

func (l *listCacheContent) debug() string {
	var buff bytes.Buffer
	for _, cache := range l.caches {
		buff.WriteString(cache.text.String())
		buff.WriteString("\n")
	}
	return buff.String()
}
