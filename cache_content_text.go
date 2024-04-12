package gopdf

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
)

const defaultCoefLineHeight = float64(1)
const defaultCoefUnderlinePosition = float64(1)
const defaultcoefUnderlineThickness = float64(1)

// ContentTypeCell cell
const ContentTypeCell = 0

// ContentTypeText text
const ContentTypeText = 1

var ErrContentTypeNotFound = errors.New("contentType not found")

type cacheContentText struct {
	//---setup---
	rectangle      *Rect
	textColor      ICacheColorText
	grayFill       float64
	txtColorMode   string
	fontCountIndex int //Curr.FontFontCount+1
	fontSize       float64
	fontStyle      int
	charSpacing    float64
	setXCount      int //จำนวนครั้งที่ใช้ setX
	x, y           float64
	fontSubset     *SubsetFontObj
	pageheight     float64
	contentType    int
	cellOpt        CellOption
	lineWidth      float64
	text           string
	//---result---
	cellWidthPdfUnit, textWidthPdfUnit float64
	cellHeightPdfUnit                  float64
	isPlaceHolder                      bool
}

func (c *cacheContentText) isSame(cache cacheContentText) bool {
	if c.rectangle != nil {
		//if rectangle != nil we assume this is not same content
		return false
	}

	// if both colors are nil we assume them equal
	if ((c.textColor == nil && cache.textColor == nil) ||
		(c.textColor != nil && c.textColor.equal(cache.textColor))) &&
		c.grayFill == cache.grayFill &&
		c.fontCountIndex == cache.fontCountIndex &&
		c.fontSize == cache.fontSize &&
		c.fontStyle == cache.fontStyle &&
		c.charSpacing == cache.charSpacing &&
		c.setXCount == cache.setXCount &&
		c.y == cache.y &&
		c.isPlaceHolder == cache.isPlaceHolder {
		return true
	}

	return false
}

func (c *cacheContentText) setPageHeight(pageheight float64) {
	c.pageheight = pageheight
}

func (c *cacheContentText) pageHeight() float64 {
	return c.pageheight //841.89
}

func convertTypoUnit(val float64, unitsPerEm uint, fontSize float64) float64 {
	val = val * 1000.00 / float64(unitsPerEm)
	return val * fontSize / 1000.0
}

func (c *cacheContentText) calTypoAscender() float64 {
	return convertTypoUnit(float64(c.fontSubset.ttfp.TypoAscender()), c.fontSubset.ttfp.UnitsPerEm(), float64(c.fontSize))
}

func (c *cacheContentText) calTypoDescender() float64 {
	return convertTypoUnit(float64(c.fontSubset.ttfp.TypoDescender()), c.fontSubset.ttfp.UnitsPerEm(), float64(c.fontSize))
}

func (c *cacheContentText) calY() (float64, error) {
	pageHeight := c.pageHeight()
	if c.contentType == ContentTypeText {
		return pageHeight - c.y, nil
	} else if c.contentType == ContentTypeCell {
		y := float64(0.0)
		if c.cellOpt.Align&Bottom == Bottom {
			y = pageHeight - c.y - c.cellHeightPdfUnit - c.calTypoDescender()
		} else if c.cellOpt.Align&Middle == Middle {
			y = pageHeight - c.y - c.cellHeightPdfUnit*0.5 - (c.calTypoDescender()+c.calTypoAscender())*0.5
		} else {
			//top
			y = pageHeight - c.y - c.calTypoAscender()
		}

		return y, nil
	}
	return 0.0, ErrContentTypeNotFound
}

func (c *cacheContentText) calX() (float64, error) {
	if c.contentType == ContentTypeText {
		return c.x, nil
	} else if c.contentType == ContentTypeCell {
		x := float64(0.0)
		if c.cellOpt.Align&Right == Right {
			x = c.x + c.cellWidthPdfUnit - c.textWidthPdfUnit
		} else if c.cellOpt.Align&Center == Center {
			x = c.x + c.cellWidthPdfUnit*0.5 - c.textWidthPdfUnit*0.5
		} else {
			x = c.x
		}
		return x, nil
	}
	return 0.0, ErrContentTypeNotFound
}

// FormatFloatTrim converts a float64 into a string, like Sprintf("%.3f")
// but with trailing zeroes (and possibly ".") removed
func FormatFloatTrim(floatval float64) (formatted string) {
	const precisionFactor = 1000.0
	roundedFontSize := math.Round(precisionFactor*floatval) / precisionFactor
	return strconv.FormatFloat(roundedFontSize, 'f', -1, 64)
}

