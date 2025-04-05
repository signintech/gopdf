package gopdf

import (
	"bufio"
	"bytes"
	"compress/zlib" // for constants
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/phpdave11/gofpdi"
)

const subsetFont = "SubsetFont"

const colorSpace = "ColorSpace"

// the default margin if no margins are set
const defaultMargin = 10.0 //for backward compatible

var ErrEmptyString = errors.New("empty string")

var ErrMissingFontFamily = errors.New("font family not found")

var ErrMissingColorSpace = errors.New("color space not found")

var ErrExistsColorSpace = errors.New("color space already exists")

var ErrUndefinedCacheContentImage = errors.New("cacheContentImage is undefined")

var ErrInvalidRectangleCoordinates = errors.New("Invalid coordinates for the rectangle")

var ErrInvalidRectangleRadius = errors.New("Radius length cannot exceed rectangle height or width")

// GoPdf : A simple library for generating PDF written in Go lang
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

	//number of pages
	numOfPagesObj int

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
	outlines           *OutlinesObj
	indexOfOutlinesObj int

	//header and footer functions
	headerFunc func()
	footerFunc func()

	// gofpdi free pdf document importer
	fpdi *gofpdi.Importer

	//placeholder text
	placeHolderTexts map[string]([]placeHolderTextInfo)
}

type DrawableRectOptions struct {
	Rect
	X            float64
	Y            float64
	PaintStyle   PaintStyle
	Transparency *Transparency

	extGStateIndexes []int
}

type CropOptions struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type ImageOptions struct {
	DegreeAngle    float64
	VerticalFlip   bool
	HorizontalFlip bool
	X              float64
	Y              float64
	Rect           *Rect
	Mask           *MaskOptions
	Crop           *CropOptions
	Transparency   *Transparency

	extGStateIndexes []int
}
type ImageFromOption struct {
	Format string //jpeg,png
	X      float64
	Y      float64
	Rect   *Rect
}

type MaskOptions struct {
	ImageOptions
	BBox   *[4]float64
	Holder ImageHolder
}

type lineOptions struct {
	extGStateIndexes []int
}

type polygonOptions struct {
	extGStateIndexes []int
}

// SetLineWidth : set line width
func (gp *GoPdf) SetLineWidth(width float64) {
	gp.curr.lineWidth = gp.UnitsToPoints(width)
	gp.getContent().AppendStreamSetLineWidth(gp.UnitsToPoints(width))
}

// SetCompressLevel : set compress Level for content streams
// Possible values for level:
//
//	-2 HuffmanOnly, -1 DefaultCompression (which is level 6)
//	 0 No compression,
//	 1 fastest compression, but not very good ratio
//	 9 best compression, but slowest
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

// SetNoCompression : compressLevel = 0
func (gp *GoPdf) SetNoCompression() {
	gp.compressLevel = zlib.NoCompression
}

// SetLineType : set line type  ("dashed" ,"dotted")
//
//	Usage:
//	pdf.SetLineType("dashed")
//	pdf.Line(50, 200, 550, 200)
//	pdf.SetLineType("dotted")
//	pdf.Line(50, 400, 550, 400)
func (gp *GoPdf) SetLineType(linetype string) {
	gp.getContent().AppendStreamSetLineType(linetype)
}

// SetCustomLineType : set custom line type
//
//	Usage:
//	pdf.SetCustomLineType([]float64{0.8, 0.8}, 0)
//	pdf.Line(50, 200, 550, 200)
func (gp *GoPdf) SetCustomLineType(dashArray []float64, dashPhase float64) {
	for i := range dashArray {
		gp.UnitsToPointsVar(&dashArray[i])
	}
	gp.UnitsToPointsVar(&dashPhase)
	gp.getContent().AppendStreamSetCustomLineType(dashArray, dashPhase)
}

// Line : draw line
//
//	Usage:
//	pdf.SetTransparency(gopdf.Transparency{Alpha: 0.5,BlendModeType: gopdf.ColorBurn})
//	pdf.SetLineType("dotted")
//	pdf.SetStrokeColor(255, 0, 0)
//	pdf.SetLineWidth(2)
//	pdf.Line(10, 30, 585, 30)
//	pdf.ClearTransparency()
func (gp *GoPdf) Line(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.UnitsToPointsVar(&x1, &y1, &x2, &y2)
	transparency, err := gp.getCachedTransparency(nil)
	if err != nil {
		transparency = nil
	}
	var opts = lineOptions{}
	if transparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, transparency.extGStateIndex)
	}
	gp.getContent().AppendStreamLine(x1, y1, x2, y2, opts)
}

// RectFromLowerLeft : draw rectangle from lower-left corner (x, y)
func (gp *GoPdf) RectFromLowerLeft(x float64, y float64, wdth float64, hght float64) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)

	opts := DrawableRectOptions{
		X:          x,
		Y:          y,
		PaintStyle: DrawPaintStyle,
		Rect:       Rect{W: wdth, H: hght},
	}

	gp.getContent().AppendStreamRectangle(opts)
}

// RectFromUpperLeft : draw rectangle from upper-left corner (x, y)
func (gp *GoPdf) RectFromUpperLeft(x float64, y float64, wdth float64, hght float64) {
	gp.UnitsToPointsVar(&x, &y, &wdth, &hght)

	opts := DrawableRectOptions{
		X:          x,
		Y:          y + hght,
		PaintStyle: DrawPaintStyle,
		Rect:       Rect{W: wdth, H: hght},
	}

	gp.getContent().AppendStreamRectangle(opts)
}

// RectFromLowerLeftWithStyle : draw rectangle from lower-left corner (x, y)
//   - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
func (gp *GoPdf) RectFromLowerLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	opts := DrawableRectOptions{
		X: x,
		Y: y,
		Rect: Rect{
			H: hght,
			W: wdth,
		},
		PaintStyle: parseStyle(style),
	}
	gp.RectFromLowerLeftWithOpts(opts)
}

