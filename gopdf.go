package gopdf

import (
	"bytes"
	"compress/zlib" // for constants
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/phpdave11/gofpdi"
)

const subsetFont = "SubsetFont"

// the default margin if no margins are set
const defaultMargin = 10.0 //for backward compatible

//GoPdf : A simple library for generating PDF written in Go lang
type GoPdf struct {

	//page Margin
	//leftMargin float64
	//topMargin  float64
	margins Margins

	pdfObjs []IObj
	config  Config
	anchors map[string]anchorOption

	indexOfCatalogObj int

	/*---index ของ obj สำคัญๆ เก็บเพื่อลด loop ตอนค้นหา---*/
	//index ของ obj pages
	indexOfPagesObj int

	//index ของ obj page อันแรก
	indexOfFirstPageObj int

	//ต่ำแหน่งปัจจุบัน
	curr Current

	indexEncodingObjFonts []int
	indexOfContent        int

	//index ของ procset ซึ่งควรจะมีอันเดียว
	indexOfProcSet int

	//IsUnderline bool

	// Buffer for io.Reader compliance
	buf bytes.Buffer

	//pdf PProtection
	pdfProtection   *PDFProtection
	encryptionObjID int

	// content streams only
	compressLevel int

	//info
	isUseInfo bool
	info      *PdfInfo

        //outlines
        outlines *OutlinesObj
        indexOfOutlinesObj int

	// gofpdi free pdf document importer
	fpdi *gofpdi.Importer
}

//SetLineWidth : set line width
func (gp *GoPdf) SetLineWidth(width float64) {
	gp.curr.lineWidth = gp.UnitsToPoints(width)
	gp.getContent().AppendStreamSetLineWidth(gp.UnitsToPoints(width))
}

//SetCompressLevel : set compress Level for content streams
// Possible values for level:
//    -2 HuffmanOnly, -1 DefaultCompression (which is level 6)
//     0 No compression,
//     1 fastest compression, but not very good ratio
//     9 best compression, but slowest
func (gp *GoPdf) SetCompressLevel(level int) {
	errfmt := "compress level too %s, using %s instead\n"
	if level < -2 { //-2 = zlib.HuffmanOnly
		fmt.Fprintf(os.Stderr, errfmt, "small", "DefaultCompression")
		level = zlib.DefaultCompression
	} else if level > zlib.BestCompression {
		fmt.Fprintf(os.Stderr, errfmt, "big", "BestCompression")
		level = zlib.BestCompression
		return
	}
	// sanity check complete
	gp.compressLevel = level
}

//SetNoCompression : compressLevel = 0
func (gp *GoPdf) SetNoCompression() {
	gp.compressLevel = zlib.NoCompression
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
	gp.UnitsToPointsVar(&x1, &y1, &x2, &y2)
	gp.getContent().AppendStreamLine(x1, y1, x2, y2)
}

//RectFromLowerLeft : draw rectangle from lower-left corner (x, y)
func (gp *GoPdf) RectFromLowerLeft(x float64, y float64, wdth float64, hght float64) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)
	gp.getContent().AppendStreamRectangle(x, y, wdth, hght, "")
}

//RectFromUpperLeft : draw rectangle from upper-left corner (x, y)
func (gp *GoPdf) RectFromUpperLeft(x float64, y float64, wdth float64, hght float64) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)
	gp.getContent().AppendStreamRectangle(x, y+hght, wdth, hght, "")
}

//RectFromLowerLeftWithStyle : draw rectangle from lower-left corner (x, y)
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//		D or empty string: draw. This is the default value.
//		F: fill
//		DF or FD: draw and fill
func (gp *GoPdf) RectFromLowerLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)
	gp.getContent().AppendStreamRectangle(x, y, wdth, hght, style)
}

//RectFromUpperLeftWithStyle : draw rectangle from upper-left corner (x, y)
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//		D or empty string: draw. This is the default value.
//		F: fill
//		DF or FD: draw and fill
func (gp *GoPdf) RectFromUpperLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)
	gp.getContent().AppendStreamRectangle(x, y+hght, wdth, hght, style)
}

//Oval : draw oval
func (gp *GoPdf) Oval(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.UnitsToPointsVar(&x1, &y1, &x2, &y2)
	gp.getContent().AppendStreamOval(x1, y1, x2, y2)
}

//Br : new line
func (gp *GoPdf) Br(h float64) {
	gp.UnitsToPointsVar(&h)
	gp.curr.Y += h
	gp.curr.X = gp.margins.Left
}

