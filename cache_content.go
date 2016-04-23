package gopdf

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type cacheContent struct {
	rectangle      *Rect
	textColor      Rgb
	grayFill       float64
	fontCountIndex int //Curr.Font_FontCount+1
	fontSize       int
	fontStyle      string
	x, y           float64
	fontSubset     *SubsetFontObj
	//
	content          bytes.Buffer
	text             bytes.Buffer
	textWidthPdfUnit float64
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

func (c *cacheContent) pageHeight() float64 {
	return 841.89 //TODO fix this //c.getRoot().config.PageSize.H
}

func (c *cacheContent) toStream() (*bytes.Buffer, error) {
	var stream bytes.Buffer

	pageHeight := c.pageHeight()
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
		//TODO fix this
	}

	stream.WriteString("[<" + c.content.String() + ">] TJ\n")
	stream.WriteString("ET\n")

	if strings.ToUpper(c.fontStyle) == "U" {
		underlineStream, err := c.underline(c.x, c.y, c.x+c.textWidthPdfUnit, c.y)
		if err != nil {
			return nil, err
		}
		_, err = underlineStream.WriteTo(&stream)
		if err != nil {
			return nil, err
		}
	}

	return &stream, nil
}

func (c *cacheContent) underline(startX float64, y float64, endX float64, endY float64) (*bytes.Buffer, error) {

	h := c.pageHeight()
	ut := int(0)

	if c.fontSubset != nil {
		ut = int(c.fontSubset.GetUt())
	} else {
		return nil, errors.New("error AppendUnderline not found font")
	}

	var buff bytes.Buffer
	textH := ContentObj_CalTextHeight(c.fontSize)
	arg3 := float64(h) - float64(y) - textH - textH*0.07
	arg4 := (float64(ut) / 1000.00) * float64(c.fontSize)
	buff.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f -%0.2f re f\n", startX, arg3, endX-startX, arg4))

	return &buff, nil
}

func (c *cacheContent) createContent() (float64, error) {

	sumWidth := uint(0)
	c.content.Truncate(0) //clear
	text := c.text.String()
	for _, r := range text {

		glyphindex, err := c.fontSubset.CharIndex(r)
		if err != nil {
			return 0, err
		}
		c.content.WriteString(fmt.Sprintf("%04X", glyphindex))
		width, err := c.fontSubset.CharWidth(r)
		if err != nil {
			return 0, err
		}
		sumWidth += width
	}

	/*err := c.text.UnreadRune() //move read rune ponter to first
	if err != nil {
		return 0, err
	}*/

	fmt.Printf(">>>>>%s\n", c.content.String())
	textWidthPdfUnit := float64(0)
	if c.rectangle == nil {
		textWidthPdfUnit = float64(sumWidth) * (float64(c.fontSize) / 1000.0)
	} else {
		textWidthPdfUnit = c.rectangle.W
	}
	c.textWidthPdfUnit = textWidthPdfUnit

	return textWidthPdfUnit, nil
}