func (gp *GoPdf) RectFromLowerLeftWithOpts(opts DrawableRectOptions) error {
	gp.UnitsToPointsVar(&opts.X, &opts.Y, &opts.W, &opts.H)

	imageTransparency, err := gp.getCachedTransparency(opts.Transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	gp.getContent().AppendStreamRectangle(opts)

	return nil
}

// RectFromUpperLeftWithStyle : draw rectangle from upper-left corner (x, y)
//   - style: Style of rectangule (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
func (gp *GoPdf) RectFromUpperLeftWithStyle(x float64, y float64, wdth float64, hght float64, style string) {
	opts := DrawableRectOptions{
		X: x,
		Y: y,
		Rect: Rect{
			H: hght,
			W: wdth,
		},
		PaintStyle: parseStyle(style),
	}
	gp.RectFromUpperLeftWithOpts(opts)
}

func (gp *GoPdf) RectFromUpperLeftWithOpts(opts DrawableRectOptions) error {
	gp.UnitsToPointsVar(&opts.X, &opts.Y, &opts.W, &opts.H)

	opts.Y += opts.H

	imageTransparency, err := gp.getCachedTransparency(opts.Transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	gp.getContent().AppendStreamRectangle(opts)

	return nil
}

// Oval : draw oval
func (gp *GoPdf) Oval(x1 float64, y1 float64, x2 float64, y2 float64) {
	gp.UnitsToPointsVar(&x1, &y1, &x2, &y2)
	gp.getContent().AppendStreamOval(x1, y1, x2, y2)
}

// Br : new line
func (gp *GoPdf) Br(h float64) {
	gp.UnitsToPointsVar(&h)
	gp.curr.Y += h
	gp.curr.X = gp.margins.Left
}

// SetGrayFill set the grayscale for the fill, takes a float64 between 0.0 and 1.0
func (gp *GoPdf) SetGrayFill(grayScale float64) {
	gp.curr.txtColorMode = "gray"
	gp.curr.grayFill = grayScale
	gp.getContent().AppendStreamSetGrayFill(grayScale)
}

// SetGrayStroke set the grayscale for the stroke, takes a float64 between 0.0 and 1.0
func (gp *GoPdf) SetGrayStroke(grayScale float64) {
	gp.curr.grayStroke = grayScale
	gp.getContent().AppendStreamSetGrayStroke(grayScale)
}

// SetX : set current position X
func (gp *GoPdf) SetX(x float64) {
	gp.UnitsToPointsVar(&x)
	gp.curr.setXCount++
	gp.curr.X = x
}

// GetX : get current position X
func (gp *GoPdf) GetX() float64 {
	return gp.PointsToUnits(gp.curr.X)
}

// SetNewY : set current position y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element(such as text) to be inserted is 25px,
// because 10<25, you need to add another page and set y to 20px.
// Because of called AddPage(), X is set to MarginLeft, so you should specify X if needed,
// or make sure SetX() is after SetNewY(), or using SetNewXY().
// SetNewYIfNoOffset is more suitable for scenarios where the offset does not change, such as pdf.Image().
func (gp *GoPdf) SetNewY(y float64, h float64) {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if gp.curr.Y+h > gp.curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
}

// SetNewYIfNoOffset : set current position y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element(such as image) to be inserted is 200px,
// because 10<200, you need to add another page and set y to 20px.
// Tips: gp.curr.X and gp.curr.Y do not change when pdf.Image() is called.
func (gp *GoPdf) SetNewYIfNoOffset(y float64, h float64) {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if y+h > gp.curr.pageSize.H-gp.MarginBottom() { // using new y(*y) instead of gp.curr.Y
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
}

// SetNewXY : set current position x and y, and modified y if add a new page.
// Example:
// For example, if the page height is set to 841px, MarginTop is 20px,
// MarginBottom is 10px, and the height of the element to be inserted is 25px,
// because 10<25, you need to add another page and set y to 20px.
// Because of AddPage(), X is set to MarginLeft, so you should specify X if needed,
// or make sure SetX() is after SetNewY().
func (gp *GoPdf) SetNewXY(y float64, x, h float64) {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if gp.curr.Y+h > gp.curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
	gp.SetX(x)
}

/*
//experimental
func (gp *GoPdf) SetNewY(y float64, h float64) float64 {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if gp.curr.Y+h > gp.curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
	return gp.GetY()
}

//experimental
func (gp *GoPdf) SetNewYIfNoOffset(y float64, h float64) float64 {
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if y+h > gp.curr.pageSize.H-gp.MarginBottom() { // using new y(*y) instead of gp.curr.Y
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
	return gp.GetY()
}

//experimental
func (gp *GoPdf) SetNewXY(y float64, x, h float64) float64{
	gp.UnitsToPointsVar(&y)
	gp.UnitsToPointsVar(&h)
	if gp.curr.Y+h > gp.curr.pageSize.H-gp.MarginBottom() {
		gp.AddPage()
		y = gp.MarginTop() // reset to top of the page.
	}
	gp.curr.Y = y
	gp.SetX(x)
	return gp.GetY()
}
*/

// SetY : set current position y
func (gp *GoPdf) SetY(y float64) {
	gp.UnitsToPointsVar(&y)
	gp.curr.Y = y
}

// GetY : get current position y
func (gp *GoPdf) GetY() float64 {
	return gp.PointsToUnits(gp.curr.Y)
}

// SetXY : set current position x and y
func (gp *GoPdf) SetXY(x, y float64) {
	gp.UnitsToPointsVar(&x)
	gp.curr.setXCount++
	gp.curr.X = x

	gp.UnitsToPointsVar(&y)
	gp.curr.Y = y
}

// ImageByHolder : draw image by ImageHolder
func (gp *GoPdf) ImageByHolder(img ImageHolder, x float64, y float64, rect *Rect) error {
	gp.UnitsToPointsVar(&x, &y)

	rect = rect.UnitsToPoints(gp.config.Unit)

	imageOptions := ImageOptions{
		X:    x,
		Y:    y,
		Rect: rect,
	}

	return gp.imageByHolder(img, imageOptions)
}

func (gp *GoPdf) ImageByHolderWithOptions(img ImageHolder, opts ImageOptions) error {
	gp.UnitsToPointsVar(&opts.X, &opts.Y)

	opts.Rect = opts.Rect.UnitsToPoints(gp.config.Unit)

	imageTransparency, err := gp.getCachedTransparency(opts.Transparency)
	if err != nil {
		return err
	}

	if imageTransparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, imageTransparency.extGStateIndex)
	}

	if opts.Mask != nil {
		maskTransparency, err := gp.getCachedTransparency(opts.Mask.ImageOptions.Transparency)
		if err != nil {
			return err
		}

		if maskTransparency != nil {
			opts.Mask.ImageOptions.extGStateIndexes = append(opts.Mask.ImageOptions.extGStateIndexes, maskTransparency.extGStateIndex)
		}

		gp.UnitsToPointsVar(&opts.Mask.ImageOptions.X, &opts.Mask.ImageOptions.Y)
		opts.Mask.ImageOptions.Rect = opts.Mask.ImageOptions.Rect.UnitsToPoints(gp.config.Unit)

		extGStateIndex, err := gp.maskHolder(opts.Mask.Holder, *opts.Mask)
		if err != nil {
			return err
		}

		opts.extGStateIndexes = append(opts.extGStateIndexes, extGStateIndex)
	}

	return gp.imageByHolder(img, opts)
}

func (gp *GoPdf) maskHolder(img ImageHolder, opts MaskOptions) (int, error) {
	var cacheImage *ImageCache
	var cacheContentImage *cacheContentImage

	for _, imgcache := range gp.curr.ImgCaches {
		if img.ID() == imgcache.Path {
			cacheImage = &imgcache
			break
		}
	}

	if cacheImage == nil {
		maskImgobj := &ImageObj{IsMask: true}
		maskImgobj.init(func() *GoPdf {
			return gp
		})
		maskImgobj.setProtection(gp.protection())

		err := maskImgobj.SetImage(img)
		if err != nil {
			return 0, err
		}

		if opts.Rect == nil {
			if opts.Rect, err = maskImgobj.getRect(); err != nil {
				return 0, err
			}
		}

		if err := maskImgobj.parse(); err != nil {
			return 0, err
		}

		if gp.indexOfProcSet != -1 {
			index := gp.addObj(maskImgobj)
			cacheContentImage = gp.getContent().GetCacheContentImage(index, opts.ImageOptions)
			procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
			procset.RelateXobjs = append(procset.RelateXobjs, RelateXobject{IndexOfObj: index})

			imgcache := ImageCache{
				Index: index,
				Path:  img.ID(),
				Rect:  opts.Rect,
			}
			gp.curr.ImgCaches[index] = imgcache
			gp.curr.CountOfImg++
		}
	} else {
		if opts.Rect == nil {
			opts.Rect = gp.curr.ImgCaches[cacheImage.Index].Rect
		}

		cacheContentImage = gp.getContent().GetCacheContentImage(cacheImage.Index, opts.ImageOptions)
	}

	if cacheContentImage != nil {
		extGStateInd, err := gp.createTransparencyXObjectGroup(cacheContentImage, opts)
		if err != nil {
			return 0, err
		}

		return extGStateInd, nil
	}

	return 0, ErrUndefinedCacheContentImage
}

func (gp *GoPdf) createTransparencyXObjectGroup(image *cacheContentImage, opts MaskOptions) (int, error) {
	bbox := opts.BBox
	if bbox == nil {
		bbox = &[4]float64{
			// correct BBox values is [opts.X, gp.curr.pageSize.H - opts.Y - opts.Rect.H, opts.X + opts.Rect.W, gp.curr.pageSize.H - opts.Y]
			// but if compress pdf through ghostscript result file can't open correctly in mac viewer, because mac viewer can't parse BBox value correctly
			// all other viewers parse BBox correctly (like Adobe Acrobat Reader, Chrome, even Internet Explorer)
			// that's why we need to set [0, 0, gp.curr.pageSize.W, gp.curr.pageSize.H]
			-gp.curr.pageSize.W * 2,
			-gp.curr.pageSize.H * 2,
			gp.curr.pageSize.W * 2,
			gp.curr.pageSize.H * 2,
			// Also, Chrome pdf viewer incorrectly recognize BBox value, that's why we need to set twice as much value
			// for every mask element will be displayed
		}
	}

	groupOpts := TransparencyXObjectGroupOptions{
		BBox:             *bbox,
		ExtGStateIndexes: opts.extGStateIndexes,
		XObjects:         []cacheContentImage{*image},
	}

	transparencyXObjectGroup, err := GetCachedTransparencyXObjectGroup(groupOpts, gp)
	if err != nil {
		return 0, err
	}

	sMaskOptions := SMaskOptions{
		Subtype:                       SMaskLuminositySubtype,
		TransparencyXObjectGroupIndex: transparencyXObjectGroup.Index,
	}
	sMask := GetCachedMask(sMaskOptions, gp)

	extGStateOpts := ExtGStateOptions{SMaskIndex: &sMask.Index}
	extGState, err := GetCachedExtGState(extGStateOpts, gp)
	if err != nil {
		return 0, err
	}

	return extGState.Index + 1, nil
}

func (gp *GoPdf) imageByHolder(img ImageHolder, opts ImageOptions) error {
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
		if opts.Mask != nil {
			imgobj.SplittedMask = true
		}

		imgobj.init(func() *GoPdf {
			return gp
		})
		imgobj.setProtection(gp.protection())

		err := imgobj.SetImage(img)
		if err != nil {
			return err
		}

		if opts.Rect == nil {
			if opts.Rect, err = imgobj.getRect(); err != nil {
				return err
			}
		}

		err = imgobj.parse()
		if err != nil {
			return err
		}
		index := gp.addObj(imgobj)
		if gp.indexOfProcSet != -1 {
			//ยัดรูป
			procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)
			gp.getContent().AppendStreamImage(index, opts)
			procset.RelateXobjs = append(procset.RelateXobjs, RelateXobject{IndexOfObj: index})
			//เก็บข้อมูลรูปเอาไว้
			var imgcache ImageCache
			imgcache.Index = index
			imgcache.Path = img.ID()
			imgcache.Rect = opts.Rect
			gp.curr.ImgCaches[index] = imgcache
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
		if opts.Rect == nil {
			opts.Rect = gp.curr.ImgCaches[cacheImageIndex].Rect
		}

		gp.getContent().AppendStreamImage(cacheImageIndex, opts)
	}
	return nil
}

