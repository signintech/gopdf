package gopdf

import "io"

type listCacheContent struct {
	caches []iCacheContent
}

func (l *listCacheContent) last() iCacheContent {
	max := len(l.caches)
	if max > 0 {
		return l.caches[max-1]
	}
	return nil
}

func (l *listCacheContent) append(cache iCacheContent) {
	l.caches = append(l.caches, cache)
}

func (l *listCacheContent) appendContentText(cache cacheContentText, text string) (float64, float64, error) {

	x := cache.x
	y := cache.y

	mustMakeNewCache := true
	var cacheFont *cacheContentText
	var ok bool
	last := l.last()
	if cacheFont, ok = last.(*cacheContentText); ok {
		if cacheFont != nil {
			if cacheFont.isSame(cache) {
				mustMakeNewCache = false
			}
		}
	}

	if mustMakeNewCache { //make new cell
		l.caches = append(l.caches, &cache)
		cacheFont = &cache
	}

	//start add text
	cacheFont.text += text

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

func (l *listCacheContent) write(w io.Writer, protection *PDFProtection) error {
	for _, cache := range l.caches {
		if err := cache.write(w, protection); err != nil {
			return err
		}
	}
	return nil
}