//SetGrayFill set the grayscale for the fill, takes a float64 between 0.0 and 1.0
func (gp *GoPdf) SetGrayFill(grayScale float64) {
	gp.curr.txtColorMode = "gray"
	gp.curr.grayFill = grayScale
	gp.getContent().AppendStreamSetGrayFill(grayScale)
}

//SetGrayStroke set the grayscale for the stroke, takes a float64 between 0.0 and 1.0
func (gp *GoPdf) SetGrayStroke(grayScale float64) {
	gp.curr.grayStroke = grayScale
	gp.getContent().AppendStreamSetGrayStroke(grayScale)
}

//SetX : set current position X
func (gp *GoPdf) SetX(x float64) {
	gp.UnitsToPointsVar(&x)
	gp.curr.setXCount++
	gp.curr.X = x
}

//GetX : get current position X
func (gp *GoPdf) GetX() float64 {
	return gp.PointsToUnits(gp.curr.X)
}

//SetY : set current position y
func (gp *GoPdf) SetY(y float64) {
	gp.UnitsToPointsVar(&y)
	gp.curr.Y = y
}

//GetY : get current position y
func (gp *GoPdf) GetY() float64 {
	return gp.PointsToUnits(gp.curr.Y)
}

//ImageByHolder : draw image by ImageHolder
func (gp *GoPdf) ImageByHolder(img ImageHolder, x float64, y float64, rect *Rect) error {
	gp.UnitsToPointsVar(&x, &y)
	rect = rect.UnitsToPoints(gp.config.Unit)

	return gp.imageByHolder(img, x, y, rect)
}

func (gp *GoPdf) imageByHolder(img ImageHolder, x float64, y float64, rect *Rect) error {
	cacheImageIndex := -1
	for _, imgcache := range gp.curr.ImgCaches {
		if img.ID() == imgcache.Path {
			cacheImageIndex = imgcache.Index
			break
		}
	}

	if cacheImageIndex == -1 { //new image

		//create img object
		imgobj := new(ImageObj)
		imgobj.init(func() *GoPdf {
			return gp
		})
		imgobj.setProtection(gp.protection())

		err := imgobj.SetImage(img)
		if err != nil {
			return err
		}

		var imgRect *Rect
		if rect == nil {
			imgRect, err = imgobj.getRect()
			if err != nil {
				return err
			}
		} else {
			imgRect = rect
		}

		err = imgobj.parse()
		if err != nil {
			return err
		}
		index := gp.addObj(imgobj)
		if gp.indexOfProcSet != -1 {
			//ยัดรูป
			procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
			gp.getContent().AppendStreamImage(gp.curr.CountOfImg, x, y, imgRect)
			procset.RelateXobjs = append(procset.RelateXobjs, RelateXobject{IndexOfObj: index})
			//เก็บข้อมูลรูปเอาไว้
			var imgcache ImageCache
			imgcache.Index = gp.curr.CountOfImg
			imgcache.Path = img.ID()
			imgcache.Rect = imgRect
			gp.curr.ImgCaches = append(gp.curr.ImgCaches, imgcache)
			gp.curr.CountOfImg++
		}

		if imgobj.haveSMask() {
			smaskObj, err := imgobj.createSMask()
			if err != nil {
				return err
			}
			imgobj.imginfo.smarkObjID = gp.addObj(smaskObj)
		}

		if imgobj.isColspaceIndexed() {
			dRGB, err := imgobj.createDeviceRGB()
			if err != nil {
				return err
			}
			dRGB.getRoot = func() *GoPdf {
				return gp
			}
			imgobj.imginfo.deviceRGBObjID = gp.addObj(dRGB)
		}

	} else { //same img
		var imgRect *Rect
		if rect == nil {
			imgRect = gp.curr.ImgCaches[cacheImageIndex].Rect
		} else {
			imgRect = rect
		}
		gp.getContent().AppendStreamImage(cacheImageIndex, x, y, imgRect)
	}
	return nil
}

//Image : draw image
func (gp *GoPdf) Image(picPath string, x float64, y float64, rect *Rect) error {
	gp.UnitsToPointsVar(&x, &y)
	rect = rect.UnitsToPoints(gp.config.Unit)
	imgh, err := ImageHolderByPath(picPath)
	if err != nil {
		return err
	}
	return gp.imageByHolder(imgh, x, y, rect)
}

//AddPage : add new page
func (gp *GoPdf) AddPage() {
	emptyOpt := PageOption{}
	gp.AddPageWithOption(emptyOpt)
}

