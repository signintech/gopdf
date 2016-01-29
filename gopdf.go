package gopdf

import (
	"bytes"
	"errors"
	ioutil "io/ioutil"
	"log"
	"os"
	//"container/list"
	"fmt"
	"strconv"
	"strings"
)

//GoPdf : A simple library for generating PDF written in Go lang
type GoPdf struct {

	//page Margin
	leftMargin float64
	topMargin  float64

	pdfObjs []IObj
	config  Config

	/*---index ของ obj สำคัญๆ เก็บเพื่อลด loop ตอนค้นหา---*/
	//index ของ obj pages
	indexOfPagesObj int

	//index ของ obj page อันแรก
	indexOfFirstPageObj int

	//ต่ำแหน่งปัจจุบัน
	Curr Current

	indexEncodingObjFonts []int
	indexOfContent        int

	//index ของ procset ซึ่งควรจะมีอันเดียว
	indexOfProcSet int

	//IsUnderline bool
}

//SetLineWidth : set line width
func (gp *GoPdf) SetLineWidth(width float64) {
	gp.getContent().AppendStreamSetLineWidth(width)
}

//SetLineType : set line type  ("dashed" ,"dotted")
//  Usage:
//  pdf.SetLineType("dashed")
//  pdf.Line(50, 200, 550, 200)
//  pdf.SetLineType("dotted")
//  pdf.Line(50, 400, 550, 400)
func (gp *GoPdf) SetLineType(linetype string) {
	gp.getContent().AppendStreamSetLineType(linetype)
}

//Line : draw line
func (gp *GoPdf) Line(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.getContent().AppendStreamLine(x1, y1, x2, y2)
}

//Oval : draw oval
func (gp *GoPdf) Oval(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.getContent().AppendStreamOval(x1, y1, x2, y2)
}

//Br : new line
func (gp *GoPdf) Br(h float64) {
	gp.Curr.Y += h
	gp.Curr.X = gp.leftMargin
}

//SetGrayFill set the grayscale for the fill, takes a float64 between 0.0 and 1.0
func (gp *GoPdf) SetGrayFill(grayScale float64) {
	gp.Curr.grayFill = grayScale
	gp.getContent().AppendStreamSetGrayFill(grayScale)
}

//SetGrayStroke set the grayscale for the stroke, takes a float64 between 0.0 and 1.0
func (gp *GoPdf) SetGrayStroke(grayScale float64) {
	gp.Curr.grayStroke = grayScale
	gp.getContent().AppendStreamSetGrayStroke(grayScale)
}

//SetLeftMargin : set left margin
func (gp *GoPdf) SetLeftMargin(margin float64) {
	gp.leftMargin = margin
}

//SetTopMargin : set top margin
func (gp *GoPdf) SetTopMargin(margin float64) {
	gp.topMargin = margin
}

//SetX : set current position X
func (gp *GoPdf) SetX(x float64) {
	gp.Curr.X = x
}

//GetX : get current position X
func (gp *GoPdf) GetX() float64 {
	return gp.Curr.X
}

//SetY : set current position y
func (gp *GoPdf) SetY(y float64) {
	gp.Curr.Y = y
}

//GetY : get current position y
func (gp *GoPdf) GetY() float64 {
	return gp.Curr.Y
}

//Image : draw image
func (gp *GoPdf) Image(picPath string, x float64, y float64, rect *Rect) {

	//check
	cacheImageIndex := -1
	for _, imgcache := range gp.Curr.ImgCaches {
		if picPath == imgcache.Path {
			cacheImageIndex = imgcache.Index
			break
		}
	}

	//create img object
	imgobj := new(ImageObj)
	imgobj.init(func() *GoPdf {
		return gp
	})
	imgobj.SetImagePath(picPath)
	if rect == nil {
		rect = imgobj.GetRect()
	}

	if cacheImageIndex == -1 { //new image

		index := gp.addObj(imgobj)
		if gp.indexOfProcSet != -1 {
			//ยัดรูป
			procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
			gp.getContent().AppendStreamImage(gp.Curr.CountOfImg, x, y, rect)
			procset.RealteXobjs = append(procset.RealteXobjs, RealteXobject{IndexOfObj: index})
			//เก็บข้อมูลรูปเอาไว้
			var imgcache ImageCache
			imgcache.Index = gp.Curr.CountOfImg
			imgcache.Path = picPath
			gp.Curr.ImgCaches = append(gp.Curr.ImgCaches, imgcache)
			gp.Curr.CountOfImg++
		}

	} else { //sagp img
		gp.getContent().AppendStreamImage(cacheImageIndex, x, y, rect)
	}

}

