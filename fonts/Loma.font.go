package fonts //change this
import (
	"github.com/signintech/gopdf"
)

type Loma struct {
	family   string
	fonttype string
	name     string
	desc     []gopdf.FontDescItem
	up       int
	ut       int
	cw       gopdf.FontCw
	enc      string
	diff     string
}

func (me *Loma) Init() {
	me.cw = make(gopdf.FontCw)
	me.cw[gopdf.Chr(0)] = 750
	me.cw[gopdf.Chr(1)] = 750
	me.cw[gopdf.Chr(2)] = 750
	me.cw[gopdf.Chr(3)] = 750
	me.cw[gopdf.Chr(4)] = 750
	me.cw[gopdf.Chr(5)] = 750
	me.cw[gopdf.Chr(6)] = 750
	me.cw[gopdf.Chr(7)] = 750
	me.cw[gopdf.Chr(8)] = 750
	me.cw[gopdf.Chr(9)] = 750
	me.cw[gopdf.Chr(10)] = 750
	me.cw[gopdf.Chr(11)] = 750
	me.cw[gopdf.Chr(12)] = 750
	me.cw[gopdf.Chr(13)] = 750
	me.cw[gopdf.Chr(14)] = 750
	me.cw[gopdf.Chr(15)] = 750
	me.cw[gopdf.Chr(16)] = 750
	me.cw[gopdf.Chr(17)] = 750
	me.cw[gopdf.Chr(18)] = 750
	me.cw[gopdf.Chr(19)] = 750
	me.cw[gopdf.Chr(20)] = 750
	me.cw[gopdf.Chr(21)] = 750
	me.cw[gopdf.Chr(22)] = 750
	me.cw[gopdf.Chr(23)] = 750
	me.cw[gopdf.Chr(24)] = 750
	me.cw[gopdf.Chr(25)] = 750
	me.cw[gopdf.Chr(26)] = 750
	me.cw[gopdf.Chr(27)] = 750
	me.cw[gopdf.Chr(28)] = 750
	me.cw[gopdf.Chr(29)] = 750
	me.cw[gopdf.Chr(30)] = 750
	me.cw[gopdf.Chr(31)] = 750
	me.cw[gopdf.ToByte(" ")] = 500
	me.cw[gopdf.ToByte("!")] = 278
	me.cw[gopdf.ToByte("\"")] = 355
	me.cw[gopdf.ToByte("#")] = 556
	me.cw[gopdf.ToByte("$")] = 556
	me.cw[gopdf.ToByte("%")] = 889
	me.cw[gopdf.ToByte("&")] = 667
	me.cw[gopdf.ToByte("'")] = 191
	me.cw[gopdf.ToByte("(")] = 333
	me.cw[gopdf.ToByte(")")] = 333
	me.cw[gopdf.ToByte("*")] = 389
	me.cw[gopdf.ToByte("+")] = 584
	me.cw[gopdf.ToByte(",")] = 278
	me.cw[gopdf.ToByte("-")] = 333
	me.cw[gopdf.ToByte(".")] = 278
	me.cw[gopdf.ToByte("/")] = 291
	me.cw[gopdf.ToByte("0")] = 556
	me.cw[gopdf.ToByte("1")] = 556
	me.cw[gopdf.ToByte("2")] = 556
	me.cw[gopdf.ToByte("3")] = 556
	me.cw[gopdf.ToByte("4")] = 556
	me.cw[gopdf.ToByte("5")] = 556
	me.cw[gopdf.ToByte("6")] = 556
	me.cw[gopdf.ToByte("7")] = 556
	me.cw[gopdf.ToByte("8")] = 556
	me.cw[gopdf.ToByte("9")] = 556
	me.cw[gopdf.ToByte(":")] = 278
	me.cw[gopdf.ToByte(";")] = 278
	me.cw[gopdf.ToByte("<")] = 584
	me.cw[gopdf.ToByte("=")] = 584
	me.cw[gopdf.ToByte(">")] = 584
	me.cw[gopdf.ToByte("?")] = 556
	me.cw[gopdf.ToByte("@")] = 1015
	me.cw[gopdf.ToByte("A")] = 667
	me.cw[gopdf.ToByte("B")] = 667
	me.cw[gopdf.ToByte("C")] = 722
	me.cw[gopdf.ToByte("D")] = 722
	me.cw[gopdf.ToByte("E")] = 667
	me.cw[gopdf.ToByte("F")] = 611
	me.cw[gopdf.ToByte("G")] = 722
	me.cw[gopdf.ToByte("H")] = 722
	me.cw[gopdf.ToByte("I")] = 317
	me.cw[gopdf.ToByte("J")] = 500
	me.cw[gopdf.ToByte("K")] = 667
	me.cw[gopdf.ToByte("L")] = 556
	me.cw[gopdf.ToByte("M")] = 833
	me.cw[gopdf.ToByte("N")] = 722
	me.cw[gopdf.ToByte("O")] = 778
	me.cw[gopdf.ToByte("P")] = 667
	me.cw[gopdf.ToByte("Q")] = 810
	me.cw[gopdf.ToByte("R")] = 684
	me.cw[gopdf.ToByte("S")] = 667
	me.cw[gopdf.ToByte("T")] = 611
	me.cw[gopdf.ToByte("U")] = 722
	me.cw[gopdf.ToByte("V")] = 667
	me.cw[gopdf.ToByte("W")] = 944
	me.cw[gopdf.ToByte("X")] = 667
	me.cw[gopdf.ToByte("Y")] = 667
	me.cw[gopdf.ToByte("Z")] = 611
	me.cw[gopdf.ToByte("[")] = 278
	me.cw[gopdf.ToByte("\\")] = 278
	me.cw[gopdf.ToByte("]")] = 278
	me.cw[gopdf.ToByte("^")] = 469
	me.cw[gopdf.ToByte("_")] = 556
	me.cw[gopdf.ToByte("`")] = 333
	me.cw[gopdf.ToByte("a")] = 556
	me.cw[gopdf.ToByte("b")] = 556
	me.cw[gopdf.ToByte("c")] = 500
	me.cw[gopdf.ToByte("d")] = 556
	me.cw[gopdf.ToByte("e")] = 556
	me.cw[gopdf.ToByte("f")] = 278
	me.cw[gopdf.ToByte("g")] = 556
	me.cw[gopdf.ToByte("h")] = 556
	me.cw[gopdf.ToByte("i")] = 222
	me.cw[gopdf.ToByte("j")] = 222
	me.cw[gopdf.ToByte("k")] = 500
	me.cw[gopdf.ToByte("l")] = 222
	me.cw[gopdf.ToByte("m")] = 833
	me.cw[gopdf.ToByte("n")] = 556
	me.cw[gopdf.ToByte("o")] = 556
	me.cw[gopdf.ToByte("p")] = 556
	me.cw[gopdf.ToByte("q")] = 556
	me.cw[gopdf.ToByte("r")] = 333
	me.cw[gopdf.ToByte("s")] = 500
	me.cw[gopdf.ToByte("t")] = 278
	me.cw[gopdf.ToByte("u")] = 556
	me.cw[gopdf.ToByte("v")] = 500
	me.cw[gopdf.ToByte("w")] = 722
	me.cw[gopdf.ToByte("x")] = 500
	me.cw[gopdf.ToByte("y")] = 500
	me.cw[gopdf.ToByte("z")] = 500
	me.cw[gopdf.ToByte("{")] = 334
	me.cw[gopdf.ToByte("|")] = 260
	me.cw[gopdf.ToByte("}")] = 334
	me.cw[gopdf.ToByte("~")] = 584
	me.cw[gopdf.Chr(127)] = 750
	me.cw[gopdf.Chr(128)] = 750
	me.cw[gopdf.Chr(129)] = 750
	me.cw[gopdf.Chr(130)] = 750
	me.cw[gopdf.Chr(131)] = 750
	me.cw[gopdf.Chr(132)] = 750
	me.cw[gopdf.Chr(133)] = 806
	me.cw[gopdf.Chr(134)] = 750
	me.cw[gopdf.Chr(135)] = 750
	me.cw[gopdf.Chr(136)] = 750
	me.cw[gopdf.Chr(137)] = 750
	me.cw[gopdf.Chr(138)] = 750
	me.cw[gopdf.Chr(139)] = 750
	me.cw[gopdf.Chr(140)] = 750
	me.cw[gopdf.Chr(141)] = 750
	me.cw[gopdf.Chr(142)] = 750
	me.cw[gopdf.Chr(143)] = 750
	me.cw[gopdf.Chr(144)] = 750
	me.cw[gopdf.Chr(145)] = 220
	me.cw[gopdf.Chr(146)] = 283
	me.cw[gopdf.Chr(147)] = 415
	me.cw[gopdf.Chr(148)] = 488
	me.cw[gopdf.Chr(149)] = 464
	me.cw[gopdf.Chr(150)] = 549
	me.cw[gopdf.Chr(151)] = 921
	me.cw[gopdf.Chr(152)] = 750
	me.cw[gopdf.Chr(153)] = 750
	me.cw[gopdf.Chr(154)] = 750
	me.cw[gopdf.Chr(155)] = 750
	me.cw[gopdf.Chr(156)] = 750
	me.cw[gopdf.Chr(157)] = 750
	me.cw[gopdf.Chr(158)] = 750
	me.cw[gopdf.Chr(159)] = 750
	me.cw[gopdf.Chr(160)] = 278
	me.cw[gopdf.Chr(161)] = 605
	me.cw[gopdf.Chr(162)] = 684
	me.cw[gopdf.Chr(163)] = 708
	me.cw[gopdf.Chr(164)] = 669
	me.cw[gopdf.Chr(165)] = 669
	me.cw[gopdf.Chr(166)] = 742
	me.cw[gopdf.Chr(167)] = 488
	me.cw[gopdf.Chr(168)] = 586
	me.cw[gopdf.Chr(169)] = 681
	me.cw[gopdf.Chr(170)] = 679
	me.cw[gopdf.Chr(171)] = 679
	me.cw[gopdf.Chr(172)] = 854
	me.cw[gopdf.Chr(173)] = 852
	me.cw[gopdf.Chr(174)] = 671
	me.cw[gopdf.Chr(175)] = 671
	me.cw[gopdf.Chr(176)] = 552
	me.cw[gopdf.Chr(177)] = 830
	me.cw[gopdf.Chr(178)] = 903
	me.cw[gopdf.Chr(179)] = 928
	me.cw[gopdf.Chr(180)] = 649
	me.cw[gopdf.Chr(181)] = 649
	me.cw[gopdf.Chr(182)] = 605
	me.cw[gopdf.Chr(183)] = 659
	me.cw[gopdf.Chr(184)] = 610
	me.cw[gopdf.Chr(185)] = 684
	me.cw[gopdf.Chr(186)] = 635
	me.cw[gopdf.Chr(187)] = 635
	me.cw[gopdf.Chr(188)] = 586
	me.cw[gopdf.Chr(189)] = 586
	me.cw[gopdf.Chr(190)] = 659
	me.cw[gopdf.Chr(191)] = 708
	me.cw[gopdf.Chr(192)] = 659
	me.cw[gopdf.Chr(193)] = 659
	me.cw[gopdf.Chr(194)] = 586
	me.cw[gopdf.Chr(195)] = 537
	me.cw[gopdf.Chr(196)] = 605
	me.cw[gopdf.Chr(197)] = 613
	me.cw[gopdf.Chr(198)] = 659
	me.cw[gopdf.Chr(199)] = 562
	me.cw[gopdf.Chr(200)] = 635
	me.cw[gopdf.Chr(201)] = 659
	me.cw[gopdf.Chr(202)] = 610
	me.cw[gopdf.Chr(203)] = 659
	me.cw[gopdf.Chr(204)] = 684
	me.cw[gopdf.Chr(205)] = 601
	me.cw[gopdf.Chr(206)] = 610
	me.cw[gopdf.Chr(207)] = 562
	me.cw[gopdf.Chr(208)] = 537
	me.cw[gopdf.Chr(209)] = 0
	me.cw[gopdf.Chr(210)] = 562
	me.cw[gopdf.Chr(211)] = 562
	me.cw[gopdf.Chr(212)] = 0
	me.cw[gopdf.Chr(213)] = 0
	me.cw[gopdf.Chr(214)] = 0
	me.cw[gopdf.Chr(215)] = 0
	me.cw[gopdf.Chr(216)] = 0
	me.cw[gopdf.Chr(217)] = 0
	me.cw[gopdf.Chr(218)] = 0
	me.cw[gopdf.Chr(219)] = 750
	me.cw[gopdf.Chr(220)] = 750
	me.cw[gopdf.Chr(221)] = 750
	me.cw[gopdf.Chr(222)] = 750
	me.cw[gopdf.Chr(223)] = 610
	me.cw[gopdf.Chr(224)] = 342
	me.cw[gopdf.Chr(225)] = 645
	me.cw[gopdf.Chr(226)] = 537
	me.cw[gopdf.Chr(227)] = 488
	me.cw[gopdf.Chr(228)] = 503
	me.cw[gopdf.Chr(229)] = 488
	me.cw[gopdf.Chr(230)] = 537
	me.cw[gopdf.Chr(231)] = 0
	me.cw[gopdf.Chr(232)] = 0
	me.cw[gopdf.Chr(233)] = 0
	me.cw[gopdf.Chr(234)] = 0
	me.cw[gopdf.Chr(235)] = 0
	me.cw[gopdf.Chr(236)] = 0
	me.cw[gopdf.Chr(237)] = 0
	me.cw[gopdf.Chr(238)] = 0
	me.cw[gopdf.Chr(239)] = 610
	me.cw[gopdf.Chr(240)] = 610
	me.cw[gopdf.Chr(241)] = 635
	me.cw[gopdf.Chr(242)] = 659
	me.cw[gopdf.Chr(243)] = 684
	me.cw[gopdf.Chr(244)] = 757
	me.cw[gopdf.Chr(245)] = 757
	me.cw[gopdf.Chr(246)] = 635
	me.cw[gopdf.Chr(247)] = 752
	me.cw[gopdf.Chr(248)] = 771
	me.cw[gopdf.Chr(249)] = 732
	me.cw[gopdf.Chr(250)] = 684
	me.cw[gopdf.Chr(251)] = 1157
	me.cw[gopdf.Chr(252)] = 750
	me.cw[gopdf.Chr(253)] = 750
	me.cw[gopdf.Chr(254)] = 750
	me.cw[gopdf.Chr(255)] = 750
	me.up = -88
	me.ut = 0
	me.fonttype = "TrueType"
	me.name = "Loma"
	me.enc = "cp874"
	me.diff = "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef"
	me.desc = make([]gopdf.FontDescItem, 8)
	me.desc[0] = gopdf.FontDescItem{Key: "Ascent", Val: "0"}
	me.desc[1] = gopdf.FontDescItem{Key: "Descent", Val: "-200"}
	me.desc[2] = gopdf.FontDescItem{Key: "CapHeight", Val: "0"}
	me.desc[3] = gopdf.FontDescItem{Key: "Flags", Val: "33"}
	me.desc[4] = gopdf.FontDescItem{Key: "FontBBox", Val: "[31257 31560 1338 1146]"}
	me.desc[5] = gopdf.FontDescItem{Key: "ItalicAngle", Val: "0"}
	me.desc[6] = gopdf.FontDescItem{Key: "StemV", Val: "70"}
	me.desc[7] = gopdf.FontDescItem{Key: "MissingWidth", Val: "750"}
}
func (me *Loma) GetType() string {
	return me.fonttype
}
func (me *Loma) GetName() string {
	return me.name
}
func (me *Loma) GetDesc() []gopdf.FontDescItem {
	return me.desc
}
func (me *Loma) GetUp() int {
	return me.up
}
func (me *Loma) GetUt() int {
	return me.ut
}
func (me *Loma) GetCw() gopdf.FontCw {
	return me.cw
}
func (me *Loma) GetEnc() string {
	return me.enc
}
func (me *Loma) GetDiff() string {
	return me.diff
}
func (me *Loma) GetOriginalsize() int {
	return 98764
}
func (me *Loma) SetFamily(family string) {
	me.family = family
}
func (me *Loma) GetFamily() string {
	return me.family
}
