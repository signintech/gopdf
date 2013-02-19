package fonts

import (

)

type THSarabun struct{
	family string
	fonttype string 
	name string
	desc  []FontDescItem
	up int
	ut int
	cw map[string]int
	enc string
	diff string
}

func (me * THSarabun)Init(){

	me.cw = make(map[string]int)
	//me.cw = cw
	me.cw[Chr(0)]=692;
	me.cw[Chr(1)]=692;
	me.cw[Chr(2)]=692;
	me.cw[Chr(3)]=692;
	me.cw[Chr(4)]=692;
	me.cw[Chr(5)]=692;
	me.cw[Chr(6)]=692;
	me.cw[Chr(7)]=692;
	me.cw[Chr(8)]=692;
	me.cw[Chr(9)]=692;
	me.cw[Chr(10)]=692;
	me.cw[Chr(11)]=692;
	me.cw[Chr(12)]=692;
	me.cw[Chr(13)]=692;
	me.cw[Chr(14)]=692;
	me.cw[Chr(15)]=692;
	me.cw[Chr(16)]=692;
	me.cw[Chr(17)]=692;
	me.cw[Chr(18)]=692;
	me.cw[Chr(19)]=692;
	me.cw[Chr(20)]=692;
	me.cw[Chr(21)]=692;
	me.cw[Chr(22)]=692;
	me.cw[Chr(23)]=692;
	me.cw[Chr(24)]=692;
	me.cw[Chr(25)]=692;
	me.cw[Chr(26)]=692;
	me.cw[Chr(27)]=692;
	me.cw[Chr(28)]=692;
	me.cw[Chr(29)]=692;
	me.cw[Chr(30)]=692;
	me.cw[Chr(31)]=692;
	me.cw[" "]=216;
	me.cw["!"]=147;
	me.cw["\""]=208;
	me.cw["#"]=403;
	me.cw["$"]=361;
	me.cw["%"]=585;
	me.cw["&"]=423;
	me.cw["'"]=120;
	me.cw["("]=190;
	me.cw[")"]=190;
	me.cw["*"]=285;
	me.cw["+"]=411;
	me.cw[","]=162;
	me.cw["-"]=216;
	me.cw["."]=162;
	me.cw["/"]=270;
	me.cw["0"]=362;
	me.cw["1"]=362;
	me.cw["2"]=362;
	me.cw["3"]=362;
	me.cw["4"]=362;
	me.cw["5"]=362;
	me.cw["6"]=362;
	me.cw["7"]=362;
	me.cw["8"]=362;
	me.cw["9"]=362;
	me.cw[":"]=162;
	me.cw[";"]=162;
	me.cw["<"]=411;
	me.cw["="]=411;
	me.cw[">"]=411;
	me.cw["?"]=283;
	me.cw["@"]=536;
	me.cw["A"]=400;
	me.cw["B"]=378;
	me.cw["C"]=406;
	me.cw["D"]=431;
	me.cw["E"]=351;
	me.cw["F"]=351;
	me.cw["G"]=425;
	me.cw["H"]=441;
	me.cw["I"]=147;
	me.cw["J"]=264;
	me.cw["K"]=376;
	me.cw["L"]=353;
	me.cw["M"]=548;
	me.cw["N"]=441;
	me.cw["O"]=486;
	me.cw["P"]=378;
	me.cw["Q"]=487;
	me.cw["R"]=379;
	me.cw["S"]=352;
	me.cw["T"]=379;
	me.cw["U"]=466;
	me.cw["V"]=390;
	me.cw["W"]=588;
	me.cw["X"]=418;
	me.cw["Y"]=366;
	me.cw["Z"]=424;
	me.cw["["]=196;
	me.cw["\\"]=262;
	me.cw["]"]=196;
	me.cw["^"]=412;
	me.cw["_"]=352;
	me.cw["`"]=204;
	me.cw["a"]=344;
	me.cw["b"]=401;
	me.cw["c"]=331;
	me.cw["d"]=401;
	me.cw["e"]=374;
	me.cw["f"]=206;
	me.cw["g"]=311;
	me.cw["h"]=390;
	me.cw["i"]=143;
	me.cw["j"]=155;
	me.cw["k"]=316;
	me.cw["l"]=200;
	me.cw["m"]=601;
	me.cw["n"]=390;
	me.cw["o"]=398;
	me.cw["p"]=401;
	me.cw["q"]=401;
	me.cw["r"]=217;
	me.cw["s"]=282;
	me.cw["t"]=238;
	me.cw["u"]=390;
	me.cw["v"]=341;
	me.cw["w"]=507;
	me.cw["x"]=318;
	me.cw["y"]=337;
	me.cw["z"]=321;
	me.cw["{"]=208;
	me.cw["|"]=153;
	me.cw["}"]=208;
	me.cw["~"]=416;
	me.cw[Chr(127)]=692;
	me.cw[Chr(128)]=406;
	me.cw[Chr(129)]=692;
	me.cw[Chr(130)]=692;
	me.cw[Chr(131)]=692;
	me.cw[Chr(132)]=692;
	me.cw[Chr(133)]=479;
	me.cw[Chr(134)]=692;
	me.cw[Chr(135)]=692;
	me.cw[Chr(136)]=692;
	me.cw[Chr(137)]=692;
	me.cw[Chr(138)]=692;
	me.cw[Chr(139)]=692;
	me.cw[Chr(140)]=692;
	me.cw[Chr(141)]=692;
	me.cw[Chr(142)]=692;
	me.cw[Chr(143)]=692;
	me.cw[Chr(144)]=692;
	me.cw[Chr(145)]=247;
	me.cw[Chr(146)]=247;
	me.cw[Chr(147)]=370;
	me.cw[Chr(148)]=370;
	me.cw[Chr(149)]=216;
	me.cw[Chr(150)]=360;
	me.cw[Chr(151)]=720;
	me.cw[Chr(152)]=692;
	me.cw[Chr(153)]=692;
	me.cw[Chr(154)]=692;
	me.cw[Chr(155)]=692;
	me.cw[Chr(156)]=692;
	me.cw[Chr(157)]=692;
	me.cw[Chr(158)]=692;
	me.cw[Chr(159)]=692;
	me.cw[Chr(160)]=216;
	me.cw[Chr(161)]=386;
	me.cw[Chr(162)]=378;
	me.cw[Chr(163)]=382;
	me.cw[Chr(164)]=393;
	me.cw[Chr(165)]=393;
	me.cw[Chr(166)]=408;
	me.cw[Chr(167)]=294;
	me.cw[Chr(168)]=367;
	me.cw[Chr(169)]=377;
	me.cw[Chr(170)]=380;
	me.cw[Chr(171)]=384;
	me.cw[Chr(172)]=519;
	me.cw[Chr(173)]=519;
	me.cw[Chr(174)]=425;
	me.cw[Chr(175)]=425;
	me.cw[Chr(176)]=343;
	me.cw[Chr(177)]=461;
	me.cw[Chr(178)]=532;
	me.cw[Chr(179)]=543;
	me.cw[Chr(180)]=391;
	me.cw[Chr(181)]=391;
	me.cw[Chr(182)]=378;
	me.cw[Chr(183)]=430;
	me.cw[Chr(184)]=335;
	me.cw[Chr(185)]=420;
	me.cw[Chr(186)]=428;
	me.cw[Chr(187)]=428;
	me.cw[Chr(188)]=381;
	me.cw[Chr(189)]=381;
	me.cw[Chr(190)]=447;
	me.cw[Chr(191)]=447;
	me.cw[Chr(192)]=425;
	me.cw[Chr(193)]=400;
	me.cw[Chr(194)]=375;
	me.cw[Chr(195)]=322;
	me.cw[Chr(196)]=378;
	me.cw[Chr(197)]=381;
	me.cw[Chr(198)]=425;
	me.cw[Chr(199)]=335;
	me.cw[Chr(200)]=393;
	me.cw[Chr(201)]=438;
	me.cw[Chr(202)]=381;
	me.cw[Chr(203)]=427;
	me.cw[Chr(204)]=454;
	me.cw[Chr(205)]=387;
	me.cw[Chr(206)]=372;
	me.cw[Chr(207)]=391;
	me.cw[Chr(208)]=357;
	me.cw[Chr(209)]=0;
	me.cw[Chr(210)]=316;
	me.cw[Chr(211)]=316;
	me.cw[Chr(212)]=0;
	me.cw[Chr(213)]=0;
	me.cw[Chr(214)]=0;
	me.cw[Chr(215)]=0;
	me.cw[Chr(216)]=0;
	me.cw[Chr(217)]=0;
	me.cw[Chr(218)]=0;
	me.cw[Chr(219)]=692;
	me.cw[Chr(220)]=692;
	me.cw[Chr(221)]=692;
	me.cw[Chr(222)]=692;
	me.cw[Chr(223)]=411;
	me.cw[Chr(224)]=203;
	me.cw[Chr(225)]=377;
	me.cw[Chr(226)]=237;
	me.cw[Chr(227)]=242;
	me.cw[Chr(228)]=244;
	me.cw[Chr(229)]=205;
	me.cw[Chr(230)]=399;
	me.cw[Chr(231)]=0;
	me.cw[Chr(232)]=0;
	me.cw[Chr(233)]=0;
	me.cw[Chr(234)]=0;
	me.cw[Chr(235)]=0;
	me.cw[Chr(236)]=0;
	me.cw[Chr(237)]=0;
	me.cw[Chr(238)]=0;
	me.cw[Chr(239)]=450;
	me.cw[Chr(240)]=449;
	me.cw[Chr(241)]=449;
	me.cw[Chr(242)]=449;
	me.cw[Chr(243)]=449;
	me.cw[Chr(244)]=449;
	me.cw[Chr(245)]=449;
	me.cw[Chr(246)]=449;
	me.cw[Chr(247)]=449;
	me.cw[Chr(248)]=449;
	me.cw[Chr(249)]=449;
	me.cw[Chr(250)]=522;
	me.cw[Chr(251)]=697;
	me.cw[Chr(252)]=692;
	me.cw[Chr(253)]=692;
	me.cw[Chr(254)]=692;
	me.cw[Chr(255)]=692;
	
	me.up = -35;
	me.ut = 30;
	me.fonttype = "TrueType"
	me.name = "THSarabunPSK"
	me.enc = "cp874"
	me.diff =   "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef /.notdef"

	me.desc = make([]FontDescItem,8)
	me.desc[0] =  FontDescItem{ Key:"Ascent",Val : "850" }
	me.desc[1] =  FontDescItem{ Key: "Descent", Val : "-250" }
	me.desc[2] =  FontDescItem{ Key:"CapHeight", Val :  "476"}
	me.desc[3] =  FontDescItem{ Key: "Flags", Val :  "32"}
	me.desc[4] =  FontDescItem{ Key:"FontBBox", Val :  "[-427 -421 947 836]"}
	me.desc[5] =  FontDescItem{ Key:"ItalicAngle", Val :  "0"}
	me.desc[6] =  FontDescItem{ Key:"StemV", Val :  "70"}
	me.desc[7] =  FontDescItem{ Key:"MissingWidth", Val :  "692"}
}
func (me * THSarabun)GetType() string{
	return me.fonttype
}
func (me * THSarabun)GetName() string{
	return me.name
}	
func (me * THSarabun)GetDesc() []FontDescItem{
	return me.desc
}
func (me * THSarabun)GetUp() int{
	return me.up
}
func (me * THSarabun)GetUt()  int{
	return me.ut
}
func (me * THSarabun)GetCw() map[string]int{
	return me.cw
}
func (me * THSarabun)GetEnc() string{
	return me.enc
}
func (me * THSarabun)GetDiff() string {
	return me.diff
}

func (me * THSarabun) GetOriginalsize() int{
	return 98764
}

func (me * THSarabun)  SetFamily(family string){
	me.family = family
}

func (me * THSarabun) 	GetFamily() string{
	return me.family
}
