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

func (c *ContentObj) init(funcGetRoot func() *GoPdf) {
	c.getRoot = funcGetRoot
}

func (c *ContentObj) build() error {
	//fmt.Printf("%s\n", c.listCache.debug())
	buff, err := c.listCache.toStream()
	if err != nil {
		return err
	}
	//fmt.Printf("%s\n", buff.String())
	//c.stream.WriteTo(buff)
	buff.WriteTo(&c.stream)
	streamlen := c.stream.Len()
	c.buffer.WriteString("<<\n")
	c.buffer.WriteString("/Length " + strconv.Itoa(streamlen) + "\n")
	c.buffer.WriteString(">>\n")
	c.buffer.WriteString("stream\n")
	c.buffer.Write(c.stream.Bytes())
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

	h := c.getRoot().config.PageSize.H
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", x1, h-y1, x2, h-y2))
}

//AppendStreamRectangle : draw rectangle from lower-left corner (x, y) with specif width/height
func (c *ContentObj) AppendStreamRectangle(x float64, y float64, wdth float64, hght float64) {
	h := c.getRoot().config.PageSize.H
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re s\n", x, h-y, wdth, hght))
}

//AppendStreamOval append oval
func (c *ContentObj) AppendStreamOval(x1 float64, y1 float64, x2 float64, y2 float64) {
	h := c.getRoot().config.PageSize.H
	cp := 0.55228                              // Magnification of the control point
	v1 := [2]float64{x1 + (x2-x1)/2, h - y2}   // Vertex of the lower
	v2 := [2]float64{x2, h - (y1 + (y2-y1)/2)} // .. Right
	v3 := [2]float64{x1 + (x2-x1)/2, h - y1}   // .. Upper
	v4 := [2]float64{x1, h - (y1 + (y2-y1)/2)} // .. Left

	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m\n", v1[0], v1[1]))
	c.stream.WriteString(fmt.Sprintf(
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v1[0]+(x2-x1)/2*cp, v1[1], v2[0], v2[1]-(y2-y1)/2*cp, v2[0], v2[1],
	))
	c.stream.WriteString(fmt.Sprintf(
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v2[0], v2[1]+(y2-y1)/2*cp, v3[0]+(x2-x1)/2*cp, v3[1], v3[0], v3[1],
	))
	c.stream.WriteString(fmt.Sprintf(
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v3[0]-(x2-x1)/2*cp, v3[1], v4[0], v4[1]+(y2-y1)/2*cp, v4[0], v4[1],
	))
	c.stream.WriteString(fmt.Sprintf(
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c S\n",
		v4[0], v4[1]-(y2-y1)/2*cp, v1[0]-(x2-x1)/2*cp, v1[1], v1[0], v1[1],
	))
}

//AppendStreamCurve draw curve
// - x0, y0: Start point
// - x1, y1: Control point 1
// - x2, y2: Control point 2
// - x3, y3: End point
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
func (c *ContentObj) AppendStreamCurve(x0 float64, y0 float64, x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64, style string) {
	h := c.getRoot().config.PageSize.H
	//cp := 0.55228
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m\n", x0, h-y0))
	c.stream.WriteString(fmt.Sprintf(
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c",
		x1, h-y1, x2, h-y2, x3, h-y3,
	))

	style = strings.TrimSpace(style)
	op := "S"
	if style == "F" {
		op = "f"
	} else if style == "FD" || style == "DF" {
		op = "B"
	}
	c.stream.WriteString(fmt.Sprintf(" %s\n", op))
}

//AppendStreamSetLineWidth : set line width
func (c *ContentObj) AppendStreamSetLineWidth(w float64) {

	c.stream.WriteString(fmt.Sprintf("%.2f w\n", w))

}

//AppendStreamSetLineType : Set linetype [solid, dashed, dotted]
func (c *ContentObj) AppendStreamSetLineType(t string) {
	switch t {
	case "dashed":
		c.stream.WriteString(fmt.Sprint("[5] 2 d\n"))
	case "dotted":
		c.stream.WriteString(fmt.Sprint("[2 3] 11 d\n"))
	default:
		c.stream.WriteString(fmt.Sprint("[] 0 d\n"))
	}

}

//AppendStreamSetGrayFill  set the grayscale fills
func (c *ContentObj) AppendStreamSetGrayFill(w float64) {
	w = fixRange10(w)
	c.stream.WriteString(fmt.Sprintf("%.2f g\n", w))
}

//AppendStreamSetGrayStroke  set the grayscale stroke
func (c *ContentObj) AppendStreamSetGrayStroke(w float64) {
	w = fixRange10(w)
	c.stream.WriteString(fmt.Sprintf("%.2f G\n", w))
}

//AppendStreamSetColorStroke  set the color stroke
func (c *ContentObj) AppendStreamSetColorStroke(r uint8, g uint8, b uint8) {
	//w = fixRange10(w)
	rFloat := float64(r) * 0.00392156862745
	gFloat := float64(g) * 0.00392156862745
	bFloat := float64(b) * 0.00392156862745
	c.stream.WriteString(fmt.Sprintf("%.2f %.2f %.2f RG\n", rFloat, gFloat, bFloat))
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
