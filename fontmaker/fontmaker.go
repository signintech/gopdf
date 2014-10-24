package fontmaker

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/signintech/gopdf"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ERROR_FONT_LICENSE_DOES_NOT_ALLOW_EMBEDDING = errors.New("Font license does not allow embedding")

type FontMaker struct {
}

func NewFontMaker() *FontMaker {
	return new(FontMaker)
}

func (me *FontMaker) MakeFont(fontpath string, encodingpath string, outfolderpath string) (gopdf.IFont, error) {

	fmt.Println("start")

	//encode
	encode := filepath.Base(encodingpath)
	encode = strings.Replace(encode, ".map", "", -1)

	//read font file
	if _, err := os.Stat(fontpath); os.IsNotExist(err) {
		return nil, err
	}

	fileext := filepath.Ext(fontpath)
	if strings.ToLower(fileext) != ".ttf" {
		//now support only ttf :-P
		return nil, errors.New("support only ttf ")
	}

	fontmaps, err := me.LoadMap(encodingpath)
	if err != nil {
		return nil, err
	}

	info, err := me.GetInfoFromTrueType(fontpath, fontmaps)
	if err != nil {
		return nil, err
	}

	//zip
	basename := filepath.Base(fontpath)
	tmp := strings.Split(basename, ".")
	basename = strings.Replace(tmp[0], " ", "_", -1)
	gzfilename := basename + ".z"
	var buff bytes.Buffer
	gzipwriter := gzip.NewWriter(&buff)

	fontbytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return nil, err
	}

	_, err = gzipwriter.Write(fontbytes)
	if err != nil {
		return nil, err
	}
	gzipwriter.Close()

	err = ioutil.WriteFile(outfolderpath+"/"+gzfilename, buff.Bytes(), 0666)
	if err != nil {
		return nil, err
	}
	info.PushString("File", gzfilename)

	err = me.MakeDefinitionFile(me.GoStructName(basename), outfolderpath+"/"+basename+".font.go", encode, fontmaps, info)
	if err != nil {
		return nil, err
	}

	fmt.Println("end")
	return nil, nil
}

func (me *FontMaker) GoStructName(name string) string {
	goname := strings.ToUpper(name[0:1]) + name[1:]
	return goname
}

func (me *FontMaker) MakeDefinitionFile(gofontname string, exportfile string, encode string, fontmaps []FontMap, info TtfInfo) error {

	fonttype := "TrueType"
	str := ""
	str += "package fonts //change this\n"
	str += "import (\n"
	str += "	\"github.com/signintech/gopdf\"\n"
	str += ")\n"
	str += "type " + gofontname + " struct {\n"
	str += "\tfamily string\n"
	str += "\tfonttype string\n"
	str += "\tname string\n"
	str += "\tdesc  []gopdf.FontDescItem\n"
	str += "\tup int\n"
	str += "\tut int\n"
	str += "\tcw gopdf.FontCw\n"
	str += "\tenc string\n"
	str += "\tdiff string\n"
	str += "}\n"

	str += "func (me * " + gofontname + ") Init(){\n"
	widths, err := info.GetMapIntInt64("Widths")
	if err != nil {
		return err
	}
	tmp, err := me.MakeWidthArray(widths)
	if err != nil {
		return err
	}
	str += tmp
	//str += "\tme.up = " + info["UnderlinePosition"] + "\n"
	//str += "\tme.ut = " + info["UnderlineThickness"] + "\n"
	str += "\tme.fonttype = \"" + fonttype + "\"\n"
	//str += "\tme.name = \"" + info["FontName"] + "\"\n"
	str += "\tme.enc = \"" + encode + "\"\n"
	//diff := MakeFontEncoding(fontmaps)
	//if diff {
	//	str += "\tme.diff = \"" + diff + "\"\n"
	//}
	//str += MakeFontDescriptor(info)
	str += "}\n"

	str += "func (me * " + gofontname + ")GetType() string{\n"
	str += "\treturn me.fonttype\n"
	str += "}\n"
	str += "func (me * " + gofontname + ")GetName() string{\n"
	str += "\treturn me.name\n"
	str += "}	\n"
	str += "func (me * " + gofontname + ")GetDesc() []gopdf.FontDescItem{\n"
	str += "\treturn me.desc\n"
	str += "}\n"
	str += "func (me * " + gofontname + ")GetUp() int{\n"
	str += "\treturn me.up\n"
	str += "}\n"
	str += "func (me * " + gofontname + ")GetUt()  int{\n"
	str += "\treturn me.ut\n"
	str += "}\n"
	str += "func (me * " + gofontname + ")GetCw() gopdf.FontCw{\n"
	str += "\treturn me.cw\n"
	str += "}\n"
	str += "func (me * " + gofontname + ")GetEnc() string{\n"
	str += "\treturn me.enc\n"
	str += "}\n"
	str += "func (me * " + gofontname + ")GetDiff() string {\n"
	str += "\treturn me.diff\n"
	str += "}\n"

	str += "func (me * " + gofontname + ") GetOriginalsize() int{\n"
	str += "\treturn 98764\n"
	str += "}\n"

	str += "func (me * " + gofontname + ")  SetFamily(family string){\n"
	str += "\tme.family = family\n"
	str += "}\n"

	str += "func (me * " + gofontname + ") 	GetFamily() string{\n"
	str += "\treturn me.family\n"
	str += "}\n"

	fmt.Printf("%s\n", str)
	return nil
}

