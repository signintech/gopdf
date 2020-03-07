package gopdf

import (
	"compress/zlib"
	"fmt"
	"io"
	"strings"
)

//ContentObj content object
type ContentObj struct { //impl IObj
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

func (c *ContentObj) write(w io.Writer, objID int) error {
	buff := GetBuffer()
	defer PutBuffer(buff)

	isFlate := (c.getRoot().compressLevel != zlib.NoCompression)
	if isFlate {
		ww, err := zlib.NewWriterLevel(buff, c.getRoot().compressLevel)
		if err != nil {
			// should never happen...
			return err
		}
		if err := c.listCache.write(ww, c.protection()); err != nil {
			return err
		}
		ww.Close()
	} else {
		if err := c.listCache.write(buff, c.protection()); err != nil {
			return err
		}
	}

	streamlen := buff.Len()

	io.WriteString(w, "<<\n")
	if isFlate {
		io.WriteString(w, "/Filter/FlateDecode")
	}
	fmt.Fprintf(w, "/Length %d\n", streamlen)
	io.WriteString(w, ">>\n")
	io.WriteString(w, "stream\n")
	if c.protection() != nil {
		tmp, err := rc4Cip(c.protection().objectkey(objID), buff.Bytes())
		if err != nil {
			return err
		}
		w.Write(tmp)
		io.WriteString(w, "\n")
	} else {
		buff.WriteTo(w)
		if isFlate {
			io.WriteString(w, "\n")
		}
	}
	io.WriteString(w, "endstream\n")

	return nil
}

func (c *ContentObj) getType() string {
	return "Content"
}

//AppendStreamText append text
func (c *ContentObj) AppendStreamText(text string) error {

	//support only CURRENT_FONT_TYPE_SUBSET
	textColor := c.getRoot().curr.textColor()
	grayFill := c.getRoot().curr.grayFill
	fontCountIndex := c.getRoot().curr.FontFontCount + 1
	fontSize := c.getRoot().curr.FontSize
	fontStyle := c.getRoot().curr.FontStyle
	x := c.getRoot().curr.X
	y := c.getRoot().curr.Y
	setXCount := c.getRoot().curr.setXCount
	fontSubset := c.getRoot().curr.FontISubset
	transparency := c.getRoot().curr.transparency

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
		pageheight:     c.getRoot().curr.pageSize.H,
		contentType:    ContentTypeText,
		lineWidth:      c.getRoot().curr.lineWidth,
		txtColorMode:   c.getRoot().curr.txtColorMode,
		transparency:   transparency,
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
	fontCountIndex := c.getRoot().curr.FontFontCount + 1
	fontSize := c.getRoot().curr.FontSize
	fontStyle := c.getRoot().curr.FontStyle
	x := c.getRoot().curr.X
	y := c.getRoot().curr.Y
	setXCount := c.getRoot().curr.setXCount
	fontSubset := c.getRoot().curr.FontISubset
	transparency := c.getRoot().curr.transparency

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
		pageheight:     c.getRoot().curr.pageSize.H,
		contentType:    ContentTypeCell,
		cellOpt:        cellOpt,
		lineWidth:      c.getRoot().curr.lineWidth,
		txtColorMode:   c.getRoot().curr.txtColorMode,
		transparency:   transparency,
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
	cache.pageHeight = c.getRoot().curr.pageSize.H
	cache.x1 = x1
	cache.y1 = y1
	cache.x2 = x2
	cache.y2 = y2
	c.listCache.append(&cache)
}

//AppendStreamImportedTemplate append imported template
func (c *ContentObj) AppendStreamImportedTemplate(tplName string, scaleX float64, scaleY float64, tX float64, tY float64) {
	var cache cacheContentImportedTemplate
	cache.pageHeight = c.getRoot().curr.pageSize.H
	cache.tplName = tplName
	cache.scaleX = scaleX
	cache.scaleY = scaleY
	cache.tX = tX
	cache.tY = tY
	c.listCache.append(&cache)
}

//AppendStreamRectangle : draw rectangle from lower-left corner (x, y) with specif width/height
func (c *ContentObj) AppendStreamRectangle(x float64, y float64, wdth float64, hght float64, style string) {
	var cache cacheContentRectangle
	cache.pageHeight = c.getRoot().curr.pageSize.H
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
	cache.pageHeight = c.getRoot().curr.pageSize.H
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
	cache.pageHeight = c.getRoot().curr.pageSize.H
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
	h := c.getRoot().curr.pageSize.H
	var cache cacheContentImage
	cache.h = h
	cache.x = x
	cache.y = y
	cache.rect = *rect
	cache.index = index
	c.listCache.append(&cache)
	//c.stream.WriteString(fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", rect.W, rect.H, x, h-(y+rect.H), index+1))
}

func (c *ContentObj) appendRotate(angle, x, y float64) {
	var cache cacheContentRotate
	cache.isReset = false
	cache.pageHeight = c.getRoot().curr.pageSize.H
	cache.angle = angle
	cache.x = x
	cache.y = y
	c.listCache.append(&cache)
}

func (c *ContentObj) appendRotateReset() {
	var cache cacheContentRotate
	cache.isReset = true
	c.listCache.append(&cache)
}

//ContentObjCalTextHeight calculates height of text.
func ContentObjCalTextHeight(fontsize int) float64 {
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