func (c *cacheContentText) write(w io.Writer, protection *PDFProtection) error {
	x, err := c.calX()
	if err != nil {
		return err
	}
	y, err := c.calY()
	if err != nil {
		return err
	}

	for _, extGStateIndex := range c.cellOpt.extGStateIndexes {
		linkToGSObj := fmt.Sprintf("/GS%d gs\n", extGStateIndex)
		if _, err := io.WriteString(w, linkToGSObj); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "BT\n"); err != nil {
		return err
	}

	fmt.Fprintf(w, "%0.2f %0.2f TD\n", x, y)
	fmt.Fprintf(w, "/F%d %s Tf %s Tc\n", c.fontCountIndex, FormatFloatTrim(c.fontSize), FormatFloatTrim(c.charSpacing))

	if c.txtColorMode == "color" {
		c.textColor.write(w, protection)
	}
	io.WriteString(w, "[<")

	unitsPerEm := int(c.fontSubset.ttfp.UnitsPerEm())
	var leftRune rune
	var leftRuneIndex uint
	for i, r := range c.text {

		glyphindex, err := c.fontSubset.CharIndex(r)
		if err == ErrCharNotFound {
			continue
		} else if err != nil {
			return err
		}

		pairvalPdfUnit := 0
		if i > 0 && c.fontSubset.ttfFontOption.UseKerning { //kerning
			pairval := kern(c.fontSubset, leftRune, r, leftRuneIndex, glyphindex)
			pairvalPdfUnit = convertTTFUnit2PDFUnit(int(pairval), unitsPerEm)
			if pairvalPdfUnit != 0 {
				fmt.Fprintf(w, ">%d<", (-1)*pairvalPdfUnit)
			}
		}

		fmt.Fprintf(w, "%04X", glyphindex)
		leftRune = r
		leftRuneIndex = glyphindex
	}

	io.WriteString(w, ">] TJ\n")
	io.WriteString(w, "ET\n")

	if c.fontStyle&Underline == Underline {
		if err := c.underline(w); err != nil {
			return err
		}
	}

	c.drawBorder(w)

	return nil
}

