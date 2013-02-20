package gopdf

import (
	"bytes"
	ioutil "io/ioutil"
	//"container/list"
	"fmt"
	"strconv"
	"gopdf/fonts"
)

type GoPdf struct {
	
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
}

/*---public---*/

//เพิ่ม page
func (me *GoPdf) AddPage() {
	page := new(PageObj)
	page.Init(func()(*GoPdf){
		return me
	})
	index := me.addObj(page)
	if me.indexOfFirstPageObj == -1 {
		me.indexOfFirstPageObj = index
	}
	me.Curr.IndexOfPageObj = index
	
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
}



//set font 
func (me *GoPdf) SetFont(family string, style string, size int){
	//ต้องแน่ใจว่ามีการ add font แล้ว
	me.Curr.FontSize = size
	font := new(FontObj)
	font.Init(func()(*GoPdf){
		return me
	})
	font.Family = family
	index := me.addObj(font)
	
	i := 0 
	max := len(me.indexEncodingObjFonts)
	for i < max {
		
		i++
	}
	
	
	if me.Curr.IndexOfPageObj != -1 {
	 	pageobj := me.pdfObjs[me.Curr.IndexOfPageObj].(*PageObj)
	 	pageobj.realtes = append(pageobj.realtes,"/F1 "+ strconv.Itoa(index+1) + " 0 R ")
	}
}

//สร้าง pdf to file
func (me *GoPdf) WritePdf(pdfPath string) {
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
	fmt.Printf("%s\n", buff.String())
	ioutil.WriteFile(pdfPath, buff.Bytes(), 0644)
}

//สร้าง cell ของ text
func (me *GoPdf) Cell(pos Rect, text string) {
	var content *ContentObj
	if me.indexOfContent == -1{ 
		content = new(ContentObj)
		content.Init(func()(*GoPdf){
			return me
		})	
		me.indexOfContent = me.addObj(content)
	}else{
		content = me.pdfObjs[me.indexOfContent].(*ContentObj)
	}
	content.AppendText(text)
	
}

func (me *GoPdf) AddFont(family string  ,ifont fonts.IFont, zfontpath string){
	encoding := new(EncodingObj)
	ifont.Init()
	ifont.SetFamily(family)
	encoding.SetFont(ifont)
	me.indexEncodingObjFonts = append(me.indexEncodingObjFonts, me.addObj(encoding))
	
	fontWidth := new(BasicObj)
	fontWidth.Init(func()(*GoPdf){
		return me
	})
	fontWidth.Data = "["+ fonts.FontConvertHelper_Cw2Str(ifont.GetCw())+"]\n"
	me.addObj(fontWidth)
	
	fontDesc := new(FontDescriptorObj)
	fontDesc.Init(func()(*GoPdf){
		return me
	})
	fontDesc.SetFont(ifont)
	me.addObj(fontDesc)
	
	embedfont := new(EmbedFontObj)
	embedfont.Init(func()(*GoPdf){
		return me
	})
	embedfont.SetFont(ifont,zfontpath)	
	index := me.addObj(embedfont)
	
	fontDesc.SetFontFileObjRelate( strconv.Itoa(index + 1)  + " 0 R")
}

/*---private---*/

//init
func (me *GoPdf) init() {
	me.Curr.X = 10.0
	me.Curr.Y = 10.0
	me.indexOfPagesObj = -1
	me.indexOfFirstPageObj = -1
	me.Curr.IndexOfPageObj = -1
	me.indexOfContent = -1
}

func (me *GoPdf) prepare() {
	
	if me.indexOfPagesObj != -1 {
		indexCurrPage := -1
		var pagesObj *PagesObj
		pagesObj = me.pdfObjs[me.indexOfPagesObj].(*PagesObj)
		i := me.indexOfFirstPageObj
		max := len(me.pdfObjs)
		for i < max {
			objtype := me.pdfObjs[i].GetType()
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
					if tmpfont.Family == tmpencoding.GetFamily() {
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



/*
//append text to buffer
func ( me *GoPdf ) appendBufferln(str string){
	me.buffer.WriteString(str);
	me.buffer.WriteString("\n");
}

//แปลง buffer เป็น string
func ( me *GoPdf) bufferToString() string{
	return me.buffer.String()
}
*/