// Image : draw image
func (gp *GoPdf) Image(picPath string, x float64, y float64, rect *Rect) error {
	gp.UnitsToPointsVar(&x, &y)
	rect = rect.UnitsToPoints(gp.config.Unit)
	imgh, err := ImageHolderByPath(picPath)
	if err != nil {
		return err
	}

	imageOptions := ImageOptions{
		X:    x,
		Y:    y,
		Rect: rect,
	}

	return gp.imageByHolder(imgh, imageOptions)
}

func (gp *GoPdf) ImageFrom(img image.Image, x float64, y float64, rect *Rect) error {
	return gp.ImageFromWithOption(img, ImageFromOption{
		Format: "png",
		X:      x,
		Y:      y,
		Rect:   rect,
	})
}

func (gp *GoPdf) ImageFromWithOption(img image.Image, opts ImageFromOption) error {
	if img == nil {
		return errors.New("Invalid image")
	}

	gp.UnitsToPointsVar(&opts.X, &opts.Y)
	opts.Rect = opts.Rect.UnitsToPoints(gp.config.Unit)
	r, w := io.Pipe()
	go func() {
		bw := bufio.NewWriter(w)
		var err error
		switch opts.Format {
		case "png":
			err = png.Encode(bw, img)
		case "jpeg":
			err = jpeg.Encode(bw, img, nil)
		}

		bw.Flush()
		if err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()

	imgh, err := ImageHolderByReader(bufio.NewReader(r))
	if err != nil {
		return err
	}

	imageOptions := ImageOptions{
		X:    opts.X,
		Y:    opts.Y,
		Rect: opts.Rect,
	}

	return gp.imageByHolder(imgh, imageOptions)
}

// AddPage : add new page
func (gp *GoPdf) AddPage() {
	emptyOpt := PageOption{}
	gp.AddPageWithOption(emptyOpt)
}

// AddPageWithOption  : add new page with option
func (gp *GoPdf) AddPageWithOption(opt PageOption) {
	opt.TrimBox = opt.TrimBox.UnitsToPoints(gp.config.Unit)
	opt.PageSize = opt.PageSize.UnitsToPoints(gp.config.Unit)

	page := new(PageObj)
	page.init(func() *GoPdf {
		return gp
	})

	if !opt.isEmpty() { //use page option
		page.setOption(opt)
		gp.curr.pageSize = opt.PageSize

		if opt.isTrimBoxSet() {
			gp.curr.trimBox = opt.TrimBox
		}
	} else { //use default
		gp.curr.pageSize = &gp.config.PageSize
		gp.curr.trimBox = &gp.config.TrimBox
	}

	page.ResourcesRelate = strconv.Itoa(gp.indexOfProcSet+1) + " 0 R"
	index := gp.addObj(page)
	if gp.indexOfFirstPageObj == -1 {
		gp.indexOfFirstPageObj = index
	}
	gp.curr.IndexOfPageObj = index

	gp.numOfPagesObj++

	//reset
	gp.indexOfContent = -1
	gp.resetCurrXY()

	if gp.headerFunc != nil {
		gp.headerFunc()
		gp.resetCurrXY()
	}

	if gp.footerFunc != nil {
		gp.footerFunc()
		gp.resetCurrXY()
	}
}

func (gp *GoPdf) AddOutline(title string) {
	gp.outlines.AddOutline(gp.curr.IndexOfPageObj+1, title)
}

// AddOutlineWithPosition add an outline with position
func (gp *GoPdf) AddOutlineWithPosition(title string) *OutlineObj {
	return gp.outlines.AddOutlinesWithPosition(gp.curr.IndexOfPageObj+1, title, gp.config.PageSize.H-gp.curr.Y+20)
}

// AddHeader - add a header function, if present this will be automatically called by AddPage()
func (gp *GoPdf) AddHeader(f func()) {
	gp.headerFunc = f
}

// AddFooter - add a footer function, if present this will be automatically called by AddPage()
func (gp *GoPdf) AddFooter(f func()) {
	gp.footerFunc = f
}

// Start : init gopdf
func (gp *GoPdf) Start(config Config) {

	gp.start(config)

}

func (gp *GoPdf) StartWithImporter(config Config, importer *gofpdi.Importer) {

	gp.start(config, importer)

}

func (gp *GoPdf) start(config Config, importer ...*gofpdi.Importer) {

	gp.config = config
	gp.init(importer...)
	//init all basic obj
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

	gp.placeHolderTexts = make(map[string][]placeHolderTextInfo)

}

// convertNumericToFloat64 : accept numeric types, return float64-value
func convertNumericToFloat64(size interface{}) (fontSize float64, err error) {
	switch size := size.(type) {
	case float32:
		return float64(size), nil
	case float64:
		return float64(size), nil
	case int:
		return float64(size), nil
	case int16:
		return float64(size), nil
	case int32:
		return float64(size), nil
	case int64:
		return float64(size), nil
	case int8:
		return float64(size), nil
	case uint:
		return float64(size), nil
	case uint16:
		return float64(size), nil
	case uint32:
		return float64(size), nil
	case uint64:
		return float64(size), nil
	case uint8:
		return float64(size), nil
	default:
		return 0.0, fmt.Errorf("fontSize must be of type (u)int* or float*, not %T", size)
	}
}

// SetFontWithStyle : set font style support Regular or Underline
// for Bold|Italic should be loaded appropriate fonts with same styles defined
// size MUST be uint*, int* or float64*
func (gp *GoPdf) SetFontWithStyle(family string, style int, size interface{}) error {
	fontSize, err := convertNumericToFloat64(size)
	if err != nil {
		return err
	}
	found := false
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].getType() == subsetFont {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*SubsetFontObj)
			if ok {
				if sub.GetFamily() == family && sub.GetTtfFontOption().Style == style&^Underline {
					gp.curr.FontSize = fontSize
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
		return ErrMissingFontFamily
	}

	return nil
}

// SetFont : set font style support "" or "U"
// for "B" and "I" should be loaded appropriate fonts with same styles defined
// size MUST be uint*, int* or float64*
func (gp *GoPdf) SetFont(family string, style string, size interface{}) error {
	return gp.SetFontWithStyle(family, getConvertedStyle(style), size)
}

// SetFontSize : set the font size (and only the font size) of the currently
// active font
func (gp *GoPdf) SetFontSize(fontSize float64) error {
	gp.curr.FontSize = fontSize
	return nil
}

// SetCharSpacing : set the character spacing of the currently active font
func (gp *GoPdf) SetCharSpacing(charSpacing float64) error {
	gp.UnitsToPointsVar(&charSpacing)
	gp.curr.CharSpacing = charSpacing
	return nil
}

// WritePdf : write pdf file
func (gp *GoPdf) WritePdf(pdfPath string) error {
	return os.WriteFile(pdfPath, gp.GetBytesPdf(), 0644)
}

// WriteTo implements the io.WriterTo interface and can
// be used to stream the PDF as it is compiled to an io.Writer.
func (gp *GoPdf) WriteTo(w io.Writer) (n int64, err error) {
	return gp.compilePdf(w)
}

// Write streams the pdf as it is compiled to an io.Writer
//
// Deprecated: use the WriteTo method instead.
func (gp *GoPdf) Write(w io.Writer) error {
	_, err := gp.compilePdf(w)
	return err
}

func (gp *GoPdf) Read(p []byte) (int, error) {
	if gp.buf.Len() == 0 && gp.buf.Cap() == 0 {
		if _, err := gp.compilePdf(&gp.buf); err != nil {
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

func (gp *GoPdf) compilePdf(w io.Writer) (n int64, err error) {
	gp.prepare()
	err = gp.Close()
	if err != nil {
		return 0, err
	}
	max := len(gp.pdfObjs)
	writer := newCountingWriter(w)
	fmt.Fprint(writer, "%PDF-1.7\n%����\n\n")
	linelens := make([]int64, max)
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
	return writer.offset, nil
}

type (
	countingWriter struct {
		offset int64
		writer io.Writer
	}
)

func newCountingWriter(w io.Writer) *countingWriter {
	return &countingWriter{writer: w}
}

func (cw *countingWriter) Write(b []byte) (int, error) {
	n, err := cw.writer.Write(b)
	cw.offset += int64(n)
	return n, err
}

// GetBytesPdfReturnErr : get bytes of pdf file
func (gp *GoPdf) GetBytesPdfReturnErr() ([]byte, error) {
	err := gp.Close()
	if err != nil {
		return nil, err
	}
	_, err = gp.compilePdf(&gp.buf)
	return gp.buf.Bytes(), err
}

// GetBytesPdf : get bytes of pdf file
func (gp *GoPdf) GetBytesPdf() []byte {
	b, err := gp.GetBytesPdfReturnErr()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return b
}

// Text write text start at current x,y ( current y is the baseline of text )
func (gp *GoPdf) Text(text string) error {

	text, err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}

	err = gp.getContent().AppendStreamText(text)
	if err != nil {
		return err
	}

	return nil
}

// CellWithOption create cell of text ( use current x,y is upper-left corner of cell)
func (gp *GoPdf) CellWithOption(rectangle *Rect, text string, opt CellOption) error {
	transparency, err := gp.getCachedTransparency(opt.Transparency)
	if err != nil {
		return err
	}

	if transparency != nil {
		opt.extGStateIndexes = append(opt.extGStateIndexes, transparency.extGStateIndex)
	}

	rectangle = rectangle.UnitsToPoints(gp.config.Unit)
	text, err = gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}
	if err := gp.getContent().AppendStreamSubsetFont(rectangle, text, opt); err != nil {
		return err
	}

	return nil
}

// Cell : create cell of text ( use current x,y is upper-left corner of cell)
// Note that this has no effect on Rect.H pdf (now). Fix later :-)
func (gp *GoPdf) Cell(rectangle *Rect, text string) error {
	rectangle = rectangle.UnitsToPoints(gp.config.Unit)
	defaultopt := CellOption{
		Align:  Left | Top,
		Border: 0,
		Float:  Right,
	}

	text, err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}
	err = gp.getContent().AppendStreamSubsetFont(rectangle, text, defaultopt)
	if err != nil {
		return err
	}

	return nil
}

// MultiCell : create of text with line breaks ( use current x,y is upper-left corner of cell)
func (gp *GoPdf) MultiCell(rectangle *Rect, text string) error {
	var line []rune
	x := gp.GetX()
	var totalLineHeight float64
	length := len([]rune(text))

	// get lineHeight
	text, err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}
	_, lineHeight, _, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, gp.curr.CharSpacing, nil)
	if err != nil {
		return err
	}
	gp.PointsToUnitsVar(&lineHeight)

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

