package gopdf

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

//ContentObj content object
type ContentObj struct { //impl IObj
	buffer    bytes.Buffer
	stream    bytes.Buffer
	listCache listCacheContent
	//text bytes.Buffer
	getRoot func() *GoPdf
}

func (c *ContentObj) protection() *PDFProtection {
	return c.getRoot().protection()
}

func (c *ContentObj) init(funcGetRoot func() *GoPdf) {
	c.getRoot = funcGetRoot
}

func (c *ContentObj) build(objID int) error {
	buff, err := c.listCache.toStream(c.protection())
	if err != nil {
		return err
	}
	buff.WriteTo(&c.stream)

	streamlen := c.stream.Len()
	c.buffer.WriteString("<<\n")
	c.buffer.WriteString("/Length " + strconv.Itoa(streamlen) + "\n")
	c.buffer.WriteString(">>\n")
	c.buffer.WriteString("stream\n")
	if c.protection() != nil {
		tmp, err := rc4Cip(c.protection().objectkey(objID), c.stream.Bytes())
		if err != nil {
			return err
		}
		c.buffer.Write(tmp)
		c.buffer.WriteString("\n")
	} else {
		c.buffer.Write(c.stream.Bytes())
	}
	c.buffer.WriteString("endstream\n")
	return nil
}

func (c *ContentObj) getType() string {
	return "Content"
}

func (c *ContentObj) getObjBuff() *bytes.Buffer {
	return &(c.buffer)
}

//AppendStreamText append text
func (c *ContentObj) AppendStreamText(text string) error {

	//support only CURRENT_FONT_TYPE_SUBSET
	textColor := c.getRoot().curr.textColor()
	grayFill := c.getRoot().curr.grayFill
	fontCountIndex := c.getRoot().curr.Font_FontCount + 1
	fontSize := c.getRoot().curr.Font_Size
	fontStyle := c.getRoot().curr.Font_Style
	x := c.getRoot().curr.X
	y := c.getRoot().curr.Y
	setXCount := c.getRoot().curr.setXCount
	fontSubset := c.getRoot().curr.Font_ISubset

	cache := cacheContentText{
		fontSubset:     fontSubset,
		rectangle:      nil,
		textColor:      textColor,
		grayFill:       grayFill,
		fontCountIndex: fontCountIndex,
		fontSize:       fontSize,
		fontStyle:      fontStyle,
		setXCount:      setXCount,
		x:              x,
		y:              y,
		pageheight:     c.getRoot().config.PageSize.H,
		contentType:    ContentTypeText,
		lineWidth:      c.getRoot().curr.lineWidth,
	}

	var err error
	c.getRoot().curr.X, c.getRoot().curr.Y, err = c.listCache.appendContentText(cache, text)
	if err != nil {
		return err
	}

	return nil
}

//AppendStreamSubsetFont add stream of text
func (c *ContentObj) AppendStreamSubsetFont(rectangle *Rect, text string, cellOpt CellOption) error {

	textColor := c.getRoot().curr.textColor()
	grayFill := c.getRoot().curr.grayFill
	fontCountIndex := c.getRoot().curr.Font_FontCount + 1
	fontSize := c.getRoot().curr.Font_Size
	fontStyle := c.getRoot().curr.Font_Style
	x := c.getRoot().curr.X
	y := c.getRoot().curr.Y
	setXCount := c.getRoot().curr.setXCount
	fontSubset := c.getRoot().curr.Font_ISubset

	cache := cacheContentText{
		fontSubset:     fontSubset,
		rectangle:      rectangle,
		textColor:      textColor,
		grayFill:       grayFill,
		fontCountIndex: fontCountIndex,
		fontSize:       fontSize,
		fontStyle:      fontStyle,
		setXCount:      setXCount,
		x:              x,
		y:              y,
		pageheight:     c.getRoot().config.PageSize.H,
		contentType:    ContentTypeCell,
		cellOpt:        cellOpt,
		lineWidth:      c.getRoot().curr.lineWidth,
	}
	var err error
	c.getRoot().curr.X, c.getRoot().curr.Y, err = c.listCache.appendContentText(cache, text)
	if err != nil {
		return err
	}
	return nil
}

//AppendStreamLine append line
func (c *ContentObj) AppendStreamLine(x1 float64, y1 float64, x2 float64, y2 float64) {
	//h := c.getRoot().config.PageSize.H
	//c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", x1, h-y1, x2, h-y2))
	var cache cacheContentLine
	cache.pageHeight = c.getRoot().config.PageSize.H
	cache.x1 = x1
	cache.y1 = y1
	cache.x2 = x2
	cache.y2 = y2
	c.listCache.append(&cache)
}

