package core

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//ErrFontLicenseDoesNotAllowEmbedding Font license does not allow embedding
var ErrFontLicenseDoesNotAllowEmbedding = errors.New("Font license does not allow embedding")

//FontMaker font maker
type FontMaker struct {
	results []string
}

//GetResults get result
func (f *FontMaker) GetResults() []string {
	return f.results
}

//NewFontMaker new FontMaker
func NewFontMaker() *FontMaker {
	return new(FontMaker)
}

func (f *FontMaker) MakeFont(fontpath string, mappath string, encode string, outfolderpath string) error {

	encodingpath := mappath + "/" + encode + ".map"

	//read font file
	if _, err := os.Stat(fontpath); os.IsNotExist(err) {
		return err
	}

	fileext := filepath.Ext(fontpath)
	if strings.ToLower(fileext) != ".ttf" {
		//now support only ttf :-P
		return errors.New("support only ttf ")
	}

	fontmaps, err := f.LoadMap(encodingpath)
	if err != nil {
		return err
	}

	info, err := f.GetInfoFromTrueType(fontpath, fontmaps)
	if err != nil {
		return err
	}

	//zip
	basename := filepath.Base(fontpath)
	tmp := strings.Split(basename, ".")
	basename = strings.Replace(tmp[0], " ", "_", -1)
	gzfilename := basename + ".z"

	var buff bytes.Buffer
	gzipwriter := zlib.NewWriter(&buff)

	fontbytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return err
	}

	_, err = gzipwriter.Write(fontbytes)
	if err != nil {
		return err
	}
	gzipwriter.Close()
	err = ioutil.WriteFile(outfolderpath+"/"+gzfilename, buff.Bytes(), 0644)
	if err != nil {
		return err
	}
	info.PushString("File", gzfilename)
	f.results = append(f.results, fmt.Sprintf("Save Z file at %s.", outfolderpath+"/"+gzfilename))

	//Definition File
	_, err = f.MakeDefinitionFile(f.GoStructName(basename), mappath, outfolderpath+"/"+basename+".font.go", encode, fontmaps, info)
	if err != nil {
		return err
	}

	return nil
}

func (f *FontMaker) GoStructName(name string) string {
	goname := strings.ToUpper(name[0:1]) + name[1:]
	return goname
}

func (f *FontMaker) MakeDefinitionFile(gofontname string, mappath string, exportfile string, encode string, fontmaps []FontMap, info TtfInfo) (string, error) {

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
		return "", err
	}

	tmpStr, err := f.MakeWidthArray(widths)
	if err != nil {
		return "", err
	}
	str += tmpStr

	tmpInt64, err := info.GetInt64("UnderlinePosition")
	if err != nil {
		return "", err
	}
	str += fmt.Sprintf("\tme.up = %d\n", tmpInt64)

	tmpInt64, err = info.GetInt64("UnderlineThickness")
	if err != nil {
		return "", err
	}
	str += fmt.Sprintf("\tme.ut = %d\n", tmpInt64)

	str += "\tme.fonttype = \"" + fonttype + "\"\n"

	tmpStr, err = info.GetString("FontName")
	if err != nil {
		return "", err
	}
	str += fmt.Sprintf("\tme.name = \"%s\"\n", tmpStr)

	str += "\tme.enc = \"" + encode + "\"\n"
	diff, err := f.MakeFontEncoding(mappath, fontmaps)
	if err != nil {
		return "", err
	}
	if diff != "" {
		str += "\tme.diff = \"" + diff + "\"\n"
	}

	fd, err := f.MakeFontDescriptor(info)
	if err != nil {
		return "", err
	}
	str += fd

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

	err = ioutil.WriteFile(exportfile, []byte(str), 0666)
	if err != nil {
		return "", err
	}
	f.results = append(f.results, fmt.Sprintf("Save GO file at %s.", exportfile))
	return str, nil
}

