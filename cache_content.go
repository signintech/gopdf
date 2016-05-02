package gopdf

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//ContentTypeCell cell
const ContentTypeCell = 0

//ContentTypeText text
const ContentTypeText = 1

type cacheContent struct {
	//---setup---
	rectangle      *Rect
	textColor      Rgb
	grayFill       float64
	fontCountIndex int //Curr.Font_FontCount+1
	fontSize       int
	fontStyle      string
	setXCount      int
	x, y           float64
	fontSubset     *SubsetFontObj
	pageheight     float64
	contentType    int
	cellOpt        CellOption
	lineWidth      float64
	//---result---
	content                             bytes.Buffer
	text                                bytes.Buffer
	textWidthPdfUnit, textHeightPdfUnit float64
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
		c.setXCount == cache.setXCount &&
		//c.x == cache.x &&
		c.y == cache.y {
		return true
	}

	return false
}

func (c *cacheContent) setPageHeight(pageheight float64) {
	c.pageheight = pageheight
}

func (c *cacheContent) pageHeight() float64 {
	return c.pageheight //841.89
}

func convertTypoUnit(val float64, unitsPerEm uint, fontSize float64) float64 {
	val = val * 1000.00 / float64(unitsPerEm)
	return val * fontSize / 1000.0
}

func (c *cacheContent) calTypoAscender() float64 {
	/*
		typoAsc := float64(c.fontSubset.ttfp.TypoAscender()) * 1000.00 / float64(c.fontSubset.ttfp.UnitsPerEm())
		return typoAsc * float64(c.fontSize) / 1000.0
	*/
	return convertTypoUnit(float64(c.fontSubset.ttfp.TypoAscender()), c.fontSubset.ttfp.UnitsPerEm(), float64(c.fontSize))
}

func (c *cacheContent) calTypoDescender() float64 {
	//typoAsc := float64(c.fontSubset.ttfp.TypoDescender()) * 1000.00 / float64(c.fontSubset.ttfp.UnitsPerEm())
	//return typoAsc * float64(c.fontSize) / 1000.0
	return convertTypoUnit(float64(c.fontSubset.ttfp.TypoDescender()), c.fontSubset.ttfp.UnitsPerEm(), float64(c.fontSize))
}

func (c *cacheContent) calY() (float64, error) {
	pageHeight := c.pageHeight()
	if c.contentType == ContentTypeText {
		return pageHeight - c.y, nil
	} else if c.contentType == ContentTypeCell {

		y := pageHeight - c.y - c.calTypoAscender()
		return y, nil
	}
	return 0.0, errors.New("contentType not found")
}

func (c *cacheContent) calX() (float64, error) {
	return c.x, nil
}

func (c *cacheContent) toStream() (*bytes.Buffer, error) {

	var stream bytes.Buffer
	r := c.textColor.r
	g := c.textColor.g
	b := c.textColor.b
	x, err := c.calX()
	if err != nil {
		return nil, err
	}
	y, err := c.calY()
	if err != nil {
		return nil, err
	}

	stream.WriteString("BT\n")
	stream.WriteString(fmt.Sprintf("%0.2f %0.2f TD\n", x, y))
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

	c.drawBorder(&stream)

	return &stream, nil
}