// IsFitMultiCell : check whether the rectangle's area is big enough for the text
func (gp *GoPdf) IsFitMultiCell(rectangle *Rect, text string) (bool, float64, error) {
	var line []rune
	var totalLineHeight float64
	length := len([]rune(text))

	// get lineHeight
	text, err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return false, totalLineHeight, err
	}
	_, lineHeight, _, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, gp.curr.CharSpacing, nil)

	if err != nil {
		return false, totalLineHeight, err
	}
	gp.PointsToUnitsVar(&lineHeight)

	for i, v := range []rune(text) {
		if totalLineHeight+lineHeight > rectangle.H {
			return false, totalLineHeight, nil
		}
		lineWidth, _ := gp.MeasureTextWidth(string(line))
		runeWidth, _ := gp.MeasureTextWidth(string(v))

		if lineWidth+runeWidth > rectangle.W {
			totalLineHeight += lineHeight
			line = nil
		}

		line = append(line, v)

		if i == length-1 {
			totalLineHeight += lineHeight
		}
	}

	ok := true
	if totalLineHeight > rectangle.H {
		ok = false
	}

	return ok, totalLineHeight, nil
}

// IsFitMultiCellWithNewline : similar to IsFitMultiCell, but process char newline as Br
func (gp *GoPdf) IsFitMultiCellWithNewline(rectangle *Rect, text string) (bool, float64, error) {
	r := *rectangle
	strs := strings.Split(text, "\n")

	for _, s := range strs {
		ok, height, err := gp.IsFitMultiCell(&r, s)
		if err != nil || !ok {
			return false, 0, err
		}
		r.H -= height
	}

	return true, rectangle.H - r.H, nil
}