//AddPage : add new page
func (gp *GoPdf) AddPage() {
	page := new(PageObj)
	page.init(func() *GoPdf {
		return gp
	})
	page.ResourcesRelate = strconv.Itoa(gp.indexOfProcSet+1) + " 0 R"
	index := gp.addObj(page)
	if gp.indexOfFirstPageObj == -1 {
		gp.indexOfFirstPageObj = index
	}
	gp.Curr.IndexOfPageObj = index

	//reset
	gp.indexOfContent = -1
	gp.resetCurrXY()
}

//Start : init gopdf
func (gp *GoPdf) Start(config Config) {

	gp.config = config
	gp.init()
	//สร้าง obj พื้นฐาน
	catalog := new(CatalogObj)
	catalog.init(func() *GoPdf {
		return gp
	})
	pages := new(PagesObj)
	pages.init(func() *GoPdf {
		return gp
	})
	gp.addObj(catalog)
	gp.indexOfPagesObj = gp.addObj(pages)

	//indexOfProcSet
	procset := new(ProcSetObj)
	procset.init(func() *GoPdf {
		return gp
	})
	gp.indexOfProcSet = gp.addObj(procset)
}

//SetFont : set font style support "" or "U"
func (gp *GoPdf) SetFont(family string, style string, size int) error {

	found := false
	i := 0
	max := len(gp.indexEncodingObjFonts)
	for i < max {
		ifont := gp.pdfObjs[gp.indexEncodingObjFonts[i]].(*EncodingObj).GetFont()
		if ifont.GetFamily() == family {
			gp.Curr.Font_Size = size
			gp.Curr.Font_Style = style
			gp.Curr.Font_FontCount = gp.pdfObjs[gp.indexEncodingObjFonts[i]+4].(*FontObj).CountOfFont
			gp.Curr.Font_Type = CURRENT_FONT_TYPE_IFONT
			gp.Curr.Font_IFont = ifont
			gp.Curr.Font_ISubset = nil
			found = true
			break
		}
		i++
	}

	if !found { //find SubsetFont
		i = 0
		max = len(gp.pdfObjs)
		for i < max {
			if gp.pdfObjs[i].getType() == "SubsetFont" {
				obj := gp.pdfObjs[i]
				sub, ok := obj.(*SubsetFontObj)
				if ok {
					if sub.GetFamily() == family {
						gp.Curr.Font_Size = size
						gp.Curr.Font_Style = style
						gp.Curr.Font_FontCount = sub.CountOfFont
						gp.Curr.Font_Type = CURRENT_FONT_TYPE_SUBSET
						gp.Curr.Font_IFont = nil
						gp.Curr.Font_ISubset = sub
						found = true
						break
					}
				}
			}
			i++
		}
	}

	if !found {
		return errors.New("not found font family")
	}

	return nil
}

//WritePdf : wirte pdf file
func (gp *GoPdf) WritePdf(pdfPath string) {
	ioutil.WriteFile(pdfPath, gp.GetBytesPdf(), 0644)
}

//GetBytesPdfReturnErr : get bytes of pdf file
func (gp *GoPdf) GetBytesPdfReturnErr() ([]byte, error) {
	gp.prepare()
	buff := new(bytes.Buffer)
	i := 0
	max := len(gp.pdfObjs)
	buff.WriteString("%PDF-1.7\n\n")
	linelens := make([]int, max)
	for i < max {
		linelens[i] = buff.Len()
		pdfObj := gp.pdfObjs[i]
		err := pdfObj.build()
		if err != nil {
			return nil, err
		}
		buff.WriteString(strconv.Itoa(i+1) + " 0 obj\n")
		buffbyte := pdfObj.getObjBuff().Bytes()
		buff.Write(buffbyte)
		buff.WriteString("endobj\n\n")
		i++
	}
	gp.xref(linelens, buff, &i)
	return buff.Bytes(), nil
}

//GetBytesPdf : get bytes of pdf file
func (gp *GoPdf) GetBytesPdf() []byte {
	b, err := gp.GetBytesPdfReturnErr()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return b
}

//Cell : create cell of text
//Note that this has no effect on Rect.H pdf (now). Fix later :-)
func (gp *GoPdf) Cell(rectangle *Rect, text string) {

	//undelineOffset := ContentObj_CalTextHeight(gp.Curr.Font_Size) + 1

	startX := gp.Curr.X
	startY := gp.Curr.Y
	if gp.Curr.Font_Type == CURRENT_FONT_TYPE_IFONT {
		gp.getContent().AppendStream(rectangle, text)
	} else if gp.Curr.Font_Type == CURRENT_FONT_TYPE_SUBSET {
		gp.Curr.Font_ISubset.AddChars(text)
		gp.getContent().AppendStreamSubsetFont(rectangle, text)
	}
	endX := gp.Curr.X
	endY := gp.Curr.Y

	//underline
	if strings.Contains(strings.ToUpper(gp.Curr.Font_Style), "U") {
		//gp.Line(x1,y1+undelineOffset,x2,y2+undelineOffset)
		gp.getContent().AppendUnderline(startX, startY, endX, endY, text)
	}

}