//AppendStreamRectangle : draw rectangle from lower-left corner (x, y) with specif width/height
func (c *ContentObj) AppendStreamRectangle(x float64, y float64, wdth float64, hght float64, style string) {
	var cache cacheContentRectangle
	cache.pageHeight = c.getRoot().config.PageSize.H
	cache.x = x
	cache.y = y
	cache.width = wdth
	cache.height = hght
	cache.style = style
	c.listCache.append(&cache)
}

//AppendStreamOval append oval
func (c *ContentObj) AppendStreamOval(x1 float64, y1 float64, x2 float64, y2 float64) {
	var cache cacheContentOval
	cache.pageHeight = c.getRoot().config.PageSize.H
	cache.x1 = x1
	cache.y1 = y1
	cache.x2 = x2
	cache.y2 = y2
	c.listCache.append(&cache)
}

//AppendStreamCurve draw curve
// - x0, y0: Start point
// - x1, y1: Control point 1
// - x2, y2: Control point 2
// - x3, y3: End point
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//		D or empty string: draw. This is the default value.
//		F: fill
//		DF or FD: draw and fill
func (c *ContentObj) AppendStreamCurve(x0 float64, y0 float64, x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64, style string) {
	var cache cacheContentCurve
	cache.pageHeight = c.getRoot().config.PageSize.H
	cache.x0 = x0
	cache.y0 = y0
	cache.x1 = x1
	cache.y1 = y1
	cache.x2 = x2
	cache.y2 = y2
	cache.x3 = x3
	cache.y3 = y3
	cache.style = strings.ToUpper(strings.TrimSpace(style))
	c.listCache.append(&cache)
}

//AppendStreamSetLineWidth : set line width
func (c *ContentObj) AppendStreamSetLineWidth(w float64) {
	var cache cacheContentLineWidth
	cache.width = w
	c.listCache.append(&cache)
}

//AppendStreamSetLineType : Set linetype [solid, dashed, dotted]
func (c *ContentObj) AppendStreamSetLineType(t string) {
	var cache cacheContentLineType
	cache.lineType = t
	c.listCache.append(&cache)

}

//AppendStreamSetGrayFill  set the grayscale fills
func (c *ContentObj) AppendStreamSetGrayFill(w float64) {
	w = fixRange10(w)
	var cache cacheContentGray
	cache.grayType = grayTypeFill
	cache.scale = w
	c.listCache.append(&cache)
}

//AppendStreamSetGrayStroke  set the grayscale stroke
func (c *ContentObj) AppendStreamSetGrayStroke(w float64) {
	w = fixRange10(w)
	var cache cacheContentGray
	cache.grayType = grayTypeStroke
	cache.scale = w
	c.listCache.append(&cache)
}

//AppendStreamSetColorStroke  set the color stroke
func (c *ContentObj) AppendStreamSetColorStroke(r uint8, g uint8, b uint8) {
	var cache cacheContentColor
	cache.colorType = colorTypeStroke
	cache.r = r
	cache.g = g
	cache.b = b
	c.listCache.append(&cache)
}

//AppendStreamSetColorFill  set the color fill
func (c *ContentObj) AppendStreamSetColorFill(r uint8, g uint8, b uint8) {
	var cache cacheContentColor
	cache.colorType = colorTypeFill
	cache.r = r
	cache.g = g
	cache.b = b
	c.listCache.append(&cache)
}

//AppendStreamImage append image
func (c *ContentObj) AppendStreamImage(index int, x float64, y float64, rect *Rect) {
	//fmt.Printf("index = %d",index)
	h := c.getRoot().config.PageSize.H
	c.stream.WriteString(fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", rect.W, rect.H, x, h-(y+rect.H), index+1))
}

//ContentObj_CalTextHeight calculate height of text
func ContentObj_CalTextHeight(fontsize int) float64 {
	return (float64(fontsize) * 0.7)
}

// When setting colour and grayscales the value has to be between 0.00 and 1.00
// This function takes a float64 and returns 0.0 if it is less than 0.0 and 1.0 if it
// is more than 1.0
func fixRange10(val float64) float64 {
	if val < 0.0 {
		return 0.0
	}
	if val > 1.0 {
		return 1.0
	}
	return val
}

func convertTTFUnit2PDFUnit(n int, upem int) int {
	var ret int
	if n < 0 {
		rest1 := n % upem
		storrest := 1000 * rest1
		//ledd2 := (storrest != 0 ? rest1 / storrest : 0);
		ledd2 := 0
		if storrest != 0 {
			ledd2 = rest1 / storrest
		} else {
			ledd2 = 0
		}
		ret = -((-1000*n)/upem - int(ledd2))
	} else {
		ret = (n/upem)*1000 + ((n%upem)*1000)/upem
	}
	return ret
}