//AddPageWithOption  : add new page with option
func (gp *GoPdf) AddPageWithOption(opt PageOption) {
	opt.PageSize = opt.PageSize.UnitsToPoints(gp.config.Unit)

	page := new(PageObj)
	page.init(func() *GoPdf {
		return gp
	})

	if !opt.isEmpty() { //use page option
		page.setOption(opt)
		gp.curr.pageSize = opt.PageSize
	} else { //use default
		gp.curr.pageSize = &gp.config.PageSize
	}

	page.ResourcesRelate = strconv.Itoa(gp.indexOfProcSet+1) + " 0 R"
	index := gp.addObj(page)
	if gp.indexOfFirstPageObj == -1 {
		gp.indexOfFirstPageObj = index
	}
	gp.curr.IndexOfPageObj = index

	//reset
	gp.indexOfContent = -1
	gp.resetCurrXY()
}

func (gp *GoPdf) AddOutline(title string) {
        gp.outlines.AddOutline(gp.curr.IndexOfPageObj + 1, title)
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
	gp.outlines = new(OutlinesObj)
	gp.outlines.init(func() *GoPdf {
		return gp
	})
        gp.indexOfCatalogObj = gp.addObj(catalog)
        gp.indexOfPagesObj = gp.addObj(pages)
	gp.indexOfOutlinesObj = gp.addObj(gp.outlines)
	gp.outlines.SetIndexObjOutlines(gp.indexOfOutlinesObj)

	//indexOfProcSet
	procset := new(ProcSetObj)
	procset.init(func() *GoPdf {
		return gp
	})
	gp.indexOfProcSet = gp.addObj(procset)

	if gp.isUseProtection() {
		gp.pdfProtection = gp.createProtection()
	}

}

// SetFontWithStyle : set font style support Regular or Underline
// for Bold|Italic should be loaded apropriate fonts with same styles defined
func (gp *GoPdf) SetFontWithStyle(family string, style int, size int) error {

	found := false
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].getType() == subsetFont {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*SubsetFontObj)
			if ok {
				if sub.GetFamily() == family && sub.GetTtfFontOption().Style == style&^Underline {
					gp.curr.FontSize = size
					gp.curr.FontStyle = style
					gp.curr.FontFontCount = sub.CountOfFont
					gp.curr.FontISubset = sub
					found = true
					break
				}
			}
		}
		i++
	}

	if !found {
		return errors.New("not found font family")
	}

	return nil
}

//SetFont : set font style support "" or "U"
// for "B" and "I" should be loaded apropriate fonts with same styles defined
func (gp *GoPdf) SetFont(family string, style string, size int) error {
	return gp.SetFontWithStyle(family, getConvertedStyle(style), size)
}

//WritePdf : wirte pdf file
func (gp *GoPdf) WritePdf(pdfPath string) error {
	return ioutil.WriteFile(pdfPath, gp.GetBytesPdf(), 0644)
}

func (gp *GoPdf) Write(w io.Writer) error {
	return gp.compilePdf(w)
}

func (gp *GoPdf) Read(p []byte) (int, error) {
	if gp.buf.Len() == 0 && gp.buf.Cap() == 0 {
		if err := gp.compilePdf(&gp.buf); err != nil {
			return 0, err
		}
	}
	return gp.buf.Read(p)
}

// Close clears the gopdf buffer.
func (gp *GoPdf) Close() error {
	gp.buf = bytes.Buffer{}
	return nil
}

func (gp *GoPdf) compilePdf(w io.Writer) error {
	gp.prepare()
	err := gp.Close()
	if err != nil {
		return err
	}
	max := len(gp.pdfObjs)
	writer := newCountingWriter(w)
	//io.WriteString(w, "%PDF-1.7\n\n")
	fmt.Fprint(writer, "%PDF-1.7\n\n")
	linelens := make([]int, max)
	i := 0

	for i < max {
		objID := i + 1
		linelens[i] = writer.offset
		pdfObj := gp.pdfObjs[i]
		fmt.Fprintf(writer, "%d 0 obj\n", objID)
		pdfObj.write(writer, objID)
		io.WriteString(writer, "endobj\n\n")
		i++
	}
	gp.xref(writer, writer.offset, linelens, i)
	return nil
}

type (
	countingWriter struct {
		offset int
		writer io.Writer
	}
)

func newCountingWriter(w io.Writer) *countingWriter {
	return &countingWriter{writer: w}
}

func (cw *countingWriter) Write(b []byte) (int, error) {
	n, err := cw.writer.Write(b)
	cw.offset += n
	return n, err
}

//GetBytesPdfReturnErr : get bytes of pdf file
func (gp *GoPdf) GetBytesPdfReturnErr() ([]byte, error) {
	err := gp.Close()
	if err != nil {
		return nil, err
	}
	err = gp.compilePdf(&gp.buf)
	return gp.buf.Bytes(), err
}