// MultiCellWithOption create of text with line breaks ( use current x,y is upper-left corner of cell)
func (gp *GoPdf) MultiCellWithOption(rectangle *Rect, text string, opt CellOption) error {
	if opt.BreakOption == nil {
		opt.BreakOption = &DefaultBreakOption
	}

	transparency, err := gp.getCachedTransparency(opt.Transparency)
	if err != nil {
		return err
	}

	if transparency != nil {
		opt.extGStateIndexes = append(opt.extGStateIndexes, transparency.extGStateIndex)
	}

	x := gp.GetX()

	// get lineHeight
	itext, err := gp.curr.FontISubset.AddChars(text)
	if err != nil {
		return err
	}
	_, lineHeight, _, err := createContent(gp.curr.FontISubset, itext, gp.curr.FontSize, gp.curr.CharSpacing, nil)
	if err != nil {
		return err
	}
	gp.PointsToUnitsVar(&lineHeight)

	textSplits, err := gp.SplitTextWithOption(text, rectangle.W, opt.BreakOption)
	if err != nil {
		return err
	}

	startHeight := rectangle.H
	if l := len(textSplits); l > 1 {
		shiftLines := l / 2
		if l%2 != 0 {
			shiftLines += 1
		}
		startHeight = rectangle.H - (lineHeight+1.5)*float64(shiftLines)
	}

	for _, text := range textSplits {
		gp.CellWithOption(&Rect{W: rectangle.W, H: startHeight}, string(text), opt)
		gp.Br(lineHeight)
		gp.SetX(x)
	}

	return nil
}

// SplitText splits text into multiple lines based on width performing potential mid-word breaks.
func (gp *GoPdf) SplitText(text string, width float64) ([]string, error) {
	return gp.SplitTextWithOption(text, width, &DefaultBreakOption)
}

// SplitTextWithWordWrap behaves the same way SplitText does but performs a word-wrap considering spaces in case
// a text line split would split a word.
func (gp *GoPdf) SplitTextWithWordWrap(text string, width float64) ([]string, error) {
	return gp.SplitTextWithOption(text, width, &BreakOption{
		Mode:           BreakModeIndicatorSensitive,
		BreakIndicator: ' ',
	})
}

// SplitTextWithOption splits a text into multiple lines based on the current font size of the document.
// BreakOptions allow to define the behavior of the split (strict or sensitive). For more information see BreakOption.
func (gp *GoPdf) SplitTextWithOption(text string, width float64, opt *BreakOption) ([]string, error) {
	// fallback to default break option
	if opt == nil {
		opt = &DefaultBreakOption
	}
	var lineText []rune
	var lineTexts []string
	utf8Texts := []rune(text)
	utf8TextsLen := len(utf8Texts) // utf8 string quantity
	if utf8TextsLen == 0 {
		return lineTexts, ErrEmptyString
	}
	separatorWidth, err := gp.MeasureTextWidth(opt.Separator)
	if err != nil {
		return nil, err
	}
	// possible (not conflicting) position of the separator within the currently processed line
	separatorIdx := 0
	for i := 0; i < utf8TextsLen; i++ {
		lineWidth, err := gp.MeasureTextWidth(string(lineText))
		if err != nil {
			return nil, err
		}
		runeWidth, err := gp.MeasureTextWidth(string(utf8Texts[i]))
		if err != nil {
			return nil, err
		}
		// mid-word break required since the max width of the given rect is exceeded
		if lineWidth+runeWidth > width && utf8Texts[i] != '\n' {
			// forceBreak will be set to true in case an indicator sensitive break was not possible which will cause
			// strict break to not exceed the desired width
			forceBreak := false
			if opt.Mode == BreakModeIndicatorSensitive {
				forceBreak = !performIndicatorSensitiveLineBreak(&lineTexts, &lineText, &i, opt)
			}
			// BreakModeStrict breaks immediately with an optionally available separator
			if opt.Mode == BreakModeStrict || forceBreak {
				performStrictLineBreak(&lineTexts, &lineText, &i, separatorIdx, opt)
			}
			continue
		}
		// regular break due to a new line rune
		if utf8Texts[i] == '\n' {
			lineTexts = append(lineTexts, string(lineText))
			lineText = lineText[0:0]
			continue
		}
		// end of text
		if i == utf8TextsLen-1 {
			lineText = append(lineText, utf8Texts[i])
			lineTexts = append(lineTexts, string(lineText))
		}
		// store overall index when separator would still fit in the currently processed text-line
		if opt.HasSeparator() && lineWidth+runeWidth+separatorWidth <= width {
			separatorIdx = i
		}
		lineText = append(lineText, utf8Texts[i])
	}
	return lineTexts, nil
}

// [experimental]
// PlaceHolderText Create a text placehold for fillin text later with function FillInPlaceHoldText.
func (gp *GoPdf) PlaceHolderText(placeHolderName string, placeHolderWidth float64) error {

	//placeHolderText := fmt.Sprintf("{%s}", placeHolderName)
	_, err := gp.curr.FontISubset.AddChars("")
	if err != nil {
		return err
	}

	gp.PointsToUnitsVar(&placeHolderWidth)
	err = gp.getContent().appendStreamPlaceHolderText(placeHolderWidth)
	if err != nil {
		return err
	}

	content := gp.pdfObjs[gp.indexOfContent].(*ContentObj)
	indexInContent := len(content.listCache.caches) - 1
	indexOfContent := gp.indexOfContent
	fontISubset := gp.curr.FontISubset

	gp.placeHolderTexts[placeHolderName] = append(
		gp.placeHolderTexts[placeHolderName],
		placeHolderTextInfo{
			indexOfContent:   indexOfContent,
			indexInContent:   indexInContent,
			fontISubset:      fontISubset,
			placeHolderWidth: placeHolderWidth,
			fontSize:         gp.curr.FontSize,
			charSpacing:      gp.curr.CharSpacing,
		},
	)

	return nil
}

// [experimental]
// fill in text that created by function PlaceHolderText
// align: Left,Right,Center
func (gp *GoPdf) FillInPlaceHoldText(placeHolderName string, text string, align int) error {

	infos, ok := gp.placeHolderTexts[placeHolderName]
	if !ok {
		return errors.New("placeHolderName not found")
	}

	for _, info := range infos {
		content, ok := gp.pdfObjs[info.indexOfContent].(*ContentObj)
		if !ok {
			return errors.New("gp.pdfObjs is not *ContentObj")
		}
		contentText, ok := content.listCache.caches[info.indexInContent].(*cacheContentText)
		if !ok {
			return errors.New("listCache.caches is not *cacheContentText")
		}
		info.fontISubset.AddChars(text)
		contentText.text = text

		//Calculate position
		_, _, textWidthPdfUnit, err := createContent(gp.curr.FontISubset, text, info.fontSize, info.charSpacing, nil)
		if err != nil {
			return err
		}
		width := pointsToUnits(gp.config, textWidthPdfUnit)

		if align == Right {
			diff := info.placeHolderWidth - width
			contentText.x = contentText.x + diff
		} else if align == Center {
			diff := info.placeHolderWidth - width
			contentText.x = contentText.x + diff/2
		}
	}

	return nil
}

