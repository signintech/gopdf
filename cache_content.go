package gopdf

import (
	"bytes"
	"fmt"
	"strconv"
)

type cacheContent struct {
	rectangle      *Rect
	textColor      Rgb
	grayFill       float64
	fontCountIndex int //Curr.Font_FontCount+1
	fontSize       int
	fontStyle      string
	x, y           float64
	text           bytes.Buffer
	textWidth      float64
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

func (c *cacheContent) toStream() (*bytes.Buffer, error) {
	var stream, textbuff bytes.Buffer
	text := c.text.String()
	for _, r := range text {
		glyphindex, err := c.fontSubset.CharIndex(r)
		if err != nil {
			return nil, err
		}
		textbuff.WriteString(fmt.Sprintf("%04X", glyphindex))
	}

	pageHeight := 841.89 //TODO fix this //c.getRoot().config.PageSize.H
	r := c.textColor.r
	g := c.textColor.g
	b := c.textColor.b
	x := fmt.Sprintf("%0.2f", c.x)
	y := fmt.Sprintf("%0.2f", pageHeight-c.y-(float64(c.fontSize)*0.7))

	stream.WriteString("BT\n")
	stream.WriteString(x + " " + y + " TD\n")
	stream.WriteString("/F" + strconv.Itoa(c.fontCountIndex) + " " + strconv.Itoa(c.fontSize) + " Tf\n")
	if r+g+b != 0 {
		rFloat := float64(r) * 0.00392156862745
		gFloat := float64(g) * 0.00392156862745
		bFloat := float64(b) * 0.00392156862745
		rgb := fmt.Sprintf("%0.2f %0.2f %0.2f rg\n", rFloat, gFloat, bFloat)
		stream.WriteString(rgb)
	} else {
		//c.AppendStreamSetGrayFill(grayFill)
	}

	stream.WriteString("[<" + textbuff.String() + ">] TJ\n")
	stream.WriteString("ET\n")
	return &stream, nil
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