//AddTTFFont : font use subtype font
func (gp *GoPdf) AddTTFFont(family string, ttfpath string) error {

	if _, err := os.Stat(ttfpath); os.IsNotExist(err) {
		return err
	}

	subsetFont := new(SubsetFontObj)
	subsetFont.init(func() *GoPdf {
		return gp
	})
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFByPath(ttfpath)
	if err != nil {
		return err
	}

	unicodemap := new(UnicodeMap)
	unicodemap.init(func() *GoPdf {
		return gp
	})
	unicodemap.SetPtrToSubsetFontObj(subsetFont)
	unicodeindex := gp.addObj(unicodemap)

	pdfdic := new(PdfDictionaryObj)
	pdfdic.init(func() *GoPdf {
		return gp
	})
	pdfdic.SetPtrToSubsetFontObj(subsetFont)
	pdfdicindex := gp.addObj(pdfdic)

	subfontdesc := new(SubfontDescriptorObj)
	subfontdesc.init(func() *GoPdf {
		return gp
	})
	subfontdesc.SetPtrToSubsetFontObj(subsetFont)
	subfontdesc.SetIndexObjPdfDictionary(pdfdicindex)
	subfontdescindex := gp.addObj(subfontdesc)

	cidfont := new(CIDFontObj)
	cidfont.init(func() *GoPdf {
		return gp
	})
	cidfont.SetPtrToSubsetFontObj(subsetFont)
	cidfont.SetIndexObjSubfontDescriptor(subfontdescindex)
	cidindex := gp.addObj(cidfont)

	subsetFont.SetIndexObjCIDFont(cidindex)
	subsetFont.SetIndexObjUnicodeMap(unicodeindex)
	index := gp.addObj(subsetFont) //add หลังสุด

	if gp.indexOfProcSet != -1 {
		procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
		if !procset.Realtes.IsContainsFamily(family) {
			procset.Realtes = append(procset.Realtes, RelateFont{Family: family, IndexOfObj: index, CountOfFont: gp.Curr.CountOfFont})
			subsetFont.CountOfFont = gp.Curr.CountOfFont
			gp.Curr.CountOfFont++
		}
	}
	return nil
}

//AddFont : user embed font in zfont file ( deprecated remove sogp day )
func (gp *GoPdf) AddFont(family string, ifont IFont, zfontpath string) {
	encoding := new(EncodingObj)
	ifont.Init()
	ifont.SetFamily(family)
	encoding.SetFont(ifont)
	gp.indexEncodingObjFonts = append(gp.indexEncodingObjFonts, gp.addObj(encoding))

	fontWidth := new(BasicObj)
	fontWidth.init(func() *GoPdf {
		return gp
	})
	fontWidth.Data = "[" + FontConvertHelper_Cw2Str(ifont.GetCw()) + "]\n"
	gp.addObj(fontWidth) //1

	fontDesc := new(FontDescriptorObj)
	fontDesc.init(func() *GoPdf {
		return gp
	})
	fontDesc.SetFont(ifont)
	gp.addObj(fontDesc) //2

	embedfont := new(EmbedFontObj)
	embedfont.init(func() *GoPdf {
		return gp
	})
	embedfont.SetFont(ifont, zfontpath)
	index := gp.addObj(embedfont) //3

	fontDesc.SetFontFileObjRelate(strconv.Itoa(index+1) + " 0 R")

	//start add font obj
	font := new(FontObj)
	font.init(func() *GoPdf {
		return gp
	})
	font.Family = family
	font.Font = ifont
	index = gp.addObj(font) //4
	if gp.indexOfProcSet != -1 {
		procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
		if !procset.Realtes.IsContainsFamily(family) {
			procset.Realtes = append(procset.Realtes, RelateFont{Family: family, IndexOfObj: index, CountOfFont: gp.Curr.CountOfFont})
			font.CountOfFont = gp.Curr.CountOfFont
			gp.Curr.CountOfFont++
		}
	}
	//end add font obj
}

//SetTextColor :  function sets the text color
func (gp *GoPdf) SetTextColor(r uint8, g uint8, b uint8) {
	rgb := Rgb{
		r: r,
		g: g,
		b: b,
	}
	gp.Curr.setTextColor(rgb)
}