//GetBytesPdf : get bytes of pdf file
func (gp *GoPdf) GetBytesPdf() []byte {
	b, err := gp.GetBytesPdfReturnErr()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return b
}

//Text write text start at current x,y ( current y is the baseline of text )
func (gp *GoPdf) Text(text string) error {

	err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}

	err = gp.getContent().AppendStreamText(text)
	if err != nil {
		return err
	}

	return nil
}

//CellWithOption create cell of text ( use current x,y is upper-left corner of cell)
func (gp *GoPdf) CellWithOption(rectangle *Rect, text string, opt CellOption) error {
	rectangle = rectangle.UnitsToPoints(gp.config.Unit)
	err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}
	err = gp.getContent().AppendStreamSubsetFont(rectangle, text, opt)
	if err != nil {
		return err
	}
	return nil
}

//Cell : create cell of text ( use current x,y is upper-left corner of cell)
//Note that this has no effect on Rect.H pdf (now). Fix later :-)
func (gp *GoPdf) Cell(rectangle *Rect, text string) error {
	rectangle = rectangle.UnitsToPoints(gp.config.Unit)
	defaultopt := CellOption{
		Align:  Left | Top,
		Border: 0,
		Float:  Right,
	}

	err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}
	err = gp.getContent().AppendStreamSubsetFont(rectangle, text, defaultopt)
	if err != nil {
		return err
	}

	return nil
}

//MultiCell : create of text with line breaks ( use current x,y is upper-left corner of cell)
func (gp *GoPdf) MultiCell(rectangle *Rect, text string) error {
	var line []rune
	x := gp.GetX()
	var totalLineHeight float64
	length := len([]rune(text))

	// get lineHeight
	if err := gp.curr.FontISubset.AddChars(text); err != nil {
		return err
	}
	_, lineHeight, _, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, nil)
	if err != nil {
		return err
	}

	for i, v := range []rune(text) {
		if totalLineHeight+lineHeight > rectangle.H {
			break
		}
		lineWidth, _ := gp.MeasureTextWidth(string(line))
		runeWidth, _ := gp.MeasureTextWidth(string(v))

		if lineWidth+runeWidth > rectangle.W {
			gp.Cell(&Rect{W: rectangle.W, H: lineHeight}, string(line))
			gp.Br(lineHeight)
			gp.SetX(x)
			totalLineHeight = totalLineHeight + lineHeight
			line = nil
		}

		line = append(line, v)

		if i == length-1 {
			gp.Cell(&Rect{W: rectangle.W, H: lineHeight}, string(line))
			gp.Br(lineHeight)
			gp.SetX(x)
		}
	}
	return nil
}

// SplitText splits text into multiple lines based on width.
func (gp *GoPdf) SplitText(text string, width float64) ([]string, error) {
	var lineText []rune
	var lineTexts []string
	utf8Texts := []rune(text)
	utf8TextsLen := len(utf8Texts) // utf8 string quantity
	if utf8TextsLen == 0 {
		return lineTexts, errors.New("empty string")
	}
	for i := 0; i < utf8TextsLen; i++ {
		lineWidth, err := gp.MeasureTextWidth(string(lineText))
		if err != nil {
			return nil, err
		}
		runeWidth, err := gp.MeasureTextWidth(string(utf8Texts[i]))
		if err != nil {
			return nil, err
		}
		if lineWidth+runeWidth > width && utf8Texts[i] != '\n' {
			lineTexts = append(lineTexts, string(lineText))
			lineText = lineText[0:0]
			i--
			continue
		}
		if utf8Texts[i] == '\n' {
			lineTexts = append(lineTexts, string(lineText))
			lineText = lineText[0:0]
			continue
		}
		if i == utf8TextsLen-1 {
			lineText = append(lineText, utf8Texts[i])
			lineTexts = append(lineTexts, string(lineText))
		}
		lineText = append(lineText, utf8Texts[i])

	}
	return lineTexts, nil
}

// ImportPage imports a page and return template id.
// gofpdi code
func (gp *GoPdf) ImportPage(sourceFile string, pageno int, box string) int {
	// Set source file for fpdi
	gp.fpdi.SetSourceFile(sourceFile)

	// gofpdi needs to know where to start the object id at.
	// By default, it starts at 1, but gopdf adds a few objects initially.
	startObjID := gp.GetNextObjectID()

	// Set gofpdi next object ID to  whatever the value of startObjID is
	gp.fpdi.SetNextObjectID(startObjID)

	// Import page
	tpl := gp.fpdi.ImportPage(pageno, box)

	// Import objects into current pdf document
	tplObjIDs := gp.fpdi.PutFormXobjects()

	// Set template names and ids in gopdf
	gp.ImportTemplates(tplObjIDs)

	// Get a map[int]string of the imported objects.
	// The map keys will be the ID of each object.
	imported := gp.fpdi.GetImportedObjects()

	// Import gofpdi objects into gopdf, starting at whatever the value of startObjID is
	gp.ImportObjects(imported, startObjID)

	// Return template ID
	return tpl
}

