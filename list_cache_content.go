package gopdf

import "bytes"

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

func (l *listCacheContent) appendTextToCache(cache cacheContent, text string) (float64, float64, error) {

	x := cache.x
	y := cache.y

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
		return x, y, err
	}

	//re-create contnet
	textWidthPdfUnit, textHeightPdfUnit, err := cacheFont.createContent()
	if err != nil {
		return x, y, err
	}

	if cacheFont.cellOpt.Float == 0 || cacheFont.cellOpt.Float&Right == Right || cacheFont.contentType == ContentTypeText {
		x += textWidthPdfUnit
	}
	if cacheFont.cellOpt.Float&Bottom == Bottom {
		y += textHeightPdfUnit
	}

	return x, y, nil
}

func (l *listCacheContent) toStream() (*bytes.Buffer, error) {
	var buff bytes.Buffer
	for _, cache := range l.caches {
		stream, err := cache.toStream()
		if err != nil {
			return nil, err
		}
		_, err = stream.WriteTo(&buff)
		if err != nil {
			return nil, err
		}
	}
	return &buff, nil
}

func (l *listCacheContent) debug() string {
	var buff bytes.Buffer
	for _, cache := range l.caches {
		buff.WriteString(cache.text.String())
		buff.WriteString("\n")
	}
	return buff.String()
}