func performIndicatorSensitiveLineBreak(lineTexts *[]string, lineText *[]rune, i *int, opt *BreakOption) bool {
	brIdx := breakIndicatorIndex(*lineText, opt.BreakIndicator)
	if brIdx > 0 {
		diff := len(*lineText) - brIdx
		*lineText = (*lineText)[0:brIdx]
		*lineTexts = append(*lineTexts, string(*lineText))
		*lineText = (*lineText)[0:0]
		*i -= diff
		return true
	}
	return false
}

func performStrictLineBreak(lineTexts *[]string, lineText *[]rune, i *int, separatorIdx int, opt *BreakOption) {
	if opt.HasSeparator() && separatorIdx > -1 {
		// trim the line to the last possible index with an appended separator
		trimIdx := *i - separatorIdx
		*lineText = (*lineText)[0 : len(*lineText)-trimIdx]
		// append separator to the line
		*lineText = append(*lineText, []rune(opt.Separator)...)
		*lineTexts = append(*lineTexts, string(*lineText))
		*lineText = (*lineText)[0:0]
		*i = separatorIdx - 1
		return
	}
	*lineTexts = append(*lineTexts, string(*lineText))
	*lineText = (*lineText)[0:0]
	*i--
}

// breakIndicatorIndex returns the index where a text line (i.e. rune slice) can be split "gracefully" by checking on
// the break indicator.
// In case no possible break can be identified -1 is returned.
func breakIndicatorIndex(text []rune, bi rune) int {
	for i := len(text) - 1; i > 0; i-- {
		if text[i] == bi {
			return i
		}
	}
	return -1
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

// GetStreamPageSizes gets the sizes of the pages using a stream
// Returns a map of available pages and its box sizes starting with the first page at index 1 containing a map of boxes containing a map of size values
func (gp *GoPdf) GetStreamPageSizes(sourceStream *io.ReadSeeker) map[int]map[string]map[string]float64 {
	gp.fpdi.SetSourceStream(sourceStream)
	return gp.fpdi.GetPageSizes()
}

// GetPageSizes gets the sizes of the pages of a pdf file1
// Returns a map of available pages and its box sizes starting with the first page at index 1 containing a map of boxes containing a map of size values
func (gp *GoPdf) GetPageSizes(sourceFile string) map[int]map[string]map[string]float64 {
	gp.fpdi.SetSourceFile(sourceFile)
	return gp.fpdi.GetPageSizes()
}

// UseImportedTemplate draws an imported PDF page.
func (gp *GoPdf) UseImportedTemplate(tplid int, x float64, y float64, w float64, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)
	// Get template values to draw
	tplName, scaleX, scaleY, tX, tY := gp.fpdi.UseTemplate(tplid, x, y, w, h)
	gp.getContent().AppendStreamImportedTemplate(tplName, scaleX, scaleY, tX, tY)
}

// ImportPagesFromSource imports pages from a source pdf.
// The source can be a file path, byte slice, or (*)io.ReadSeeker.
func (gp *GoPdf) ImportPagesFromSource(source interface{}, box string) error {
	switch v := source.(type) {
	case string:
		// Set source file for fpdi
		gp.fpdi.SetSourceFile(v)
	case []byte:
		// Set source stream for fpdi
		rs := io.ReadSeeker(bytes.NewReader(v))
		gp.fpdi.SetSourceStream(&rs)
	case io.ReadSeeker:
		// Set source stream for fpdi
		gp.fpdi.SetSourceStream(&v)
	case *io.ReadSeeker:
		// Set source stream for fpdi
		gp.fpdi.SetSourceStream(v)
	default:
		return errors.New("source type not supported")
	}

	// Get number of pages from source file
	pages := gp.fpdi.GetNumPages()

	// Get page sizes from source file
	sizes := gp.fpdi.GetPageSizes()

	for i := 0; i < pages; i++ {
		pageno := i + 1

		// Get the size of the page
		size, ok := sizes[pageno][box]
		if !ok {
			return errors.New("can not get page size")
		}

		// Add a new page to the document
		gp.AddPage()

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

		// Draws the imported template on the current page
		gp.UseImportedTemplate(tpl, 0, 0, size["w"], size["h"])
	}

	return nil
}

// GetNextObjectID gets the next object ID so that gofpdi knows where to start the object IDs.
func (gp *GoPdf) GetNextObjectID() int {
	return len(gp.pdfObjs) + 1
}

// GetNumberOfPages gets the number of pages from the PDF.
func (gp *GoPdf) GetNumberOfPages() int {
	return gp.numOfPagesObj
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

	linkOpt := linkOption{x, gp.config.PageSize.H - y, w, h, url, ""}
	gp.addLink(linkOpt)
}

// AddInternalLink adds a new internal link.
func (gp *GoPdf) AddInternalLink(anchor string, x, y, w, h float64) {
	gp.UnitsToPointsVar(&x, &y, &w, &h)

	linkOpt := linkOption{x, gp.config.PageSize.H - y, w, h, "", anchor}
	gp.addLink(linkOpt)
}

func (gp *GoPdf) addLink(option linkOption) {
	page := gp.pdfObjs[gp.curr.IndexOfPageObj].(*PageObj)
	linkObj := gp.addObj(annotObj{option, func() *GoPdf {
		return gp
	}})
	page.LinkObjIds = append(page.LinkObjIds, linkObj+1)
}

// SetAnchor creates a new anchor.
func (gp *GoPdf) SetAnchor(name string) {
	y := gp.config.PageSize.H - gp.curr.Y + float64(gp.curr.FontSize)
	gp.anchors[name] = anchorOption{gp.curr.IndexOfPageObj, y}
}

// AddTTFFontByReader adds font data by reader.
func (gp *GoPdf) AddTTFFontData(family string, fontData []byte) error {
	return gp.AddTTFFontDataWithOption(family, fontData, defaultTtfFontOption())
}

// AddTTFFontDataWithOption adds font data with option.
func (gp *GoPdf) AddTTFFontDataWithOption(family string, fontData []byte, option TtfOption) error {
	subsetFont := new(SubsetFontObj)
	subsetFont.init(func() *GoPdf {
		return gp
	})
	subsetFont.SetTtfFontOption(option)
	subsetFont.SetFamily(family)
	err := subsetFont.SetTTFData(fontData)
	if err != nil {
		return err
	}

	return gp.setSubsetFontObject(subsetFont, family, option)
}

// AddTTFFontByReader adds font file by reader.
func (gp *GoPdf) AddTTFFontByReader(family string, rd io.Reader) error {
	return gp.AddTTFFontByReaderWithOption(family, rd, defaultTtfFontOption())
}

// AddTTFFontByReaderWithOption adds font file by reader with option.
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

	return gp.setSubsetFontObject(subsetFont, family, option)
}

// setSubsetFontObject sets SubsetFontObj.
// The given SubsetFontObj is expected to be configured in advance.
func (gp *GoPdf) setSubsetFontObject(subsetFont *SubsetFontObj, family string, option TtfOption) error {
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

// AddTTFFontWithOption : add font file
func (gp *GoPdf) AddTTFFontWithOption(family string, ttfpath string, option TtfOption) error {

	if _, err := os.Stat(ttfpath); os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(ttfpath)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(data)
	return gp.AddTTFFontByReaderWithOption(family, rd, option)
}

// AddTTFFont : add font file
func (gp *GoPdf) AddTTFFont(family string, ttfpath string) error {
	return gp.AddTTFFontWithOption(family, ttfpath, defaultTtfFontOption())
}

// KernOverride override kern value
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
	return ErrMissingFontFamily
}

