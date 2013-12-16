package gopdf

import (

)

type THSarabunNewBold struct {
	family string
	fonttype string
	name string
	desc  []FontDescItem
	up int
	ut int
	cw FontCw
	enc string
	diff string
}
func (me * THSarabunNewBold) Init(){
	me.cw = make(FontCw)
	me.cw[Chr(0)]=698
	me.cw[Chr(1)]=698
	me.cw[Chr(2)]=698
	me.cw[Chr(3)]=698
	me.cw[Chr(4)]=698
	me.cw[Chr(5)]=698
	me.cw[Chr(6)]=698
	me.cw[Chr(7)]=698
	me.cw[Chr(8)]=698
	me.cw[Chr(9)]=698
	me.cw[Chr(10)]=698
	me.cw[Chr(11)]=698
	me.cw[Chr(12)]=698
	me.cw[Chr(13)]=698
	me.cw[Chr(14)]=698
	me.cw[Chr(15)]=698
	me.cw[Chr(16)]=698
	me.cw[Chr(17)]=698
	me.cw[Chr(18)]=698
	me.cw[Chr(19)]=698
	me.cw[Chr(20)]=698
	me.cw[Chr(21)]=698
	me.cw[Chr(22)]=698
	me.cw[Chr(23)]=698
	me.cw[Chr(24)]=698
	me.cw[Chr(25)]=698
	me.cw[Chr(26)]=698
	me.cw[Chr(27)]=698
	me.cw[Chr(28)]=698
	me.cw[Chr(29)]=698
	me.cw[Chr(30)]=698
	me.cw[Chr(31)]=698
	me.cw[ToByte(" ")]=226
	me.cw[ToByte("!")]=168
	me.cw[ToByte("\"")]=291
	me.cw[ToByte("#")]=462
	me.cw[ToByte("$")]=367
	me.cw[ToByte("%")]=688
	me.cw[ToByte("&")]=486
	me.cw[ToByte("'")]=169
	me.cw[ToByte("(")]=244
	me.cw[ToByte(")")]=244
	me.cw[ToByte("*")]=343
	me.cw[ToByte("+")]=417
	me.cw[ToByte(",")]=172
	me.cw[ToByte("-")]=226
	me.cw[ToByte(".")]=172
	me.cw[ToByte("/")]=276
	me.cw[ToByte("0")]=378
	me.cw[ToByte("1")]=378
	me.cw[ToByte("2")]=378
	me.cw[ToByte("3")]=378
	me.cw[ToByte("4")]=378
	me.cw[ToByte("5")]=378
	me.cw[ToByte("6")]=378
	me.cw[ToByte("7")]=378
	me.cw[ToByte("8")]=378
	me.cw[ToByte("9")]=378
	me.cw[ToByte(":")]=172
	me.cw[ToByte(";")]=172
	me.cw[ToByte("<")]=417
	me.cw[ToByte("=")]=417
	me.cw[ToByte(">")]=417
	me.cw[ToByte("?")]=289
	me.cw[ToByte("@")]=572
	me.cw[ToByte("A")]=451
	me.cw[ToByte("B")]=396
	me.cw[ToByte("C")]=432
	me.cw[ToByte("D")]=459
	me.cw[ToByte("E")]=358
	me.cw[ToByte("F")]=358
	me.cw[ToByte("G")]=454
	me.cw[ToByte("H")]=472
	me.cw[ToByte("I")]=163
	me.cw[ToByte("J")]=282
	me.cw[ToByte("K")]=411
	me.cw[ToByte("L")]=359
	me.cw[ToByte("M")]=574
	me.cw[ToByte("N")]=472
	me.cw[ToByte("O")]=518
	me.cw[ToByte("P")]=390
	me.cw[ToByte("Q")]=518
	me.cw[ToByte("R")]=395
	me.cw[ToByte("S")]=365
	me.cw[ToByte("T")]=421
	me.cw[ToByte("U")]=482
	me.cw[ToByte("V")]=422
	me.cw[ToByte("W")]=612
	me.cw[ToByte("X")]=426
	me.cw[ToByte("Y")]=396
	me.cw[ToByte("Z")]=451
	me.cw[ToByte("[")]=234
	me.cw[ToByte("\\")]=268
	me.cw[ToByte("]")]=234
	me.cw[ToByte("^")]=418
	me.cw[ToByte("_")]=376
	me.cw[ToByte("`")]=210
	me.cw[ToByte("a")]=367
	me.cw[ToByte("b")]=435
	me.cw[ToByte("c")]=347
	me.cw[ToByte("d")]=435
	me.cw[ToByte("e")]=391
	me.cw[ToByte("f")]=254
	me.cw[ToByte("g")]=332
	me.cw[ToByte("h")]=423
	me.cw[ToByte("i")]=163
	me.cw[ToByte("j")]=171
	me.cw[ToByte("k")]=340
	me.cw[ToByte("l")]=213
	me.cw[ToByte("m")]=636
	me.cw[ToByte("n")]=423
	me.cw[ToByte("o")]=412
	me.cw[ToByte("p")]=431
	me.cw[ToByte("q")]=431
	me.cw[ToByte("r")]=246
	me.cw[ToByte("s")]=296
	me.cw[ToByte("t")]=258
	me.cw[ToByte("u")]=418
	me.cw[ToByte("v")]=358
	me.cw[ToByte("w")]=480
	me.cw[ToByte("x")]=349
	me.cw[ToByte("y")]=361
	me.cw[ToByte("z")]=358
	me.cw[ToByte("{")]=255
	me.cw[ToByte("|")]=177
	me.cw[ToByte("}")]=255
	me.cw[ToByte("~")]=422
	me.cw[Chr(127)]=698
	me.cw[Chr(128)]=518
	me.cw[Chr(129)]=698
	me.cw[Chr(130)]=698
	me.cw[Chr(131)]=698
	me.cw[Chr(132)]=698
	me.cw[Chr(133)]=504
	me.cw[Chr(134)]=698
	me.cw[Chr(135)]=698
	me.cw[Chr(136)]=698
	me.cw[Chr(137)]=698
	me.cw[Chr(138)]=698
	me.cw[Chr(139)]=698
	me.cw[Chr(140)]=698
	me.cw[Chr(141)]=698
	me.cw[Chr(142)]=698
	me.cw[Chr(143)]=698
	me.cw[Chr(144)]=698
	me.cw[Chr(145)]=253
	me.cw[Chr(146)]=253
	me.cw[Chr(147)]=384
	me.cw[Chr(148)]=384
	me.cw[Chr(149)]=226
	me.cw[Chr(150)]=335
	me.cw[Chr(151)]=693
	me.cw[Chr(152)]=698
	me.cw[Chr(153)]=698
	me.cw[Chr(154)]=698
	me.cw[Chr(155)]=698
	me.cw[Chr(156)]=698
	me.cw[Chr(157)]=698
	me.cw[Chr(158)]=698
	me.cw[Chr(159)]=698
	me.cw[Chr(160)]=226
	me.cw[Chr(161)]=391
	me.cw[Chr(162)]=415
	me.cw[Chr(163)]=431
	me.cw[Chr(164)]=403
	me.cw[Chr(165)]=403
	me.cw[Chr(166)]=453
	me.cw[Chr(167)]=318
	me.cw[Chr(168)]=390
	me.cw[Chr(169)]=405
	me.cw[Chr(170)]=415
	me.cw[Chr(171)]=431
	me.cw[Chr(172)]=565
	me.cw[Chr(173)]=565
	me.cw[Chr(174)]=425
	me.cw[Chr(175)]=425
	me.cw[Chr(176)]=372
	me.cw[Chr(177)]=492
	me.cw[Chr(178)]=551
	me.cw[Chr(179)]=576
	me.cw[Chr(180)]=405
	me.cw[Chr(181)]=405
	me.cw[Chr(182)]=391
	me.cw[Chr(183)]=438
	me.cw[Chr(184)]=358
	me.cw[Chr(185)]=433
	me.cw[Chr(186)]=445
	me.cw[Chr(187)]=445
	me.cw[Chr(188)]=414
	me.cw[Chr(189)]=414
	me.cw[Chr(190)]=482
	me.cw[Chr(191)]=482
	me.cw[Chr(192)]=425
	me.cw[Chr(193)]=417
	me.cw[Chr(194)]=386
	me.cw[Chr(195)]=329
	me.cw[Chr(196)]=391
	me.cw[Chr(197)]=399
	me.cw[Chr(198)]=425
	me.cw[Chr(199)]=379
	me.cw[Chr(200)]=403
	me.cw[Chr(201)]=461
	me.cw[Chr(202)]=399
	me.cw[Chr(203)]=438
	me.cw[Chr(204)]=468
	me.cw[Chr(205)]=394
	me.cw[Chr(206)]=385
	me.cw[Chr(207)]=380
	me.cw[Chr(208)]=340
	me.cw[Chr(209)]=0
	me.cw[Chr(210)]=336
	me.cw[Chr(211)]=336
	me.cw[Chr(212)]=0
	me.cw[Chr(213)]=0
	me.cw[Chr(214)]=0
	me.cw[Chr(215)]=0
	me.cw[Chr(216)]=0
	me.cw[Chr(217)]=0
	me.cw[Chr(218)]=0
	me.cw[Chr(219)]=698
	me.cw[Chr(220)]=698
	me.cw[Chr(221)]=698
	me.cw[Chr(222)]=698
	me.cw[Chr(223)]=389
	me.cw[Chr(224)]=208
	me.cw[Chr(225)]=395
	me.cw[Chr(226)]=252
	me.cw[Chr(227)]=269
	me.cw[Chr(228)]=252
	me.cw[Chr(229)]=209
	me.cw[Chr(230)]=437
	me.cw[Chr(231)]=0
	me.cw[Chr(232)]=0
	me.cw[Chr(233)]=0
	me.cw[Chr(234)]=0
	me.cw[Chr(235)]=0
	me.cw[Chr(236)]=0
	me.cw[Chr(237)]=0
	me.cw[Chr(238)]=0
	me.cw[Chr(239)]=450
	me.cw[Chr(240)]=456
	me.cw[Chr(241)]=456
	me.cw[Chr(242)]=456
	me.cw[Chr(243)]=456
	me.cw[Chr(244)]=456
	me.cw[Chr(245)]=456
	me.cw[Chr(246)]=456
	me.cw[Chr(247)]=456
	me.cw[Chr(248)]=456
	me.cw[Chr(249)]=456
	me.cw[Chr(250)]=525
	me.cw[Chr(251)]=697
	me.cw[Chr(252)]=698
	me.cw[Chr(253)]=698
	me.cw[Chr(254)]=698
	me.cw[Chr(255)]=698
	me.up = -35
	me.ut = 30
	me.fonttype = "TrueType"
	me.name = "THSarabunNew-Bold"
	me.enc = "cp874"
	me.diff = "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef /.notdef"
	me.desc = make([]FontDescItem,8)
	me.desc[0] =  FontDescItem{ Key:"Ascent",Val : "850" }
	me.desc[1] =  FontDescItem{ Key: "Descent", Val : "-250" }
	me.desc[2] =  FontDescItem{ Key:"CapHeight", Val :  "476" }
	me.desc[3] =  FontDescItem{ Key: "Flags", Val :  "32" }
	me.desc[4] =  FontDescItem{ Key:"FontBBox", Val :  "[-466 -457 947 844]" }
	me.desc[5] =  FontDescItem{ Key:"ItalicAngle", Val :  "0" }
	me.desc[6] =  FontDescItem{ Key:"StemV", Val :  "120" }
 	me.desc[7] =  FontDescItem{ Key:"MissingWidth", Val :  "698" } 
 }
func (me * THSarabunNewBold)GetType() string{
	return me.fonttype
}
func (me * THSarabunNewBold)GetName() string{
	return me.name
}	
func (me * THSarabunNewBold)GetDesc() []FontDescItem{
	return me.desc
}
func (me * THSarabunNewBold)GetUp() int{
	return me.up
}
func (me * THSarabunNewBold)GetUt()  int{
	return me.ut
}
func (me * THSarabunNewBold)GetCw() FontCw{
	return me.cw
}
func (me * THSarabunNewBold)GetEnc() string{
	return me.enc
}
func (me * THSarabunNewBold)GetDiff() string {
	return me.diff
}
func (me * THSarabunNewBold) GetOriginalsize() int{
	return 98764
}
func (me * THSarabunNewBold)  SetFamily(family string){
	me.family = family
}
func (me * THSarabunNewBold) 	GetFamily() string{
	return me.family
}