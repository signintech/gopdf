package gopdf

import (
	"bytes"
	obj "gopdf/obj"
	ioutil "io/ioutil"
	//"container/list"
	"fmt"
	"strconv"
)

type GoPdf struct {
	pdfObjs []obj.IObj

	/*index ของ obj สำคัญๆ*/
	//index ของ obj pages
	indexOfPagesObj int
	//index ของ obj page อันแรก
	indexOfFirstPageObj int
	//index ของ obj page อันปัจจุบันในขณะถูกเรียกจาก AddPage
	//indexOfCurrPageObj int
}

/*public:*/

//เพิ่ม page
func (me *GoPdf) AddPage() {
	page := new(obj.PageObj)
	page.Init()
	index := me.AddObj(page)
	if me.indexOfFirstPageObj == -1 {
		me.indexOfFirstPageObj = index
	}
}

func (me *GoPdf) Start() {
	me.init()
	catalog := new(obj.CatalogObj)
	catalog.Init()
	pages := new(obj.PagesObj)
	pages.Init()
	me.AddObj(catalog)
	me.indexOfPagesObj = me.AddObj(pages)
}

func (me *GoPdf) SetFont(family string, style string, size int){
	font := new(obj.FontObj)
	font.Init()
	me.AddObj(font)
}

func (me *GoPdf) WritePdf(pdfPath string) {
	me.prepare()
	buff := new(bytes.Buffer)
	i := 0
	max := len(me.pdfObjs)
	buff.WriteString("%PDF-1.7\n")
	linelens := make([]int, max)
	for i < max {
		pdfObj := me.pdfObjs[i]
		pdfObj.Build()
		buff.WriteString(strconv.Itoa(i+1) + " 0 obj\n<<\n")
		buffbyte := pdfObj.GetObjBuff().Bytes()
		linelens[i] = len(buffbyte)
		buff.Write(buffbyte)
		buff.WriteString(">>\nendobj\n\n")
		i++
	}
	me.xref(linelens, buff, &i)
	fmt.Printf("%s\n", buff.String())
	ioutil.WriteFile(pdfPath, buff.Bytes(), 0644)
}

func (me *GoPdf) Cell(pos Rect, txt string) {
	content := new(obj.ContentObj)
	content.Init()
	me.AddObj(content)
}

/*private:*/

//init
func (me *GoPdf) init() {
	me.indexOfPagesObj = -1
	me.indexOfFirstPageObj = -1
}

func (me *GoPdf) prepare() {
	
	if me.indexOfPagesObj != -1 {
		indexCurrPage := -1
		var pagesObj *obj.PagesObj
		pagesObj = me.pdfObjs[me.indexOfPagesObj].(*obj.PagesObj)
		i := me.indexOfFirstPageObj
		max := len(me.pdfObjs)
		for i < max {
			objtype := me.pdfObjs[i].GetType()
			if objtype == "Page" {
				pagesObj.Kids = fmt.Sprintf("%s %d 0 R ", pagesObj.Kids, i+1)
				pagesObj.PageCount++
				indexCurrPage = i
			}else if( objtype == "Content" ){
				if indexCurrPage != -1 {
					me.pdfObjs[indexCurrPage].(*obj.PageObj).Contents = fmt.Sprintf("%s %d 0 R ",me.pdfObjs[indexCurrPage].(*obj.PageObj).Contents,i+1);
				}
			}
			i++
		}
	}
}

func (me *GoPdf) xref(linelens []int, buff *bytes.Buffer, i *int) {
	buff.WriteString("xref\n")
	buff.WriteString(strconv.Itoa((*i)+1) + " 0\n")
	buff.WriteString("0000000000 65535 f\n")
	j := 0
	max := len(linelens)
	for j < max {
		linelen := linelens[j]
		buff.WriteString(strconv.Itoa(linelen) + " 00000 n\n")
		j++
	}
	buff.WriteString("trailer\n")
	buff.WriteString("<<\n")
	buff.WriteString("/Size " + strconv.Itoa(max+1) + "\n")
	buff.WriteString("/Root 1 0 R\n")
	buff.WriteString(">>\n")
	(*i)++
}

func (me *GoPdf) AddObj(iobj obj.IObj) int {
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