// SetTextColor :  function sets the text color
func (gp *GoPdf) SetTextColor(r uint8, g uint8, b uint8) {
	gp.curr.txtColorMode = "color"
	rgb := cacheContentTextColorRGB{
		r: r,
		g: g,
		b: b,
	}
	gp.curr.setTextColor(rgb)
}

func (gp *GoPdf) SetTextColorCMYK(c, m, y, k uint8) {
	gp.curr.txtColorMode = "color"
	cmyk := cacheContentTextColorCMYK{
		c: c,
		m: m,
		y: y,
		k: k,
	}
	gp.curr.setTextColor(cmyk)
}

// SetStrokeColor set the color for the stroke
func (gp *GoPdf) SetStrokeColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorStroke(r, g, b)
}

// SetFillColor set the color for the stroke
func (gp *GoPdf) SetFillColor(r uint8, g uint8, b uint8) {
	gp.getContent().AppendStreamSetColorFill(r, g, b)
}

// SetStrokeColorCMYK set the color for the stroke in CMYK color mode
func (gp *GoPdf) SetStrokeColorCMYK(c, m, y, k uint8) {
	gp.getContent().AppendStreamSetColorStrokeCMYK(c, m, y, k)
}

// SetFillColorCMYK set the color for the fill in CMYK color mode
func (gp *GoPdf) SetFillColorCMYK(c, m, y, k uint8) {
	gp.getContent().AppendStreamSetColorFillCMYK(c, m, y, k)
}

// MeasureTextWidth : measure Width of text (use current font)
func (gp *GoPdf) MeasureTextWidth(text string) (float64, error) {

	text, err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, _, textWidthPdfUnit, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, gp.curr.CharSpacing, nil)
	if err != nil {
		return 0, err
	}
	return pointsToUnits(gp.config, textWidthPdfUnit), nil
}

// MeasureCellHeightByText : measure Height of cell by text (use current font)
func (gp *GoPdf) MeasureCellHeightByText(text string) (float64, error) {

	text, err := gp.curr.FontISubset.AddChars(text) //AddChars for create CharacterToGlyphIndex
	if err != nil {
		return 0, err
	}

	_, cellHeightPdfUnit, _, err := createContent(gp.curr.FontISubset, text, gp.curr.FontSize, gp.curr.CharSpacing, nil)
	if err != nil {
		return 0, err
	}
	return pointsToUnits(gp.config, cellHeightPdfUnit), nil
}

// Curve Draws a Bézier curve (the Bézier curve is tangent to the line between the control points at either end of the curve)
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

// SetInfo set Document Information Dictionary
func (gp *GoPdf) SetInfo(info PdfInfo) {
	gp.info = &info
	gp.isUseInfo = true
}

// GetInfo get Document Information Dictionary
func (gp *GoPdf) GetInfo() PdfInfo {
	return *gp.info
}

// Rotate rotate text or image
// angle is angle in degrees.
// x, y is rotation center
func (gp *GoPdf) Rotate(angle, x, y float64) {
	gp.UnitsToPointsVar(&x, &y)
	gp.getContent().appendRotate(angle, x, y)
}

// RotateReset reset rotate
func (gp *GoPdf) RotateReset() {
	gp.getContent().appendRotateReset()
}

// Polygon : draw polygon
//   - style: Style of polygon (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
//
// Usage:
//
//	 pdf.SetStrokeColor(255, 0, 0)
//		pdf.SetLineWidth(2)
//		pdf.SetFillColor(0, 255, 0)
//		pdf.Polygon([]gopdf.Point{{X: 10, Y: 30}, {X: 585, Y: 200}, {X: 585, Y: 250}}, "DF")
func (gp *GoPdf) Polygon(points []Point, style string) {

	transparency, err := gp.getCachedTransparency(nil)
	if err != nil {
		transparency = nil
	}

	var opts = polygonOptions{}
	if transparency != nil {
		opts.extGStateIndexes = append(opts.extGStateIndexes, transparency.extGStateIndex)
	}

	var pointReals []Point
	for _, p := range points {
		x := p.X
		y := p.Y
		gp.UnitsToPointsVar(&x, &y)
		pointReals = append(pointReals, Point{X: x, Y: y})
	}
	gp.getContent().AppendStreamPolygon(pointReals, style, opts)
}

// Rectangle : draw rectangle, and add radius input to make a round corner, it helps to calculate the round corner coordinates and use Polygon functions to draw rectangle
//   - style: Style of Rectangle (draw and/or fill: D, F, DF, FD)
//     D or empty string: draw. This is the default value.
//     F: fill
//     DF or FD: draw and fill
//
// Usage:
//
//	 pdf.SetStrokeColor(255, 0, 0)
//		pdf.SetLineWidth(2)
//		pdf.SetFillColor(0, 255, 0)
//		pdf.Rectangle(196.6, 336.8, 398.3, 379.3, "DF", 3, 10)
func (gp *GoPdf) Rectangle(x0 float64, y0 float64, x1 float64, y1 float64, style string, radius float64, radiusPointNum int) error {
	if x1 <= x0 || y1 <= y0 {
		return ErrInvalidRectangleCoordinates
	}
	if radiusPointNum <= 0 || radius <= 0 {
		//draw rectangle without round corner
		points := []Point{}
		points = append(points, Point{X: x0, Y: y0})
		points = append(points, Point{X: x1, Y: y0})
		points = append(points, Point{X: x1, Y: y1})
		points = append(points, Point{X: x0, Y: y1})
		gp.Polygon(points, style)

	} else {

		if radius > (x1-x0) || radius > (y1-y0) {
			return ErrInvalidRectangleCoordinates
		}

		degrees := []float64{}
		angle := float64(90) / float64(radiusPointNum+1)
		accAngle := angle
		for accAngle < float64(90) {
			degrees = append(degrees, accAngle)
			accAngle += angle
		}

		radians := []float64{}
		for _, v := range degrees {
			radians = append(radians, v*math.Pi/180)
		}

		points := []Point{}
		points = append(points, Point{X: x0, Y: (y0 + radius)})
		for _, v := range radians {
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x0 + radius - offsetX
			y := y0 + radius - offsetY
			points = append(points, Point{X: x, Y: y})
		}
		points = append(points, Point{X: (x0 + radius), Y: y0})

		points = append(points, Point{X: (x1 - radius), Y: y0})
		for i := range radians {
			v := radians[len(radians)-1-i]
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x1 - radius + offsetX
			y := y0 + radius - offsetY
			points = append(points, Point{X: x, Y: y})
		}
		points = append(points, Point{X: x1, Y: (y0 + radius)})

		points = append(points, Point{X: x1, Y: (y1 - radius)})
		for _, v := range radians {
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x1 - radius + offsetX
			y := y1 - radius + offsetY
			points = append(points, Point{X: x, Y: y})
		}
		points = append(points, Point{X: (x1 - radius), Y: y1})

		points = append(points, Point{X: (x0 + radius), Y: y1})
		for i := range radians {
			v := radians[len(radians)-1-i]
			offsetX := radius * math.Cos(v)
			offsetY := radius * math.Sin(v)
			x := x0 + radius - offsetX
			y := y1 - radius + offsetY
			points = append(points, Point{X: x, Y: y})
		}
		points = append(points, Point{X: x0, Y: y1 - radius})

		gp.Polygon(points, style)
	}
	return nil
}

/*---private---*/