//MeasureTextWidth : measure Width of text (use current font)
func (gp *GoPdf) MeasureTextWidth(text string) (float64, error) {

	fontSize := gp.Curr.Font_Size
	sfont := gp.Curr.Font_ISubset
	sfont.AddChars(text) //AddChars for create CharacterToGlyphIndex

	sumWidth := uint64(0)
	//sum width of each rune.
	for _, r := range text {
		width, err := sfont.CharWidth(r)
		if err != nil {
			return 0, err
		}
		sumWidth = sumWidth + width
	}
	realWidth := float64(sumWidth) * (float64(fontSize) / 1000.0) //convert
	return realWidth, nil
}

/*---private---*/

//init
func (gp *GoPdf) init() {

	//default
	gp.leftMargin = 10.0
	gp.topMargin = 10.0

	//init curr
	gp.resetCurrXY()
	gp.Curr.IndexOfPageObj = -1
	gp.Curr.CountOfFont = 0
	gp.Curr.CountOfL = 0
	gp.Curr.CountOfImg = 0 //img
	gp.Curr.ImgCaches = *new([]ImageCache)

	//init index
	gp.indexOfPagesObj = -1
	gp.indexOfFirstPageObj = -1
	gp.indexOfContent = -1

	//No underline
	//gp.IsUnderline = false

}

func (gp *GoPdf) resetCurrXY() {
	gp.Curr.X = gp.leftMargin
	gp.Curr.Y = gp.topMargin
}

func (gp *GoPdf) prepare() {

	if gp.indexOfPagesObj != -1 {
		indexCurrPage := -1
		var pagesObj *PagesObj
		pagesObj = gp.pdfObjs[gp.indexOfPagesObj].(*PagesObj)
		i := 0 //gp.indexOfFirstPageObj
		max := len(gp.pdfObjs)
		for i < max {
			objtype := gp.pdfObjs[i].getType()
			//fmt.Printf(" objtype = %s , %d \n", objtype , i)
			if objtype == "Page" {
				pagesObj.Kids = fmt.Sprintf("%s %d 0 R ", pagesObj.Kids, i+1)
				pagesObj.PageCount++
				indexCurrPage = i
			} else if objtype == "Content" {
				if indexCurrPage != -1 {
					gp.pdfObjs[indexCurrPage].(*PageObj).Contents = fmt.Sprintf("%s %d 0 R ", gp.pdfObjs[indexCurrPage].(*PageObj).Contents, i+1)
				}
			} else if objtype == "Font" {
				tmpfont := gp.pdfObjs[i].(*FontObj)
				j := 0
				jmax := len(gp.indexEncodingObjFonts)
				for j < jmax {
					tmpencoding := gp.pdfObjs[gp.indexEncodingObjFonts[j]].(*EncodingObj).GetFont()
					if tmpfont.Family == tmpencoding.GetFamily() { //ใส่ ข้อมูลของ embed font
						tmpfont.IsEmbedFont = true
						tmpfont.SetIndexObjEncoding(gp.indexEncodingObjFonts[j] + 1)
						tmpfont.SetIndexObjWidth(gp.indexEncodingObjFonts[j] + 2)
						tmpfont.SetIndexObjFontDescriptor(gp.indexEncodingObjFonts[j] + 3)
						break
					}
					j++
				}
			}
			i++
		}
	}
}

func (gp *GoPdf) xref(linelens []int, buff *bytes.Buffer, i *int) {
	xrefbyteoffset := buff.Len()
	buff.WriteString("xref\n")
	buff.WriteString("0 " + strconv.Itoa((*i)+1) + "\n")
	buff.WriteString("0000000000 65535 f\n")
	j := 0
	max := len(linelens)
	for j < max {
		linelen := linelens[j]
		buff.WriteString(gp.formatXrefline(linelen) + " 00000 n\n")
		j++
	}
	buff.WriteString("trailer\n")
	buff.WriteString("<<\n")
	buff.WriteString("/Size " + strconv.Itoa(max+1) + "\n")
	buff.WriteString("/Root 1 0 R\n")
	buff.WriteString(">>\n")

	buff.WriteString("startxref\n")
	buff.WriteString(strconv.Itoa(xrefbyteoffset))
	buff.WriteString("\n%%EOF\n")

	(*i)++
}

//ปรับ xref ให้เป็น 10 หลัก
func (gp *GoPdf) formatXrefline(n int) string {
	str := strconv.Itoa(n)
	for len(str) < 10 {
		str = "0" + str
	}
	return str
}

func (gp *GoPdf) addObj(iobj IObj) int {
	index := len(gp.pdfObjs)
	gp.pdfObjs = append(gp.pdfObjs, iobj)
	return index
}

func (gp *GoPdf) getContent() *ContentObj {
	var content *ContentObj
	if gp.indexOfContent <= -1 {
		content = new(ContentObj)
		content.init(func() *GoPdf {
			return gp
		})
		gp.indexOfContent = gp.addObj(content)
	} else {
		content = gp.pdfObjs[gp.indexOfContent].(*ContentObj)
	}

	return content
}