// ImportPageStream imports page using a stream.
// Return template id after importing.
// gofpdi code
func (gp *GoPdf) ImportPageStream(sourceStream *io.ReadSeeker, pageno int, box string) int {
	// Set source file for fpdi
	gp.fpdi.SetSourceStream(sourceStream)

	// gofpdi needs to know where to start the object id at.
	// By default, it starts at 1, but gopdf adds a few objects initially.
	startObjID := gp.GetNextObjectID()

	// Set gofpdi next object ID to  whatever the value of startObjID is
	gp.fpdi.SetNextObjectID(startObjID)

	// Import page
	tpl := gp.fpdi.ImportPage(pageno, box)

	// Import objects into current pdf document
	tplObjIDs := gp.fpdi.PutFormXobjects()

	// Set template names and ids in gopdf
	gp.ImportTemplates(tplObjIDs)

	// Get a map[int]string of the imported objects.
	// The map keys will be the ID of each object.
	imported := gp.fpdi.GetImportedObjects()

	// Import gofpdi objects into gopdf, starting at whatever the value of startObjID is
	gp.ImportObjects(imported, startObjID)

	// Return template ID
	return tpl
}

// UseImportedTemplate draws an imported PDF page.
func (gp *GoPdf) UseImportedTemplate(tplid int, x float64, y float64, w float64, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)
	// Get template values to draw
	tplName, scaleX, scaleY, tX, tY := gp.fpdi.UseTemplate(tplid, x, y, w, h)
	gp.getContent().AppendStreamImportedTemplate(tplName, scaleX, scaleY, tX, tY)
}

// GetNextObjectID gets the next object ID so that gofpdi knows where to start the object IDs.
func (gp *GoPdf) GetNextObjectID() int {
	return len(gp.pdfObjs) + 1
}

// ImportObjects imports objects from gofpdi into current document.
func (gp *GoPdf) ImportObjects(objs map[int]string, startObjID int) {
	for i := startObjID; i < len(objs)+startObjID; i++ {
		if objs[i] != "" {
			gp.addObj(&ImportedObj{Data: objs[i]})
		}
	}
}

// ImportTemplates names into procset dictionary.
func (gp *GoPdf) ImportTemplates(tpls map[string]int) {
	procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
	for tplName, tplID := range tpls {
		procset.ImportedTemplateIds[tplName] = tplID
	}
}

// AddExternalLink adds a new external link.
func (gp *GoPdf) AddExternalLink(url string, x, y, w, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)
	page := gp.pdfObjs[gp.curr.IndexOfPageObj].(*PageObj)
	page.Links = append(page.Links, linkOption{x, gp.config.PageSize.H - y, w, h, url, ""})
}

// AddInternalLink adds a new internal link.
func (gp *GoPdf) AddInternalLink(anchor string, x, y, w, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)
	page := gp.pdfObjs[gp.curr.IndexOfPageObj].(*PageObj)
	page.Links = append(page.Links, linkOption{x, gp.config.PageSize.H - y, w, h, "", anchor})
}

// SetAnchor creates a new anchor.
func (gp *GoPdf) SetAnchor(name string) {
	y := gp.config.PageSize.H - gp.curr.Y + float64(gp.curr.FontSize)
	gp.anchors[name] = anchorOption{gp.curr.IndexOfPageObj, y}
}

//AddTTFFontByReader add font file
func (gp *GoPdf) AddTTFFontByReader(family string, rd io.Reader) error {
	return gp.AddTTFFontByReaderWithOption(family, rd, defaultTtfFontOption())
}

//AddTTFFontByReaderWithOption add font file
func (gp *GoPdf) AddTTFFontByReaderWithOption(family string, rd io.Reader, option TtfOption) error {
	subsetFont := new(SubsetFontObj)
	subsetFont.init(func() *GoPdf {
		return gp
	})
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFByReader(rd)
	if err != nil {
		return err
	}

	unicodemap := new(UnicodeMap)
	unicodemap.init(func() *GoPdf {
		return gp
	})
	unicodemap.setProtection(gp.protection())
	unicodemap.SetPtrToSubsetFontObj(subsetFont)
	unicodeindex := gp.addObj(unicodemap)

	pdfdic := new(PdfDictionaryObj)
	pdfdic.init(func() *GoPdf {
		return gp
	})
	pdfdic.setProtection(gp.protection())
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
		if !procset.Relates.IsContainsFamilyAndStyle(family, option.Style&^Underline) {
			procset.Relates = append(procset.Relates, RelateFont{Family: family, IndexOfObj: index, CountOfFont: gp.curr.CountOfFont, Style: option.Style &^ Underline})
			subsetFont.CountOfFont = gp.curr.CountOfFont
			gp.curr.CountOfFont++
		}
	}
	return nil
}