// init
func (gp *GoPdf) init(importer ...*gofpdi.Importer) {
	gp.pdfObjs = []IObj{}
	gp.buf = bytes.Buffer{}
	gp.indexEncodingObjFonts = []int{}
	gp.pdfProtection = nil
	gp.encryptionObjID = 0
	gp.isUseInfo = false
	gp.info = nil

	//default
	gp.margins = Margins{
		Left:   defaultMargin,
		Top:    defaultMargin,
		Right:  defaultMargin,
		Bottom: defaultMargin,
	}

	//init curr
	gp.resetCurrXY()
	gp.curr = Current{}
	gp.curr.IndexOfPageObj = -1
	gp.curr.CountOfFont = 0
	gp.curr.CountOfL = 0
	gp.curr.CountOfImg = 0                       //img
	gp.curr.ImgCaches = make(map[int]ImageCache) //= *new([]ImageCache)
	gp.curr.sMasksMap = NewSMaskMap()
	gp.curr.extGStatesMap = NewExtGStatesMap()
	gp.curr.transparencyMap = NewTransparencyMap()
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
	gp.config.PageSize = *gp.config.PageSize.unitsToPoints(gp.config)
	gp.config.TrimBox = *gp.config.TrimBox.unitsToPoints(gp.config)

	// init gofpdi free pdf document importer
	gp.fpdi = importerOrDefault(importer...)

}

func importerOrDefault(importer ...*gofpdi.Importer) *gofpdi.Importer {
	if len(importer) != 0 {
		return importer[len(importer)-1]
	}
	return gofpdi.NewImporter()
}

func (gp *GoPdf) resetCurrXY() {
	gp.curr.X = gp.margins.Left
	gp.curr.Y = gp.margins.Top
}

// UnitsToPoints converts the units to the documents unit type
func (gp *GoPdf) UnitsToPoints(u float64) float64 {
	return unitsToPoints(gp.config, u)
}

// UnitsToPointsVar converts the units to the documents unit type for all variables passed in
func (gp *GoPdf) UnitsToPointsVar(u ...*float64) {
	unitsToPointsVar(gp.config, u...)
}

// PointsToUnits converts the points to the documents unit type
func (gp *GoPdf) PointsToUnits(u float64) float64 {
	return pointsToUnits(gp.config, u)
}

// PointsToUnitsVar converts the points to the documents unit type for all variables passed in
func (gp *GoPdf) PointsToUnitsVar(u ...*float64) {
	pointsToUnitsVar(gp.config, u...)
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
		pagesObj := gp.pdfObjs[gp.indexOfPagesObj].(*PagesObj)
		i := 0 //gp.indexOfFirstPageObj
		max := len(gp.pdfObjs)
		for i < max {
			objtype := gp.pdfObjs[i].getType()
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

func (gp *GoPdf) xref(w io.Writer, xrefbyteoffset int64, linelens []int64, i int) error {

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
		fmt.Fprintf(w, "/CreationDate(D:%s)\n", infodate(gp.info.CreationDate))
	}

	io.WriteString(w, " >>\n")
}

// ปรับ xref ให้เป็น 10 หลัก
func (gp *GoPdf) formatXrefline(n int64) string {
	str := strconv.FormatInt(n, 10)
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

// SetTransparency sets transparency.
// alpha: 		value from 0 (transparent) to 1 (opaque)
// blendMode:   blend mode, one of the following:
//
//	Normal, Multiply, Screen, Overlay, Darken, Lighten, ColorDodge, ColorBurn,
//	HardLight, SoftLight, Difference, Exclusion, Hue, Saturation, Color, Luminosity
func (gp *GoPdf) SetTransparency(transparency Transparency) error {
	t, err := gp.saveTransparency(&transparency)
	if err != nil {
		return err
	}

	gp.curr.transparency = t

	return nil
}

func (gp *GoPdf) ClearTransparency() {
	gp.curr.transparency = nil
}

func (gp *GoPdf) getCachedTransparency(transparency *Transparency) (*Transparency, error) {
	if transparency == nil {
		transparency = gp.curr.transparency
	} else {
		cached, err := gp.saveTransparency(transparency)
		if err != nil {
			return nil, err
		}

		transparency = cached
	}

	return transparency, nil
}

func (gp *GoPdf) saveTransparency(transparency *Transparency) (*Transparency, error) {
	cached, ok := gp.curr.transparencyMap.Find(*transparency)
	if ok {
		return &cached, nil
	} else if transparency.Alpha != DefaultAplhaValue {
		bm := transparency.BlendModeType
		opts := ExtGStateOptions{
			BlendMode:     &bm,
			StrokingCA:    &transparency.Alpha,
			NonStrokingCa: &transparency.Alpha,
		}

		extGState, err := GetCachedExtGState(opts, gp)
		if err != nil {
			return nil, err
		}

		transparency.extGStateIndex = extGState.Index + 1

		gp.curr.transparencyMap.Save(*transparency)

		return transparency, nil
	}

	return nil, nil
}

// IsCurrFontContainGlyph defines is current font contains to a glyph
// r:           any rune
func (gp *GoPdf) IsCurrFontContainGlyph(r rune) (bool, error) {
	fontISubset := gp.curr.FontISubset
	if fontISubset == nil {
		return false, nil
	}

	glyphIndex, err := fontISubset.CharCodeToGlyphIndex(r)
	if err == ErrGlyphNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if glyphIndex == 0 {
		return false, nil
	}

	return true, nil
}

// SetPage set current page
func (gp *GoPdf) SetPage(pageno int) error {
	var pageIndex int
	for i := 0; i < len(gp.pdfObjs); i++ {
		switch gp.pdfObjs[i].(type) {
		case *ContentObj:
			pageIndex += 1
			if pageIndex == pageno {
				gp.indexOfContent = i
				return nil
			}
		}
	}

	return errors.New("invalid page number")
}

func (gp *GoPdf) SetColorSpace(name string) error {
	found := false
	i := 0
	max := len(gp.pdfObjs)
	for i < max {
		if gp.pdfObjs[i].getType() == colorSpace {
			obj := gp.pdfObjs[i]
			sub, ok := obj.(*ColorSpaceObj)
			if ok {
				if sub.Name == name {
					gp.curr.IndexOfColorSpaceObj = i
					gp.getContent().appendColorSpace(sub.CountOfSpaceColor)
					found = true
					break
				}
			}
		}
		i++
	}

	if !found {
		return ErrMissingColorSpace
	}

	return nil
}

func (gp *GoPdf) AddColorSpaceRGB(name string, r, g, b uint8) error {
	colorSpace := ColorSpaceObj{}
	colorSpace.Name = name

	colorSpace.SetColorRBG(r, g, b)

	return gp.addColorSpace(&colorSpace)
}

func (gp *GoPdf) AddColorSpaceCMYK(name string, c, m, y, k uint8) error {
	colorSpace := ColorSpaceObj{}
	colorSpace.Name = name

	colorSpace.SetColorCMYK(c, m, y, k)

	return gp.addColorSpace(&colorSpace)
}

func (gp *GoPdf) addColorSpace(colorSpace *ColorSpaceObj) error {
	index := gp.addObj(colorSpace)

	if gp.indexOfProcSet != -1 {
		procset := gp.pdfObjs[gp.indexOfProcSet].(*ProcSetObj)

		for _, relate := range procset.RelateColorSpaces {
			if relate.Name == colorSpace.Name {
				return ErrExistsColorSpace
			}
		}

		procset.RelateColorSpaces = append(procset.RelateColorSpaces, RelateColorSpace{Name: colorSpace.Name, IndexOfObj: index, CountOfColorSpace: gp.curr.CountOfColorSpace})
		colorSpace.CountOfSpaceColor = gp.curr.CountOfColorSpace
		gp.curr.CountOfColorSpace++
	}

	return nil
}

//tool for validate pdf https://www.pdf-online.com/osa/validate.aspx
