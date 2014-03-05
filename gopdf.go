package gopdf

import (
	"bytes"
	ioutil "io/ioutil"
	//"container/list"
	"fmt"
	"strconv"
	"strings"
)

type GoPdf struct {
	
	//page Margin
	leftMargin float64
	topMargin  float64
	
	pdfObjs []IObj
	config Config
	
	/*---index ของ obj สำคัญๆ เก็บเพื่อลด loop ตอนค้นหา---*/
	//index ของ obj pages
	indexOfPagesObj int
	
	//index ของ obj page อันแรก
	indexOfFirstPageObj int
	
	//ต่ำแหน่งปัจจุบัน 
	Curr Current
	
	indexEncodingObjFonts []int
	indexOfContent int
	
	//index ของ procset ซึ่งควรจะมีอันเดียว
	indexOfProcSet int
	
	//IsUnderline bool
}

/*---public---*/

func (me *GoPdf)  SetLineWidth( width float64){
	me.getContent().AppendStreamSetLineWidth(width)
}

//วาดเส้น
func (me *GoPdf) Line(x1 float64 , y1 float64, x2 float64 , y2 float64){
	me.getContent().AppendStreamLine(x1,y1,x2,y2)
}

//ขึ้นบรรทัดใหม่
func (me *GoPdf) Br( h  float64){
	me.Curr.Y += h
	me.Curr.X = me.leftMargin
}

func (me *GoPdf) SetLeftMargin(margin float64){
	me.leftMargin = margin
}

func (me *GoPdf) SetTopMargin(margin float64){
	me.topMargin = margin
}

func (me * GoPdf) SetX(x float64){
	me.Curr.X = x
}

func (me * GoPdf) GetX() float64{
	return me.Curr.X
}

func (me * GoPdf) SetY(y float64){
	me.Curr.Y = y
}

func (me * GoPdf) GetY() float64{
	return me.Curr.Y 
}

func (me * GoPdf) Image (picPath string,x float64,y float64 , rect *Rect ){

	//check
	cacheImageIndex := -1
	for _, imgcache := range me.Curr.ImgCaches {
	 	   if   picPath == imgcache.Path {
			   cacheImageIndex = imgcache.Index
			   break
		   }
	}

	//create img object
	imgobj := new (ImageObj)
	imgobj.Init(func()(*GoPdf){
		return me;
	});
	imgobj.SetImagePath(picPath)
	if rect == nil {
		rect = imgobj.GetRect()
	}

	if cacheImageIndex == -1 {  //new image

		index := me.addObj(imgobj)
		if me.indexOfProcSet != -1 {
			//ยัดรูป
			procset := me.pdfObjs[me.indexOfProcSet].(*ProcSetObj)
			me.getContent().AppendStreamImage( me.Curr.CountOfImg ,x,y,rect)
			procset.RealteXobjs = append(procset.RealteXobjs, RealteXobject{ IndexOfObj : index } )
			//เก็บข้อมูลรูปเอาไว้
			var imgcache ImageCache
			imgcache.Index = me.Curr.CountOfImg
			imgcache.Path = picPath
			me.Curr.ImgCaches = append(me.Curr.ImgCaches,imgcache)
			me.Curr.CountOfImg++
		}

	}else{    //same img
			me.getContent().AppendStreamImage( cacheImageIndex ,x,y,rect)
	}
	
}

//เพิ่ม page
func (me *GoPdf) AddPage() {
	page := new(PageObj)
	page.Init(func()(*GoPdf){
		return me
	})
	page.ResourcesRelate = strconv.Itoa(me.indexOfProcSet + 1) + " 0 R"
	index := me.addObj(page)
	if me.indexOfFirstPageObj == -1 {
		me.indexOfFirstPageObj = index
	}
	me.Curr.IndexOfPageObj = index
	
	//reset
	me.indexOfContent = -1 
	me.resetCurrXY()
}

//เริ่ม
func (me *GoPdf) Start(config Config) {

	me.config = config
	me.init()
	//สร้าง obj พื้นฐาน
	catalog := new(CatalogObj)
	catalog.Init(func()(*GoPdf){
		return me;
	});
	pages := new(PagesObj)
	pages.Init(func()(*GoPdf){
		return me
	})
	me.addObj(catalog)
	me.indexOfPagesObj = me.addObj(pages)
	
	//indexOfProcSet
	procset := new(ProcSetObj)
	procset.Init(func()(*GoPdf){
		return me;
	});
	me.indexOfProcSet = me.addObj(procset)
}



