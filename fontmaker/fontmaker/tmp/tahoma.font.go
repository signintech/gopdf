package fonts //change this
import (
	"github.com/signintech/gopdf"
)
type Tahoma struct {
	family string
	fonttype string
	name string
	desc  []gopdf.FontDescItem
	up int
	ut int
	cw gopdf.FontCw
	enc string
	diff string
}
func (me * Tahoma) Init(){
	me.cw = make(gopdf.FontCw)
	me.cw[gopdf.Chr(0)]=1000
	me.cw[gopdf.Chr(1)]=1000
	me.cw[gopdf.Chr(2)]=1000
	me.cw[gopdf.Chr(3)]=1000
	me.cw[gopdf.Chr(4)]=1000
	me.cw[gopdf.Chr(5)]=1000
	me.cw[gopdf.Chr(6)]=1000
	me.cw[gopdf.Chr(7)]=1000
	me.cw[gopdf.Chr(8)]=1000
	me.cw[gopdf.Chr(9)]=1000
	me.cw[gopdf.Chr(10)]=1000
	me.cw[gopdf.Chr(11)]=1000
	me.cw[gopdf.Chr(12)]=1000
	me.cw[gopdf.Chr(13)]=1000
	me.cw[gopdf.Chr(14)]=1000
	me.cw[gopdf.Chr(15)]=1000
	me.cw[gopdf.Chr(16)]=1000
	me.cw[gopdf.Chr(17)]=1000
	me.cw[gopdf.Chr(18)]=1000
	me.cw[gopdf.Chr(19)]=1000
	me.cw[gopdf.Chr(20)]=1000
	me.cw[gopdf.Chr(21)]=1000
	me.cw[gopdf.Chr(22)]=1000
	me.cw[gopdf.Chr(23)]=1000
	me.cw[gopdf.Chr(24)]=1000
	me.cw[gopdf.Chr(25)]=1000
	me.cw[gopdf.Chr(26)]=1000
	me.cw[gopdf.Chr(27)]=1000
	me.cw[gopdf.Chr(28)]=1000
	me.cw[gopdf.Chr(29)]=1000
	me.cw[gopdf.Chr(30)]=1000
	me.cw[gopdf.Chr(31)]=1000
	me.cw[gopdf.ToByte(" ")]=313
	me.cw[gopdf.ToByte("!")]=332
	me.cw[gopdf.ToByte("\"")]=401
	me.cw[gopdf.ToByte("#")]=728
	me.cw[gopdf.ToByte("$")]=546
	me.cw[gopdf.ToByte("%")]=977
	me.cw[gopdf.ToByte("&")]=674
	me.cw[gopdf.ToByte("'")]=211
	me.cw[gopdf.ToByte("(")]=383
	me.cw[gopdf.ToByte(")")]=383
	me.cw[gopdf.ToByte("*")]=546
	me.cw[gopdf.ToByte("+")]=728
	me.cw[gopdf.ToByte(",")]=303
	me.cw[gopdf.ToByte("-")]=363
	me.cw[gopdf.ToByte(".")]=303
	me.cw[gopdf.ToByte("/")]=382
	me.cw[gopdf.ToByte("0")]=546
	me.cw[gopdf.ToByte("1")]=546
	me.cw[gopdf.ToByte("2")]=546
	me.cw[gopdf.ToByte("3")]=546
	me.cw[gopdf.ToByte("4")]=546
	me.cw[gopdf.ToByte("5")]=546
	me.cw[gopdf.ToByte("6")]=546
	me.cw[gopdf.ToByte("7")]=546
	me.cw[gopdf.ToByte("8")]=546
	me.cw[gopdf.ToByte("9")]=546
	me.cw[gopdf.ToByte(":")]=354
	me.cw[gopdf.ToByte(";")]=354
	me.cw[gopdf.ToByte("<")]=728
	me.cw[gopdf.ToByte("=")]=728
	me.cw[gopdf.ToByte(">")]=728
	me.cw[gopdf.ToByte("?")]=474
	me.cw[gopdf.ToByte("@")]=909
	me.cw[gopdf.ToByte("A")]=600
	me.cw[gopdf.ToByte("B")]=589
	me.cw[gopdf.ToByte("C")]=601
	me.cw[gopdf.ToByte("D")]=678
	me.cw[gopdf.ToByte("E")]=561
	me.cw[gopdf.ToByte("F")]=521
	me.cw[gopdf.ToByte("G")]=667
	me.cw[gopdf.ToByte("H")]=675
	me.cw[gopdf.ToByte("I")]=373
	me.cw[gopdf.ToByte("J")]=417
	me.cw[gopdf.ToByte("K")]=588
	me.cw[gopdf.ToByte("L")]=498
	me.cw[gopdf.ToByte("M")]=771
	me.cw[gopdf.ToByte("N")]=667
	me.cw[gopdf.ToByte("O")]=708
	me.cw[gopdf.ToByte("P")]=551
	me.cw[gopdf.ToByte("Q")]=708
	me.cw[gopdf.ToByte("R")]=621
	me.cw[gopdf.ToByte("S")]=557
	me.cw[gopdf.ToByte("T")]=584
	me.cw[gopdf.ToByte("U")]=656
	me.cw[gopdf.ToByte("V")]=597
	me.cw[gopdf.ToByte("W")]=902
	me.cw[gopdf.ToByte("X")]=581
	me.cw[gopdf.ToByte("Y")]=576
	me.cw[gopdf.ToByte("Z")]=559
	me.cw[gopdf.ToByte("[")]=383
	me.cw[gopdf.ToByte("\\")]=382
	me.cw[gopdf.ToByte("]")]=383
	me.cw[gopdf.ToByte("^")]=728
	me.cw[gopdf.ToByte("_")]=546
	me.cw[gopdf.ToByte("`")]=546
	me.cw[gopdf.ToByte("a")]=525
	me.cw[gopdf.ToByte("b")]=553
	me.cw[gopdf.ToByte("c")]=461
	me.cw[gopdf.ToByte("d")]=553
	me.cw[gopdf.ToByte("e")]=526
	me.cw[gopdf.ToByte("f")]=318
	me.cw[gopdf.ToByte("g")]=553
	me.cw[gopdf.ToByte("h")]=558
	me.cw[gopdf.ToByte("i")]=229
	me.cw[gopdf.ToByte("j")]=282
	me.cw[gopdf.ToByte("k")]=498
	me.cw[gopdf.ToByte("l")]=229
	me.cw[gopdf.ToByte("m")]=840
	me.cw[gopdf.ToByte("n")]=558
	me.cw[gopdf.ToByte("o")]=543
	me.cw[gopdf.ToByte("p")]=553
	me.cw[gopdf.ToByte("q")]=553
	me.cw[gopdf.ToByte("r")]=360
	me.cw[gopdf.ToByte("s")]=446
	me.cw[gopdf.ToByte("t")]=334
	me.cw[gopdf.ToByte("u")]=558
	me.cw[gopdf.ToByte("v")]=498
	me.cw[gopdf.ToByte("w")]=742
	me.cw[gopdf.ToByte("x")]=495
	me.cw[gopdf.ToByte("y")]=498
	me.cw[gopdf.ToByte("z")]=444
	me.cw[gopdf.ToByte("{")]=480
	me.cw[gopdf.ToByte("|")]=382
	me.cw[gopdf.ToByte("}")]=480
	me.cw[gopdf.ToByte("~")]=728
	me.cw[gopdf.Chr(127)]=1000
	me.cw[gopdf.Chr(128)]=546
	me.cw[gopdf.Chr(129)]=1000
	me.cw[gopdf.Chr(130)]=1000
	me.cw[gopdf.Chr(131)]=1000
	me.cw[gopdf.Chr(132)]=1000
	me.cw[gopdf.Chr(133)]=817
	me.cw[gopdf.Chr(134)]=1000
	me.cw[gopdf.Chr(135)]=1000
	me.cw[gopdf.Chr(136)]=1000
	me.cw[gopdf.Chr(137)]=1000
	me.cw[gopdf.Chr(138)]=1000
	me.cw[gopdf.Chr(139)]=1000
	me.cw[gopdf.Chr(140)]=1000
	me.cw[gopdf.Chr(141)]=1000
	me.cw[gopdf.Chr(142)]=1000
	me.cw[gopdf.Chr(143)]=1000
	me.cw[gopdf.Chr(144)]=1000
	me.cw[gopdf.Chr(145)]=211
	me.cw[gopdf.Chr(146)]=211
	me.cw[gopdf.Chr(147)]=401
	me.cw[gopdf.Chr(148)]=401
	me.cw[gopdf.Chr(149)]=455
	me.cw[gopdf.Chr(150)]=546
	me.cw[gopdf.Chr(151)]=909
	me.cw[gopdf.Chr(152)]=1000
	me.cw[gopdf.Chr(153)]=1000
	me.cw[gopdf.Chr(154)]=1000
	me.cw[gopdf.Chr(155)]=1000
	me.cw[gopdf.Chr(156)]=1000
	me.cw[gopdf.Chr(157)]=1000
	me.cw[gopdf.Chr(158)]=1000
	me.cw[gopdf.Chr(159)]=1000
	me.cw[gopdf.Chr(160)]=313
	me.cw[gopdf.Chr(161)]=595
	me.cw[gopdf.Chr(162)]=613
	me.cw[gopdf.Chr(163)]=655
	me.cw[gopdf.Chr(164)]=615
	me.cw[gopdf.Chr(165)]=615
	me.cw[gopdf.Chr(166)]=673
	me.cw[gopdf.Chr(167)]=492
	me.cw[gopdf.Chr(168)]=556
	me.cw[gopdf.Chr(169)]=620
	me.cw[gopdf.Chr(170)]=627
	me.cw[gopdf.Chr(171)]=688
	me.cw[gopdf.Chr(172)]=862
	me.cw[gopdf.Chr(173)]=832
	me.cw[gopdf.Chr(174)]=621
	me.cw[gopdf.Chr(175)]=621
	me.cw[gopdf.Chr(176)]=556
	me.cw[gopdf.Chr(177)]=809
	me.cw[gopdf.Chr(178)]=872
	me.cw[gopdf.Chr(179)]=853
	me.cw[gopdf.Chr(180)]=615
	me.cw[gopdf.Chr(181)]=615
	me.cw[gopdf.Chr(182)]=585
	me.cw[gopdf.Chr(183)]=716
	me.cw[gopdf.Chr(184)]=540
	me.cw[gopdf.Chr(185)]=612
	me.cw[gopdf.Chr(186)]=630
	me.cw[gopdf.Chr(187)]=630
	me.cw[gopdf.Chr(188)]=632
	me.cw[gopdf.Chr(189)]=632
	me.cw[gopdf.Chr(190)]=708
	me.cw[gopdf.Chr(191)]=698
	me.cw[gopdf.Chr(192)]=621
	me.cw[gopdf.Chr(193)]=599
	me.cw[gopdf.Chr(194)]=579
	me.cw[gopdf.Chr(195)]=429
	me.cw[gopdf.Chr(196)]=567
	me.cw[gopdf.Chr(197)]=616
	me.cw[gopdf.Chr(198)]=621
	me.cw[gopdf.Chr(199)]=485
	me.cw[gopdf.Chr(200)]=615
	me.cw[gopdf.Chr(201)]=744
	me.cw[gopdf.Chr(202)]=608
	me.cw[gopdf.Chr(203)]=670
	me.cw[gopdf.Chr(204)]=715
	me.cw[gopdf.Chr(205)]=595
	me.cw[gopdf.Chr(206)]=595
	me.cw[gopdf.Chr(207)]=529
	me.cw[gopdf.Chr(208)]=448
	me.cw[gopdf.Chr(209)]=0
	me.cw[gopdf.Chr(210)]=485
	me.cw[gopdf.Chr(211)]=485
	me.cw[gopdf.Chr(212)]=0
	me.cw[gopdf.Chr(213)]=0
	me.cw[gopdf.Chr(214)]=0
	me.cw[gopdf.Chr(215)]=0
	me.cw[gopdf.Chr(216)]=0
	me.cw[gopdf.Chr(217)]=0
	me.cw[gopdf.Chr(218)]=0
	me.cw[gopdf.Chr(219)]=1000
	me.cw[gopdf.Chr(220)]=1000
	me.cw[gopdf.Chr(221)]=1000
	me.cw[gopdf.Chr(222)]=1000
	me.cw[gopdf.Chr(223)]=589
	me.cw[gopdf.Chr(224)]=332
	me.cw[gopdf.Chr(225)]=592
	me.cw[gopdf.Chr(226)]=492
	me.cw[gopdf.Chr(227)]=442
	me.cw[gopdf.Chr(228)]=511
	me.cw[gopdf.Chr(229)]=485
	me.cw[gopdf.Chr(230)]=499
	me.cw[gopdf.Chr(231)]=0
	me.cw[gopdf.Chr(232)]=0
	me.cw[gopdf.Chr(233)]=0
	me.cw[gopdf.Chr(234)]=0
	me.cw[gopdf.Chr(235)]=0
	me.cw[gopdf.Chr(236)]=0
	me.cw[gopdf.Chr(237)]=0
	me.cw[gopdf.Chr(238)]=0
	me.cw[gopdf.Chr(239)]=602
	me.cw[gopdf.Chr(240)]=602
	me.cw[gopdf.Chr(241)]=632
	me.cw[gopdf.Chr(242)]=685
	me.cw[gopdf.Chr(243)]=688
	me.cw[gopdf.Chr(244)]=733
	me.cw[gopdf.Chr(245)]=733
	me.cw[gopdf.Chr(246)]=646
	me.cw[gopdf.Chr(247)]=774
	me.cw[gopdf.Chr(248)]=774
	me.cw[gopdf.Chr(249)]=715
	me.cw[gopdf.Chr(250)]=694
	me.cw[gopdf.Chr(251)]=1065
	me.cw[gopdf.Chr(252)]=1000
	me.cw[gopdf.Chr(253)]=1000
	me.cw[gopdf.Chr(254)]=1000
	me.cw[gopdf.Chr(255)]=1000
	me.up = -83
	me.ut = 63
	me.fonttype = "TrueType"
	me.name = "Tahoma"
	me.enc = "cp874"
	me.diff = "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef /.notdef"
	me.desc = make([]gopdf.FontDescItem,8)
	me.desc[0] =  gopdf.FontDescItem{ Key:"Ascent",Val : "765" }
	me.desc[1] =  gopdf.FontDescItem{ Key: "Descent", Val : "-207" }
	me.desc[2] =  gopdf.FontDescItem{ Key:"CapHeight", Val :  "727" }
	me.desc[3] =  gopdf.FontDescItem{ Key: "Flags", Val :  "32" }
	me.desc[4] =  gopdf.FontDescItem{ Key:"FontBBox", Val :  "[-600 -216 1516 1034]" }
	me.desc[5] =  gopdf.FontDescItem{ Key:"ItalicAngle", Val :  "0" }
	me.desc[6] =  gopdf.FontDescItem{ Key:"StemV", Val :  "70" }
 	me.desc[7] =  gopdf.FontDescItem{ Key:"MissingWidth", Val :  "1000" } 
 }
func (me * Tahoma)GetType() string{
	return me.fonttype
}
func (me * Tahoma)GetName() string{
	return me.name
}	
func (me * Tahoma)GetDesc() []gopdf.FontDescItem{
	return me.desc
}
func (me * Tahoma)GetUp() int{
	return me.up
}
func (me * Tahoma)GetUt()  int{
	return me.ut
}
func (me * Tahoma)GetCw() gopdf.FontCw{
	return me.cw
}
func (me * Tahoma)GetEnc() string{
	return me.enc
}
func (me * Tahoma)GetDiff() string {
	return me.diff
}
func (me * Tahoma) GetOriginalsize() int{
	return 98764
}
func (me * Tahoma)  SetFamily(family string){
	me.family = family
}
func (me * Tahoma) 	GetFamily() string{
	return me.family
}