//AddTTFFontWithOption : add font file
func (gp *GoPdf) AddTTFFontWithOption(family string, ttfpath string, option TtfOption) error {

	if _, err := os.Stat(ttfpath); os.IsNotExist(err) {
		return err
	}
	data, err := ioutil.ReadFile(ttfpath)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(data)
	return gp.AddTTFFontByReaderWithOption(family, rd, option)
}

//AddTTFFont : add font file
func (gp *GoPdf) AddTTFFont(family string, ttfpath string) error {
	return gp.AddTTFFontWithOption(family, ttfpath, defaultTtfFontOption())
}

//KernOverride override kern value
func (gp *GoPdf) KernOverride(family string, fn FuncKernOverride) error {
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].getType() == subsetFont {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*SubsetFontObj)
			if ok {
				if sub.GetFamily() == family {
					sub.funcKernOverride = fn
					return nil
				}
			}
		}
		i++
	}
	return errors.New("font family not found")
}

//SetTextColor :  function sets the text color
func (gp *GoPdf) SetTextColor(r uint8, g uint8, b uint8) {
	gp.curr.txtColorMode = "color"
	rgb := Rgb{
		r: r,
		g: g,
		b: b,
	}
	gp.curr.setTextColor(rgb)
}

//SetStrokeColor set the color for the stroke
func (gp *GoPdf) SetStrokeColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorStroke(r, g, b)
}

//SetFillColor set the color for the stroke
func (gp *GoPdf) SetFillColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorFill(r, g, b)
}

//MeasureTextWidth : measure Width of text (use current font)
func (gp *GoPdf) MeasureTextWidth(text string) (float64, error) {

	err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, _, textWidthPdfUnit, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, nil)
	if err != nil {
		return 0, err
	}
	return PointsToUnits(gp.config.Unit, textWidthPdfUnit), nil
}

//Curve Draws a Bézier curve (the Bézier curve is tangent to the line between the control points at either end of the curve)
// Parameters:
// - x0, y0: Start point
// - x1, y1: Control point 1
// - x2, y2: Control point 2
// - x3, y3: End point
// - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
func (gp *GoPdf) Curve(x0 float64, y0 float64, x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64, style string) {
	gp.UnitsToPointsVar(&x0, &y0, &x1, &y1, &x2, &y2, &x3, &y3)
	gp.getContent().AppendStreamCurve(x0, y0, x1, y1, x2, y2, x3, y3, style)
}

/*
//SetProtection set permissions as well as user and owner passwords
func (gp *GoPdf) SetProtection(permissions int, userPass []byte, ownerPass []byte) {
	gp.pdfProtection = new(PDFProtection)
	gp.pdfProtection.setProtection(permissions, userPass, ownerPass)
}*/

//SetInfo set Document Information Dictionary
func (gp *GoPdf) SetInfo(info PdfInfo) {
	gp.info = &info
	gp.isUseInfo = true
}

//Rotate rotate text or image
// angle is angle in degrees.
// x, y is rotation center
func (gp *GoPdf) Rotate(angle, x, y float64) {
	gp.UnitsToPointsVar(&x, &y)
	gp.getContent().appendRotate(angle, x, y)
}

//RotateReset reset rotate
func (gp *GoPdf) RotateReset() {
	gp.getContent().appendRotateReset()
}

/*---private---*/

//init
func (gp *GoPdf) init() {

	//default
	gp.margins = Margins{
		Left:   defaultMargin,
		Top:    defaultMargin,
		Right:  defaultMargin,
		Bottom: defaultMargin,
	}

	//init curr
	gp.resetCurrXY()
	gp.curr.IndexOfPageObj = -1
	gp.curr.CountOfFont = 0
	gp.curr.CountOfL = 0
	gp.curr.CountOfImg = 0 //img
	gp.curr.ImgCaches = *new([]ImageCache)
	gp.curr.transparencyMap = make(map[string]Transparency)
	gp.anchors = make(map[string]anchorOption)
	gp.curr.txtColorMode = "gray"

	//init index
	gp.indexOfPagesObj = -1
	gp.indexOfFirstPageObj = -1
	gp.indexOfContent = -1

	//No underline
	//gp.IsUnderline = false
	gp.curr.lineWidth = 1

	// default to zlib.DefaultCompression
	gp.compressLevel = zlib.DefaultCompression

	// change the unit type
	gp.config.PageSize = *gp.config.PageSize.UnitsToPoints(gp.config.Unit)

	// init gofpdi free pdf document importer
	gp.fpdi = gofpdi.NewImporter()
}