func (c *cacheContent) drawBorder(stream *bytes.Buffer) error {

	//stream.WriteString(fmt.Sprintf("%.2f w\n", 0.1))
	lineOffset := c.lineWidth * 0.5

	if c.cellOpt.Border&Top == Top {

		startX := c.x - lineOffset
		startY := c.pageHeight() - c.y
		endX := c.x + c.textWidthPdfUnit + lineOffset
		endY := startY
		_, err := stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY))
		if err != nil {
			return err
		}
	}

	if c.cellOpt.Border&Left == Left {
		startX := c.x
		startY := c.pageHeight() - c.y
		endX := c.x
		endY := startY - c.textHeightPdfUnit
		_, err := stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY))
		if err != nil {
			return err
		}
	}

	if c.cellOpt.Border&Right == Right {
		startX := c.x + c.textWidthPdfUnit
		startY := c.pageHeight() - c.y
		endX := c.x + c.textWidthPdfUnit
		endY := startY - c.textHeightPdfUnit
		_, err := stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY))
		if err != nil {
			return err
		}
	}

	if c.cellOpt.Border&Bottom == Bottom {
		startX := c.x - lineOffset
		startY := c.pageHeight() - c.y - c.textHeightPdfUnit
		endX := c.x + c.textWidthPdfUnit + lineOffset
		endY := startY
		_, err := stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cacheContent) underline(startX float64, startY float64, endX float64, endY float64) (*bytes.Buffer, error) {

	if c.fontSubset == nil {
		return nil, errors.New("error AppendUnderline not found font")
	}
	unitsPerEm := float64(c.fontSubset.ttfp.UnitsPerEm())
	h := c.pageHeight()
	ut := float64(c.fontSubset.GetUt())
	up := float64(c.fontSubset.GetUp())
	var buff bytes.Buffer
	textH := ContentObj_CalTextHeight(c.fontSize)
	arg3 := float64(h) - (float64(startY) - ((up / unitsPerEm) * float64(c.fontSize))) - textH
	arg4 := (ut / unitsPerEm) * float64(c.fontSize)
	buff.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f -%0.2f re f\n", startX, arg3, endX-startX, arg4))
	//fmt.Printf("arg3=%f arg4=%f\n", arg3, arg4)

	return &buff, nil
}

func (c *cacheContent) createContent() (float64, float64, error) {

	c.content.Truncate(0) //clear
	textWidthPdfUnit, textHeightPdfUnit, err := createContent(c.fontSubset, c.text.String(), c.fontSize, c.rectangle, &c.content)
	if err != nil {
		return 0, 0, err
	}
	c.textWidthPdfUnit = textWidthPdfUnit
	c.textHeightPdfUnit = textHeightPdfUnit
	return textWidthPdfUnit, textHeightPdfUnit, nil
}

func createContent(f *SubsetFontObj, text string, fontSize int, rectangle *Rect, out *bytes.Buffer) (float64, float64, error) {

	unitsPerEm := int(f.ttfp.UnitsPerEm())
	var leftRune rune
	var leftRuneIndex uint
	sumWidth := int(0)
	//fmt.Printf("unitsPerEm = %d", unitsPerEm)
	for i, r := range text {

		glyphindex, err := f.CharIndex(r)
		if err != nil {
			return 0, 0, err
		}

		pairvalPdfUnit := 0
		if i > 0 && f.ttfFontOption.UseKerning { //kerning
			pairval := kern(f, leftRune, r, leftRuneIndex, glyphindex)
			pairvalPdfUnit = convertTTFUnit2PDFUnit(int(pairval), unitsPerEm)
			if pairvalPdfUnit != 0 && out != nil {
				out.WriteString(fmt.Sprintf(">%d<", (-1)*pairvalPdfUnit))
			}
		}

		if out != nil {
			out.WriteString(fmt.Sprintf("%04X", glyphindex))
		}
		width, err := f.CharWidth(r)
		if err != nil {
			return 0, 0, err
		}

		sumWidth += int(width) + int(pairvalPdfUnit)
		leftRune = r
		leftRuneIndex = glyphindex
	}

	textWidthPdfUnit := float64(0)
	textHeightPdfUnit := float64(0)
	if rectangle == nil {
		textWidthPdfUnit = float64(sumWidth) * (float64(fontSize) / 1000.0)
		typoAscender := convertTypoUnit(float64(f.ttfp.TypoAscender()), f.ttfp.UnitsPerEm(), float64(fontSize))
		typoDescender := convertTypoUnit(float64(f.ttfp.TypoDescender()), f.ttfp.UnitsPerEm(), float64(fontSize))
		textHeightPdfUnit = typoAscender - typoDescender
	} else {
		textWidthPdfUnit = rectangle.W
		textHeightPdfUnit = rectangle.H
	}

	return textWidthPdfUnit, textHeightPdfUnit, nil
}

func kern(f *SubsetFontObj, leftRune rune, rightRune rune, leftIndex uint, rightIndex uint) int16 {

	pairVal := int16(0)
	if haveKerning, kval := f.KernValueByLeft(leftIndex); haveKerning {
		if ok, v := kval.ValueByRight(rightIndex); ok {
			pairVal = v
		}
	}

	if f.funcKernOverride != nil {
		pairVal = f.funcKernOverride(
			leftRune,
			rightRune,
			leftIndex,
			rightIndex,
			pairVal,
		)
	}
	return pairVal
}