func (f *FontMaker) MakeFontDescriptor(info TtfInfo) (string, error) {

	fd := ""
	fd = "\tme.desc = make([]gopdf.FontDescItem,8)\n"

	// Ascent
	ascender, err := info.GetInt64("Ascender")
	if err != nil {
		return "", err
	}
	fd += fmt.Sprintf("\tme.desc[0] =  gopdf.FontDescItem{ Key:\"Ascent\",Val : \"%d\" }\n", ascender)

	// Descent
	descender, err := info.GetInt64("Descender")
	if err != nil {
		return "", err
	}
	fd += fmt.Sprintf("\tme.desc[1] =  gopdf.FontDescItem{ Key: \"Descent\", Val : \"%d\" }\n", descender)

	// CapHeight
	capHeight, err := info.GetInt64("CapHeight")
	if err == nil {
		fd += fmt.Sprintf("\tme.desc[2] =  gopdf.FontDescItem{ Key:\"CapHeight\", Val :  \"%d\" }\n", capHeight)
	} else if err == ERROR_NO_KEY_FOUND {
		fd += fmt.Sprintf("\tme.desc[2] =  gopdf.FontDescItem{ Key:\"CapHeight\", Val :  \"%d\" }\n", ascender)
	} else {
		return "", err
	}

	// Flags
	flags := 0
	isFixedPitch, err := info.GetBool("IsFixedPitch")
	if err != nil {
		return "", err
	}

	if isFixedPitch {
		flags += 1 << 0
	}
	flags += 1 << 5
	italicAngle, err := info.GetInt64("ItalicAngle")
	if italicAngle != 0 {
		flags += 1 << 6
	}
	fd += fmt.Sprintf("\tme.desc[3] =  gopdf.FontDescItem{ Key: \"Flags\", Val :  \"%d\" }\n", flags)
	//fmt.Printf("\n----\n")
	// FontBBox
	fbb, err := info.GetInt64s("FontBBox")
	if err != nil {
		return "", err
	}
	fd += fmt.Sprintf("\tme.desc[4] =  gopdf.FontDescItem{ Key:\"FontBBox\", Val :  \"[%d %d %d %d]\" }\n", fbb[0], fbb[1], fbb[2], fbb[3])

	// ItalicAngle
	fd += fmt.Sprintf("\tme.desc[5] =  gopdf.FontDescItem{ Key:\"ItalicAngle\", Val :  \"%d\" }\n", italicAngle)

	// StemV
	stdVW, err := info.GetInt64("StdVW")
	issetStdVW := false
	if err == nil {
		issetStdVW = true
	} else if err == ERROR_NO_KEY_FOUND {
		issetStdVW = false
	} else {
		return "", err
	}

	bold, err := info.GetBool("Bold")
	if err != nil {
		return "", err
	}

	stemv := int(0)
	if issetStdVW {
		stemv = stdVW
	} else if bold {
		stemv = 120
	} else {
		stemv = 70
	}
	fd += fmt.Sprintf("\tme.desc[6] =  gopdf.FontDescItem{ Key:\"StemV\", Val :  \"%d\" }\n ", stemv)

	// MissingWidth
	missingWidth, err := info.GetInt64("MissingWidth")
	if err != nil {
		return "", err
	}
	fd += fmt.Sprintf("\tme.desc[7] =  gopdf.FontDescItem{ Key:\"MissingWidth\", Val :  \"%d\" } \n ", missingWidth)
	return fd, nil
}

func (f *FontMaker) MakeFontEncoding(mappath string, fontmaps []FontMap) (string, error) {

	refpath := mappath + "/cp1252.map"
	ref, err := f.LoadMap(refpath)
	if err != nil {
		return "", err
	}
	s := ""
	last := 0
	for c := 0; c <= 255; c++ {
		if fontmaps[c].Name != ref[c].Name {
			if c != last+1 {
				s += fmt.Sprintf("%d ", c)
			}
			last = c
			s += "/" + fontmaps[c].Name + " "
		}
	}
	return strings.TrimSpace(s), nil
}

func (f *FontMaker) MakeWidthArray(widths map[int]int) (string, error) {

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

func (f *FontMaker) FileSize(path string) (int64, error) {

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

func (f *FontMaker) GetInfoFromTrueType(fontpath string, fontmaps []FontMap) (TtfInfo, error) {

	var parser TTFParser
	err := parser.Parse(fontpath)
	if err != nil {
		return nil, err
	}

	if !parser.Embeddable {
		return nil, ErrFontLicenseDoesNotAllowEmbedding
	}

	info := NewTtfInfo()

	fileContent, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return nil, err
	}
	info.PushBytes("Data", fileContent)

	size, err := f.FileSize(fontpath)
	if err != nil {
		return nil, err
	}
	info.PushInt64("OriginalSize", size)

	k := float64(1000.0 / float64(parser.unitsPerEm))
	info.PushString("FontName", parser.postScriptName)
	info.PushBool("Bold", parser.Bold)
	info.PushInt("ItalicAngle", parser.italicAngle)
	info.PushBool("IsFixedPitch", parser.isFixedPitch)
	info.PushInt("Ascender", f.MultiplyAndRound(k, parser.typoAscender))
	info.PushInt("Descender", f.MultiplyAndRound(k, parser.typoDescender))
	info.PushInt("UnderlineThickness", f.MultiplyAndRound(k, parser.underlineThickness))
	info.PushInt("UnderlinePosition", f.MultiplyAndRound(k, parser.underlinePosition))
	fontBBoxs := []int{
		f.MultiplyAndRound(k, parser.xMin),
		f.MultiplyAndRound(k, parser.yMin),
		f.MultiplyAndRound(k, parser.xMax),
		f.MultiplyAndRound(k, parser.yMax),
	}
	info.PushInt64s("FontBBox", fontBBoxs)
	info.PushInt("CapHeight", f.MultiplyAndRound(k, parser.capHeight))
	missingWidth := f.MultiplyAndRoundWithUInt64(k, parser.widths[0])
	info.PushInt("MissingWidth", missingWidth)

	widths := make(map[int]int)
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
				widths[c] = f.MultiplyAndRoundWithUInt64(k, w)
			} else {
				f.results = append(f.results, fmt.Sprintf("Warning: Character %s (%d) is missing", fontmaps[c].Name, fontmaps[c].Uv))
			}
		}
		c++
	}
	info.PushMapIntInt64("Widths", widths)

	return info, nil
}

func (f *FontMaker) MultiplyAndRoundWithUInt64(k float64, v uint) int {
	r := k * float64(v)
	return f.Round(r)
}

func (f *FontMaker) MultiplyAndRound(k float64, v int) int {
	r := k * float64(v)
	return f.Round(r)
}

func (f *FontMaker) Round(value float64) int {
	return Round(value)
}

func (f *FontMaker) LoadMap(encodingpath string) ([]FontMap, error) {

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