//set font 
func (me *GoPdf) SetFont(family string, style string, size int){

	i := 0
	max := len(me.indexEncodingObjFonts)
	for i < max {
		ifont := me.pdfObjs[me.indexEncodingObjFonts[i]].(*EncodingObj).GetFont()
		if ifont.GetFamily() == family {
			me.Curr.Font_Size   = size
			me.Curr.Font_Style = style
			me.Curr.Font_IFont = ifont
			me.Curr.Font_FontCount =   me.pdfObjs[me.indexEncodingObjFonts[i] + 4].(*FontObj).CountOfFont
			break
		}
		i++
	}


}

//create pdf file
func (me *GoPdf) WritePdf(pdfPath string) {
	ioutil.WriteFile(pdfPath, me.GetBytesPdf(), 0644)
}


//get bytes of pdf file
func (me *GoPdf) GetBytesPdf() []byte {
	me.prepare()
	buff := new(bytes.Buffer)
	i := 0
	max := len(me.pdfObjs)
	buff.WriteString("%PDF-1.7\n\n")
	linelens := make([]int, max)
	for i < max {
		linelens[i] = buff.Len()
		pdfObj := me.pdfObjs[i]
		pdfObj.Build()
		buff.WriteString(strconv.Itoa(i+1) + " 0 obj\n")
		buffbyte := pdfObj.GetObjBuff().Bytes()
		buff.Write(buffbyte)
		buff.WriteString("endobj\n\n")
		i++
	}
	me.xref(linelens, buff, &i)
	return buff.Bytes()
}

//สร้าง cell ของ text
//Note that this has no effect on Rect.H pdf (now). Fix later :-)
func (me *GoPdf) Cell(rectangle *Rect, text string) {

	//undelineOffset := ContentObj_CalTextHeight(me.Curr.Font_Size) + 1
	startX := me.Curr.X
	startY := me.Curr.Y
	me.getContent().AppendStream(rectangle,text)
	endX := me.Curr.X
	endY := me.Curr.Y

	//underline
	if strings.Contains(strings.ToUpper(me.Curr.Font_Style),"U") {
		//me.Line(x1,y1+undelineOffset,x2,y2+undelineOffset)
		me.getContent().AppendUnderline(startX,startY,endX,endY,text)
	}

}

func (me *GoPdf) AddFont(family string  ,ifont IFont, zfontpath string){
	encoding := new(EncodingObj)
	ifont.Init()
	ifont.SetFamily(family)
	encoding.SetFont(ifont)
	me.indexEncodingObjFonts = append(me.indexEncodingObjFonts, me.addObj(encoding))
	
	fontWidth := new(BasicObj)
	fontWidth.Init(func()(*GoPdf){
		return me
	})
	fontWidth.Data = "["+ FontConvertHelper_Cw2Str(ifont.GetCw())+"]\n"
	me.addObj(fontWidth)  //1
	
	fontDesc := new(FontDescriptorObj)
	fontDesc.Init(func()(*GoPdf){
		return me
	})
	fontDesc.SetFont(ifont)
	me.addObj(fontDesc) //2
	
	embedfont := new(EmbedFontObj)
	embedfont.Init(func()(*GoPdf){
		return me
	})
	embedfont.SetFont(ifont,zfontpath)	
	index := me.addObj(embedfont) //3
	
	fontDesc.SetFontFileObjRelate( strconv.Itoa(index + 1)  + " 0 R")
	
	
	//start add font obj
	font := new(FontObj)
	font.Init(func()(*GoPdf){
		return me
	})
	font.Family = family
	font.Font = ifont
	index  = me.addObj(font) //4
	if me.indexOfProcSet != -1 {
	 	procset := me.pdfObjs[me.indexOfProcSet].(*ProcSetObj)
	 	if !procset.Realtes.IsContainsFamily(family) {
	 		procset.Realtes = append(procset.Realtes,RelateFont{ Family : family, IndexOfObj : index , CountOfFont : me.Curr.CountOfFont  })
			font.CountOfFont = me.Curr.CountOfFont 
			me.Curr.CountOfFont++
		}
	}
	//end add font obj
}

