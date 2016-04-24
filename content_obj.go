package gopdf

import (
	"bytes"
	"fmt"
	"log"
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

//AppendStreamSubsetFont add stream of text
func (c *ContentObj) AppendStreamSubsetFont(rectangle *Rect, text string) error {

	textColor := c.getRoot().Curr.textColor()
	grayFill := c.getRoot().Curr.grayFill
	fontCountIndex := c.getRoot().Curr.Font_FontCount + 1
	fontSize := c.getRoot().Curr.Font_Size
	fontStyle := c.getRoot().Curr.Font_Style
	x := c.getRoot().Curr.X
	y := c.getRoot().Curr.Y
	setXCount := c.getRoot().Curr.setXCount
	fontSubset := c.getRoot().Curr.Font_ISubset

	//fmt.Printf("fontSubset = %v", fontSubset.ttfFontOption.UseKerning)

	cache := cacheContent{
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
	}
	var err error
	c.getRoot().Curr.X, c.getRoot().Curr.Y, err = c.listCache.appendTextToCache(cache, text)
	if err != nil {
		return err
	}
	return nil
}

//AppendStream add stream of text
func (c *ContentObj) AppendStream(rectangle *Rect, text string) {

	fontSize := c.getRoot().Curr.Font_Size
	r := c.getRoot().Curr.textColor().r
	g := c.getRoot().Curr.textColor().g
	b := c.getRoot().Curr.textColor().b
	grayFill := c.getRoot().Curr.grayFill

	x := fmt.Sprintf("%0.2f", c.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f", c.getRoot().config.PageSize.H-c.getRoot().Curr.Y-(float64(fontSize)*0.7))

	c.stream.WriteString("BT\n")
	c.stream.WriteString(x + " " + y + " TD\n")
	c.stream.WriteString("/F" + strconv.Itoa(c.getRoot().Curr.Font_FontCount+1) + " " + strconv.Itoa(fontSize) + " Tf\n")
	if r+g+b != 0 {
		rFloat := float64(r) * 0.00392156862745
		gFloat := float64(g) * 0.00392156862745
		bFloat := float64(b) * 0.00392156862745
		rgb := fmt.Sprintf("%0.2f %0.2f %0.2f rg\n", rFloat, gFloat, bFloat)
		c.stream.WriteString(rgb)
	} else {
		c.AppendStreamSetGrayFill(grayFill)
	}
	c.stream.WriteString("(" + text + ") Tj\n")
	c.stream.WriteString("ET\n")
	if rectangle == nil {
		c.getRoot().Curr.X += StrHelperGetStringWidth(text, fontSize, c.getRoot().Curr.Font_IFont)
	} else {
		c.getRoot().Curr.X += rectangle.W
	}

}

//AppendStreamLine append line
func (c *ContentObj) AppendStreamLine(x1 float64, y1 float64, x2 float64, y2 float64) {

	h := c.getRoot().config.PageSize.H
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", x1, h-y1, x2, h-y2))
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

//Deprecated
//AppendUnderline append underline
func (c *ContentObj) AppendUnderline(startX float64, y float64, endX float64, endY float64, text string) {

	h := c.getRoot().config.PageSize.H
	ut := int(0)
	if c.getRoot().Curr.Font_IFont != nil {
		ut = c.getRoot().Curr.Font_IFont.GetUt()
	} else if c.getRoot().Curr.Font_ISubset != nil {
		ut = int(c.getRoot().Curr.Font_ISubset.GetUt())
	} else {
		log.Fatal("error AppendUnderline not found font")
	}

	textH := ContentObj_CalTextHeight(c.getRoot().Curr.Font_Size)
	arg3 := float64(h) - float64(y) - textH - textH*0.07
	arg4 := (float64(ut) / 1000.00) * float64(c.getRoot().Curr.Font_Size)
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f -%0.2f re f\n", startX, arg3, endX-startX, arg4))
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