func (gp *GoPdf) resetCurrXY() {
	gp.curr.X = gp.margins.Left
	gp.curr.Y = gp.margins.Top
}

// UnitsToPoints converts the units to the documents unit type
func (gp *GoPdf) UnitsToPoints(u float64) float64 {
	return UnitsToPoints(gp.config.Unit, u)
}

// UnitsToPointsVar converts the units to the documents unit type for all variables passed in
func (gp *GoPdf) UnitsToPointsVar(u ...*float64) {
	UnitsToPointsVar(gp.config.Unit, u...)
}

// PointsToUnits converts the points to the documents unit type
func (gp *GoPdf) PointsToUnits(u float64) float64 {
	return PointsToUnits(gp.config.Unit, u)
}

//PointsToUnitsVar converts the points to the documents unit type for all variables passed in
func (gp *GoPdf) PointsToUnitsVar(u ...*float64) {
	PointsToUnitsVar(gp.config.Unit, u...)
}

func (gp *GoPdf) isUseProtection() bool {
	return gp.config.Protection.UseProtection
}

func (gp *GoPdf) createProtection() *PDFProtection {
	var prot PDFProtection
	prot.setProtection(
		gp.config.Protection.Permissions,
		gp.config.Protection.UserPass,
		gp.config.Protection.OwnerPass,
	)
	return &prot
}

func (gp *GoPdf) protection() *PDFProtection {
	return gp.pdfProtection
}

func (gp *GoPdf) prepare() {

	if gp.isUseProtection() {
		encObj := gp.pdfProtection.encryptionObj()
		gp.addObj(encObj)
	}

        if gp.outlines.Count() > 0 {
                catalogObj := gp.pdfObjs[gp.indexOfCatalogObj].(*CatalogObj)
                catalogObj.SetIndexObjOutlines(gp.indexOfOutlinesObj)
        }

	if gp.indexOfPagesObj != -1 {
		indexCurrPage := -1
		var pagesObj *PagesObj
		pagesObj = gp.pdfObjs[gp.indexOfPagesObj].(*PagesObj)
		i := 0 //gp.indexOfFirstPageObj
		max := len(gp.pdfObjs)
		for i < max {
			objtype := gp.pdfObjs[i].getType()
			//fmt.Printf(" objtype = %s , %d \n", objtype , i)
			switch objtype {
			case "Page":
				pagesObj.Kids = fmt.Sprintf("%s %d 0 R ", pagesObj.Kids, i+1)
				pagesObj.PageCount++
				indexCurrPage = i
			case "Content":
				if indexCurrPage != -1 {
					gp.pdfObjs[indexCurrPage].(*PageObj).Contents = fmt.Sprintf("%s %d 0 R ", gp.pdfObjs[indexCurrPage].(*PageObj).Contents, i+1)
				}
			case "Font":
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
			case "Encryption":
				gp.encryptionObjID = i + 1
			}
			i++
		}
	}
}

func (gp *GoPdf) xref(w io.Writer, xrefbyteoffset int, linelens []int, i int) error {

	io.WriteString(w, "xref\n")
	fmt.Fprintf(w, "0 %d\n", i+1)
	io.WriteString(w, "0000000000 65535 f \n")
	j := 0
	max := len(linelens)
	for j < max {
		linelen := linelens[j]
		fmt.Fprintf(w, "%s 00000 n \n", gp.formatXrefline(linelen))
		j++
	}
	io.WriteString(w, "trailer\n")
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "/Size %d\n", max+1)
	io.WriteString(w, "/Root 1 0 R\n")
	if gp.isUseProtection() {
		fmt.Fprintf(w, "/Encrypt %d 0 R\n", gp.encryptionObjID)
		io.WriteString(w, "/ID [()()]\n")
	}
	if gp.isUseInfo {
		gp.writeInfo(w)
	}
	io.WriteString(w, ">>\n")
	io.WriteString(w, "startxref\n")
	fmt.Fprintf(w, "%d", xrefbyteoffset)
	io.WriteString(w, "\n%%EOF\n")

	return nil
}