func (me *FontMaker) MakeWidthArray(widths map[int]int64) (string, error) {

	str := "\tme.cw = make(gopdf.FontCw)\n"
	for c := 0; c <= 255; c++ {
		str += "\tme.cw["
		chr := string(c)
		if chr == "\"" {
			str += "gopdf.ToByte(\"\\\"\")"
		} else if chr == "\\" {
			str += "gopdf.ToByte(\"\\\\\")"
		} else if c >= 32 && c <= 126 {
			str += "gopdf.ToByte(\"" + chr + "\")"
		} else {
			str += fmt.Sprintf("gopdf.Chr(%d)", c)
		}
		str += fmt.Sprintf("]=%d\n", widths[c])
	}
	return str, nil
}

func (me *FontMaker) FileSize(path string) (int64, error) {

	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}

func (me *FontMaker) GetInfoFromTrueType(fontpath string, fontmaps []FontMap) (TtfInfo, error) {

	var parser TTFParser
	err := parser.Parse(fontpath)
	if err != nil {
		return nil, err
	}

	if !parser.Embeddable {
		return nil, ERROR_FONT_LICENSE_DOES_NOT_ALLOW_EMBEDDING
	}

	info := NewTtfInfo()

	fileContent, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return nil, err
	}
	info.PushBytes("Data", fileContent)

	size, err := me.FileSize(fontpath)
	if err != nil {
		return nil, err
	}
	info.PushInt64("OriginalSize", size)

	k := float64(1000.0 / float64(parser.unitsPerEm))
	info.PushString("FontName", parser.postScriptName)
	info.PushBool("Bold", parser.Bold)
	info.PushInt64("ItalicAngle", parser.italicAngle)
	info.PushBool("IsFixedPitch", parser.isFixedPitch)
	info.PushInt64("Ascender", me.MultiplyAndRound(k, parser.typoAscender))
	info.PushInt64("Descender", me.MultiplyAndRound(k, parser.typoDescender))
	info.PushInt64("UnderlineThickness", me.MultiplyAndRound(k, parser.underlineThickness))

	fontBBoxs := []int64{
		me.MultiplyAndRoundWithUInt64(k, parser.xMin),
		me.MultiplyAndRoundWithUInt64(k, parser.yMin),
		me.MultiplyAndRoundWithUInt64(k, parser.xMax),
		me.MultiplyAndRoundWithUInt64(k, parser.yMax),
	}
	info.PushInt64s("FontBBox", fontBBoxs)
	info.PushInt64("CapHeight", me.MultiplyAndRound(k, parser.capHeight))
	missingWidth := me.MultiplyAndRoundWithUInt64(k, parser.widths[0])
	info.PushInt64("MissingWidth", missingWidth)

	widths := make(map[int]int64)
	max := 256
	c := 0
	for c < max {
		widths[c] = missingWidth
		c++
	}

	c = 0 //reset
	for c < max {
		if fontmaps[c].Name != ".notdef" {
			uv := fontmaps[c].Uv
			if val, ok := parser.chars[int(uv)]; ok {
				w := parser.widths[val]
				widths[c] = me.MultiplyAndRoundWithUInt64(k, w)
			} else {
				fmt.Printf("Character %s (%d) is missing\n", fontmaps[c].Name, fontmaps[c].Uv)
			}
		}
		c++
	}
	info.PushMapIntInt64("Widths", widths)

	return info, nil
}

func (me *FontMaker) MultiplyAndRoundWithUInt64(k float64, v uint64) int64 {
	r := k * float64(v)
	//fmt.Printf("%#v %#v %#v = %#v\n", k, v, r, me.Round(r))
	return me.Round(r)
}

func (me *FontMaker) MultiplyAndRound(k float64, v int64) int64 {
	r := k * float64(v)
	//fmt.Printf("%#v %#v %#v = %#v\n", k, v, r, me.Round(r))
	return me.Round(r)
}

func (me *FontMaker) Round(value float64) int64 {
	if value < 0.0 {
		value -= 0.5
	} else {
		value += 0.5
	}
	return int64(value)
}

func (me *FontMaker) LoadMap(encodingpath string) ([]FontMap, error) {

	if _, err := os.Stat(encodingpath); os.IsNotExist(err) {
		return nil, err
	}

	var fontmaps []FontMap
	i := 0
	max := 256
	for i < max {
		fontmaps = append(fontmaps, FontMap{Uv: -1, Name: ".notdef"})
		i++
	}

	file, err := os.Open(encodingpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		e := strings.Split(line, " ")
		strC := (e[0])[1:]
		strUv := (e[1])[2:]
		c, err := strconv.ParseInt(strC, 16, 0)
		if err != nil {
			return nil, err
		}
		uv, err := strconv.ParseInt(strUv, 16, 0)
		if err != nil {
			return nil, err
		}
		name := e[2]
		//fmt.Println("strC = "+strC+"strUv = "+strUv+" c=%d , uv= %d", c, uv)
		fontmaps[c].Name = name
		fontmaps[c].Uv = int(uv)
	}

	return fontmaps, nil
}

func (me *FontMaker) CompressFont(path string) (*bytes.Buffer, error) {
	rawbytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	gw := gzip.NewWriter(&buff)
	_, err = gw.Write(rawbytes)
	if err != nil {
		return nil, err
	}
	gw.Close()
	return &buff, nil
}