/*---private---*/

//init
func (me *GoPdf) init() {

	//default
	me.leftMargin = 10.0
	me.topMargin = 10.0

	//init curr
	me.resetCurrXY()
	me.Curr.IndexOfPageObj = -1
	me.Curr.CountOfFont = 0
	me.Curr.CountOfL = 0
	me.Curr.CountOfImg = 0 //img
	me.Curr.ImgCaches = *new( []ImageCache)

	//init index
	me.indexOfPagesObj = -1
	me.indexOfFirstPageObj = -1
	me.indexOfContent = -1
	
	//No underline
	//me.IsUnderline = false

}

func (me * GoPdf) resetCurrXY(){
	me.Curr.X = me.leftMargin
	me.Curr.Y = me.topMargin
}

func (me *GoPdf) prepare() {
	
	if me.indexOfPagesObj != -1 {
		indexCurrPage := -1
		var pagesObj *PagesObj
		pagesObj = me.pdfObjs[me.indexOfPagesObj].(*PagesObj)
		i := 0//me.indexOfFirstPageObj
		max := len(me.pdfObjs)
		for i < max {
			objtype := me.pdfObjs[i].GetType()
			//fmt.Printf(" objtype = %s , %d \n", objtype , i)
			if objtype == "Page" {
				pagesObj.Kids = fmt.Sprintf("%s %d 0 R ", pagesObj.Kids, i+1)
				pagesObj.PageCount++
				indexCurrPage = i
			}else if  objtype == "Content" {
				if indexCurrPage != -1 {
					me.pdfObjs[indexCurrPage].(*PageObj).Contents = fmt.Sprintf("%s %d 0 R ",me.pdfObjs[indexCurrPage].(*PageObj).Contents,i+1);
				}
			}else if  objtype == "Font" {
				tmpfont := me.pdfObjs[i].(*FontObj)
				j := 0
				jmax := len(me.indexEncodingObjFonts)
				for j < jmax {
					tmpencoding := me.pdfObjs[me.indexEncodingObjFonts[j]].(*EncodingObj).GetFont()
					if tmpfont.Family == tmpencoding.GetFamily() { //ใส่ ข้อมูลของ embed font
						tmpfont.IsEmbedFont = true
						tmpfont.SetIndexObjEncoding( me.indexEncodingObjFonts[j] + 1)
						tmpfont.SetIndexObjWidth( me.indexEncodingObjFonts[j] + 2)
						tmpfont.SetIndexObjFontDescriptor( me.indexEncodingObjFonts[j] + 3)
						break
					}
					j++
				}
			}
			i++
		}
	}
}

func (me *GoPdf) xref(linelens []int, buff *bytes.Buffer, i *int) {
	buff.WriteString("xref\n")
	buff.WriteString("0 "+strconv.Itoa((*i)+1)+"\n")
	buff.WriteString("0000000000 65535 f\n")
	j := 0
	max := len(linelens)
	for j < max {
		linelen := linelens[j]
		buff.WriteString(me.formatXrefline(linelen) + " 00000 n\n")
		j++
	}
	buff.WriteString("trailer\n")
	buff.WriteString("<<\n")
	buff.WriteString("/Size " + strconv.Itoa(max+1) + "\n")
	buff.WriteString("/Root 1 0 R\n")
	buff.WriteString(">>\n")
	(*i)++
}

//ปรับ xref ให้เป็น 10 หลัก
func (me *GoPdf) formatXrefline(n int) string{
	str := strconv.Itoa(n)
	for len(str) < 10 {
		str = "0" + str
	}
	return str
}

func (me *GoPdf) addObj(iobj IObj) int {
	index := len(me.pdfObjs)
	me.pdfObjs = append(me.pdfObjs, iobj)
	return index
}

func (me * GoPdf) getContent() *ContentObj{
	var content *ContentObj
	if me.indexOfContent <= -1{
		content = new(ContentObj)
		content.Init(func()(*GoPdf){
			return me
		})
		me.indexOfContent = me.addObj(content)
	}else{
		content = me.pdfObjs[me.indexOfContent].(*ContentObj)
	}

	return content
}