func (gp *GoPdf) writeInfo(w io.Writer) {
	var zerotime time.Time
	io.WriteString(w, "/Info <<\n")

	if gp.info.Author != "" {
		fmt.Fprintf(w, "/Author <FEFF%s>\n", encodeUtf8(gp.info.Author))
	}

	if gp.info.Title != "" {
		fmt.Fprintf(w, "/Title <FEFF%s>\n", encodeUtf8(gp.info.Title))
	}

	if gp.info.Subject != "" {
		fmt.Fprintf(w, "/Subject <FEFF%s>\n", encodeUtf8(gp.info.Subject))
	}

	if gp.info.Creator != "" {
		fmt.Fprintf(w, "/Creator <FEFF%s>\n", encodeUtf8(gp.info.Creator))
	}

	if gp.info.Producer != "" {
		fmt.Fprintf(w, "/Producer <FEFF%s>\n", encodeUtf8(gp.info.Producer))
	}

	if !zerotime.Equal(gp.info.CreationDate) {
		fmt.Fprintf(w, "/CreationDate(D:%s)>>\n", infodate(gp.info.CreationDate))
	}

	io.WriteString(w, " >>\n")
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

func encodeUtf8(str string) string {
	var buff bytes.Buffer
	for _, r := range str {
		c := fmt.Sprintf("%X", r)
		for len(c) < 4 {
			c = "0" + c
		}
		buff.WriteString(c)
	}
	return buff.String()
}

func infodate(t time.Time) string {
	ft := t.Format("20060102150405-07'00'")
	return ft
}

// SetAlpha sets transparency.
// alpha: 		value from 0 (transparent) to 1 (opaque)
// blendMode:   blend mode, one of the following:
//          		Normal, Multiply, Screen, Overlay, Darken, Lighten, ColorDodge, ColorBurn,
//          		HardLight, SoftLight, Difference, Exclusion, Hue, Saturation, Color, Luminosity
func (gp *GoPdf) SetAlpha(alpha float64, blendModeStr string) error {
	if alpha < 0.0 || alpha > 1.0 {
		return errors.Unwrap(fmt.Errorf("alpha value (0.0 - 1.0) is out of range: %.3f", alpha))
	}

	blendMode, err := getBlendMode(blendModeStr)
	if err != nil {
		return err
	}

	alphaStr := fmt.Sprintf("%.3f", alpha)
	keyStr := fmt.Sprintf("%s %s", alphaStr, blendMode)

	transparency, ok := gp.curr.transparencyMap[keyStr]
	if !ok {
		index, err := gp.addExtGStateObj(&ExtGStateObj{
			ca: alpha,
			CA: alpha,
			BM: blendMode,
		})

		if err != nil {
			return err
		}

		transparency = Transparency{IndexOfExtGState: index + 1}

		gp.curr.transparencyMap[keyStr] = transparency
	}

	gp.curr.transparency = transparency

	return err
}

func getBlendMode(blendModeStr string) (bl string, err error) {
	switch blendModeStr {
	case
		"Hue",
		"Color",
		"Normal",
		"Screen",
		"Darken",
		"Overlay",
		"Lighten",
		"Multiply",
		"ColorBurn",
		"HardLight",
		"SoftLight",
		"Exclusion",
		"Difference",
		"ColorDodge",
		"Saturation",
		"Luminosity":
		bl = "/" + blendModeStr
	case "":
		bl = "/Normal"
	default:
		err = errors.Unwrap(fmt.Errorf("blendMode is unknown"))
	}

	return bl, err
}

func (gp *GoPdf) addExtGStateObj(extGStateObj *ExtGStateObj) (index int, err error) {
	extGStateObj.init(func() *GoPdf {
		return gp
	})
	index = gp.addObj(extGStateObj)

	var pdfObj interface{}
	pdfObj = gp.pdfObjs[gp.indexOfProcSet]

	procset, ok := pdfObj.(*ProcSetObj)
	if !ok {
		err = errors.Unwrap(fmt.Errorf("can't convert pdfobject to procsetobj"))
		return index, err
	}
	procset.ExtGStates = append(procset.ExtGStates, ExtGS{Index: index})

	return index, nil
}

// IsCurrFontContainGlyph defines is current font contains to a glyph
// r:           any rune
func (gp *GoPdf) IsCurrFontContainGlyph (r rune) (bool, error) {
	fontISubset := gp.curr.FontISubset
	if fontISubset == nil {
		return false, nil
	}

	glyphIndex, err := fontISubset.CharCodeToGlyphIndex(r)
	if err != nil {
		return false, err
	}

	if glyphIndex == 0 {
		return false, nil
	}

	return true, nil
}

//tool for validate pdf https://www.pdf-online.com/osa/validate.aspx