func (c *cacheContentText) drawBorder(w io.Writer) error {

	//stream.WriteString(fmt.Sprintf("%.2f w\n", 0.1))
	lineOffset := c.lineWidth * 0.5

	if c.cellOpt.Border&Top == Top {

		startX := c.x - lineOffset
		startY := c.pageHeight() - c.y
		endX := c.x + c.cellWidthPdfUnit + lineOffset
		endY := startY
		_, err := fmt.Fprintf(w, "%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY)
		if err != nil {
			return err
		}
	}

	if c.cellOpt.Border&Left == Left {
		startX := c.x
		startY := c.pageHeight() - c.y
		endX := c.x
		endY := startY - c.cellHeightPdfUnit
		_, err := fmt.Fprintf(w, "%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY)
		if err != nil {
			return err
		}
	}

	if c.cellOpt.Border&Right == Right {
		startX := c.x + c.cellWidthPdfUnit
		startY := c.pageHeight() - c.y
		endX := c.x + c.cellWidthPdfUnit
		endY := startY - c.cellHeightPdfUnit
		_, err := fmt.Fprintf(w, "%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY)
		if err != nil {
			return err
		}
	}

	if c.cellOpt.Border&Bottom == Bottom {
		startX := c.x - lineOffset
		startY := c.pageHeight() - c.y - c.cellHeightPdfUnit
		endX := c.x + c.cellWidthPdfUnit + lineOffset
		endY := startY
		_, err := fmt.Fprintf(w, "%0.2f %0.2f m %0.2f %0.2f l s\n", startX, startY, endX, endY)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cacheContentText) underline(w io.Writer) error {
	if c.fontSubset == nil {
		return errors.New("error AppendUnderline not found font")
	}

	coefLineHeight := defaultCoefLineHeight
	if c.cellOpt.CoefLineHeight != 0 {
		coefLineHeight = c.cellOpt.CoefLineHeight
	}

	coefUnderlinePosition := defaultCoefUnderlinePosition
	if c.cellOpt.CoefUnderlinePosition != 0 {
		coefUnderlinePosition = c.cellOpt.CoefUnderlinePosition
	}

	coefUnderlineThickness := defaultcoefUnderlineThickness
	if c.cellOpt.CoefUnderlineThickness != 0 {
		coefUnderlineThickness = c.cellOpt.CoefUnderlineThickness
	}

	ascenderPx := c.fontSubset.GetAscenderPx(c.fontSize)
	descenderPx := -c.fontSubset.GetDescenderPx(c.fontSize)

	contentHeight := ascenderPx + descenderPx
	virtualHeight := coefLineHeight * float64(c.fontSize)
	leading := (contentHeight - virtualHeight) / 2

	baseline := ascenderPx + leading

	underlinePositionPx := c.fontSubset.GetUnderlinePositionPx(c.fontSize) * coefUnderlinePosition
	underlineThicknessPx := c.fontSubset.GetUnderlineThicknessPx(c.fontSize) * coefUnderlineThickness

	yUnderlinePosition := c.pageHeight() - c.y + underlinePositionPx - baseline
	if _, err := fmt.Fprintf(w, "%0.2f %0.2f %0.2f %0.2f re f\n", c.x, yUnderlinePosition, c.cellWidthPdfUnit, underlineThicknessPx); err != nil {
		return err
	}

	return nil
}

func (c *cacheContentText) createContent() (float64, float64, error) {

	cellWidthPdfUnit, cellHeightPdfUnit, textWidthPdfUnit, err := createContent(c.fontSubset, c.text, c.fontSize, c.charSpacing, c.rectangle)
	if err != nil {
		return 0, 0, err
	}
	c.cellWidthPdfUnit = cellWidthPdfUnit
	c.cellHeightPdfUnit = cellHeightPdfUnit
	c.textWidthPdfUnit = textWidthPdfUnit
	return cellWidthPdfUnit, cellHeightPdfUnit, nil
}

func createContent(f *SubsetFontObj, text string, fontSize float64, charSpacing float64, rectangle *Rect) (float64, float64, float64, error) {

	unitsPerEm := int(f.ttfp.UnitsPerEm())
	var leftRune rune
	var leftRuneIndex uint
	sumWidth := int(0)
	//fmt.Printf("unitsPerEm = %d", unitsPerEm)
	for i, r := range text {

		glyphindex, err := f.CharIndex(r)
		if err == ErrCharNotFound {
			continue
		} else if err != nil {
			return 0, 0, 0, err
		}

		pairvalPdfUnit := 0
		if i > 0 && f.ttfFontOption.UseKerning { //kerning
			pairval := kern(f, leftRune, r, leftRuneIndex, glyphindex)
			pairvalPdfUnit = convertTTFUnit2PDFUnit(int(pairval), unitsPerEm)
		}

		width, err := f.CharWidth(r)
		if err != nil {
			return 0, 0, 0, err
		}

		unitsPerPt := float64(unitsPerEm) / fontSize
		spaceWidthInPt := unitsPerPt * charSpacing
		spaceWidthPdfUnit := convertTTFUnit2PDFUnit(int(spaceWidthInPt), unitsPerEm)

		sumWidth += int(width) + int(pairvalPdfUnit) + spaceWidthPdfUnit
		leftRune = r
		leftRuneIndex = glyphindex
	}

	cellWidthPdfUnit := float64(0)
	cellHeightPdfUnit := float64(0)
	if rectangle == nil {
		cellWidthPdfUnit = float64(sumWidth) * (float64(fontSize) / 1000.0)
		typoAscender := convertTypoUnit(float64(f.ttfp.TypoAscender()), f.ttfp.UnitsPerEm(), float64(fontSize))
		typoDescender := convertTypoUnit(float64(f.ttfp.TypoDescender()), f.ttfp.UnitsPerEm(), float64(fontSize))
		cellHeightPdfUnit = typoAscender - typoDescender
	} else {
		cellWidthPdfUnit = rectangle.W
		cellHeightPdfUnit = rectangle.H
	}
	textWidthPdfUnit := float64(sumWidth) * (float64(fontSize) / 1000.0)
	return cellWidthPdfUnit, cellHeightPdfUnit, textWidthPdfUnit, nil
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

// CacheContent Export cacheContent
type CacheContent struct {
	cacheContentText
}

// Setup setup all information for cacheContent
func (c *CacheContent) Setup(rectangle *Rect,
	textColor ICacheColorText,
	grayFill float64,
	fontCountIndex int, //Curr.FontFontCount+1
	fontSize float64,
	fontStyle int,
	charSpacing float64,
	setXCount int, //จำนวนครั้งที่ใช้ setX
	x, y float64,
	fontSubset *SubsetFontObj,
	pageheight float64,
	contentType int,
	cellOpt CellOption,
	lineWidth float64,
) {
	c.cacheContentText = cacheContentText{
		fontSubset:     fontSubset,
		rectangle:      rectangle,
		textColor:      textColor,
		grayFill:       grayFill,
		fontCountIndex: fontCountIndex,
		fontSize:       fontSize,
		fontStyle:      fontStyle,
		charSpacing:    charSpacing,
		setXCount:      setXCount,
		x:              x,
		y:              y,
		pageheight:     pageheight,
		contentType:    ContentTypeCell,
		cellOpt:        cellOpt,
		lineWidth:      lineWidth,
	}
}

// WriteTextToContent write text to content
func (c *CacheContent) WriteTextToContent(text string) {
	c.cacheContentText.text += text
}
