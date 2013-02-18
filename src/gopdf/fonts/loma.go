package fonts

import (

)

/*ยังไม่เสร็จ*/

type Loma struct{
	fonttype string 
	name string
	desc map[string]string
	up int
	ut int
	cw map[string]int
	enc string
	diff string
}

func (me * Loma)Init(){

	me.cw[Chr(0)]=750;
	me.cw[Chr(1)]=750;
	me.cw[Chr(2)]=750;
	me.cw[Chr(3)]=750;
	me.cw[Chr(4)]=750;
	me.cw[Chr(5)]=750;
	me.cw[Chr(6)]=750;
	me.cw[Chr(7)]=750;
	me.cw[Chr(8)]=750;
	me.cw[Chr(9)]=750;
	me.cw[Chr(10)]=750;
	me.cw[Chr(11)]=750;
	me.cw[Chr(12)]=750;
	me.cw[Chr(13)]=750;
	me.cw[Chr(14)]=750;
	me.cw[Chr(15)]=750;
	me.cw[Chr(16)]=750;
	me.cw[Chr(17)]=750;
	me.cw[Chr(18)]=750;
	me.cw[Chr(19)]=750;
	me.cw[Chr(20)]=750;
	me.cw[Chr(21)]=750;
	me.cw[Chr(22)]=750;
	me.cw[Chr(23)]=750;
	me.cw[Chr(24)]=750;
	me.cw[Chr(25)]=750;
	me.cw[Chr(26)]=750;
	me.cw[Chr(27)]=750;
	me.cw[Chr(28)]=750;
	me.cw[Chr(29)]=750;
	me.cw[Chr(30)]=750;
	me.cw[Chr(31)]=750;
	me.cw[" "]=500;
	me.cw["!"]=278;
	me.cw["\""]=355;
	me.cw["#"]=556;
	me.cw["$"]=556;
	me.cw["%"]=889;
	me.cw["&"]=667;
	me.cw["'"]=191;
	me.cw["("]=333;
	me.cw[")"]=333;
	me.cw["*"]=389;
	me.cw["+"]=584;
	me.cw[","]=278;
	me.cw["-"]=333;
	me.cw["."]=278;
	me.cw["/"]=291;
	me.cw["0"]=556;
	me.cw["1"]=556;
	me.cw["2"]=556;
	me.cw["3"]=556;
	me.cw["4"]=556;
	me.cw["5"]=556;
	me.cw["6"]=556;
	me.cw["7"]=556;
	me.cw["8"]=556;
	me.cw["9"]=556;
	me.cw[":"]=278;
	me.cw[";"]=278;
	me.cw["<"]=584;
	me.cw["="]=584;
	me.cw[">"]=584;
	me.cw["?"]=556;
	me.cw["@"]=1015;
	me.cw["A"]=667;
	me.cw["B"]=667;
	me.cw["C"]=722;
	me.cw["D"]=722;
	me.cw["E"]=667;
	me.cw["F"]=611;
	me.cw["G"]=722;
	me.cw["H"]=722;
	me.cw["I"]=317;
	me.cw["J"]=500;
	me.cw["K"]=667;
	me.cw["L"]=556;
	me.cw["M"]=833;
	me.cw["N"]=722;
	me.cw["O"]=778;
	me.cw["P"]=667;
	me.cw["Q"]=810;
	me.cw["R"]=684;
	me.cw["S"]=667;
	me.cw["T"]=611;
	me.cw["U"]=722;
	me.cw["V"]=667;
	me.cw["W"]=944;
	me.cw["X"]=667;
	me.cw["Y"]=667;
	me.cw["Z"]=611;
	me.cw["["]=278;
	me.cw["\\"]=278;
	me.cw["]"]=278;
	me.cw["^"]=469;
	me.cw["_"]=556;
	me.cw["`"]=333;
	me.cw["a"]=556;
	me.cw["b"]=556;
	me.cw["c"]=500;
	me.cw["d"]=556;
	me.cw["e"]=556;
	me.cw["f"]=278;
	me.cw["g"]=556;
	me.cw["h"]=556;
	me.cw["i"]=222;
	me.cw["j"]=222;
	me.cw["k"]=500;
	me.cw["l"]=222;
	me.cw["m"]=833;
	me.cw["n"]=556;
	me.cw["o"]=556;
	me.cw["p"]=556;
	me.cw["q"]=556;
	me.cw["r"]=333;
	me.cw["s"]=500;
	me.cw["t"]=278;
	me.cw["u"]=556;
	me.cw["v"]=500;
	me.cw["w"]=722;
	me.cw["x"]=500;
	me.cw["y"]=500;
	me.cw["z"]=500;
	me.cw["{"]=334;
	me.cw["|"]=260;
	me.cw["}"]=334;
	me.cw["~"]=584;
	me.cw[Chr(127)]=750;
	me.cw[Chr(128)]=750;
	me.cw[Chr(129)]=750;
	me.cw[Chr(130)]=750;
	me.cw[Chr(131)]=750;
	me.cw[Chr(132)]=750;
	me.cw[Chr(133)]=806;
	me.cw[Chr(134)]=750;
	me.cw[Chr(135)]=750;
	me.cw[Chr(136)]=750;
	me.cw[Chr(137)]=750;
	me.cw[Chr(138)]=750;
	me.cw[Chr(139)]=750;
	me.cw[Chr(140)]=750;
	me.cw[Chr(141)]=750;
	me.cw[Chr(142)]=750;
	me.cw[Chr(143)]=750;
	me.cw[Chr(144)]=750;
	me.cw[Chr(145)]=220;
	me.cw[Chr(146)]=283;
	me.cw[Chr(147)]=415;
	me.cw[Chr(148)]=488;
	me.cw[Chr(149)]=464;
	me.cw[Chr(150)]=549;
	me.cw[Chr(151)]=921;
	me.cw[Chr(152)]=750;
	me.cw[Chr(153)]=750;
	me.cw[Chr(154)]=750;
	me.cw[Chr(155)]=750;
	me.cw[Chr(156)]=750;
	me.cw[Chr(157)]=750;
	me.cw[Chr(158)]=750;
	me.cw[Chr(159)]=750;
	me.cw[Chr(160)]=278;
	me.cw[Chr(161)]=605;
	me.cw[Chr(162)]=684;
	me.cw[Chr(163)]=708;
	me.cw[Chr(164)]=669;
	me.cw[Chr(165)]=669;
	me.cw[Chr(166)]=742;
	me.cw[Chr(167)]=488;
	me.cw[Chr(168)]=586;
	me.cw[Chr(169)]=681;
	me.cw[Chr(170)]=679;
	me.cw[Chr(171)]=679;
	me.cw[Chr(172)]=854;
	me.cw[Chr(173)]=852;
	me.cw[Chr(174)]=671;
	me.cw[Chr(175)]=671;
	me.cw[Chr(176)]=552;
	me.cw[Chr(177)]=830;
	me.cw[Chr(178)]=903;
	me.cw[Chr(179)]=928;
	me.cw[Chr(180)]=649;
	me.cw[Chr(181)]=649;
	me.cw[Chr(182)]=605;
	me.cw[Chr(183)]=659;
	me.cw[Chr(184)]=610;
	me.cw[Chr(185)]=684;
	me.cw[Chr(186)]=635;
	me.cw[Chr(187)]=635;
	me.cw[Chr(188)]=586;
	me.cw[Chr(189)]=586;
	me.cw[Chr(190)]=659;
	me.cw[Chr(191)]=708;
	me.cw[Chr(192)]=659;
	me.cw[Chr(193)]=659;
	me.cw[Chr(194)]=586;
	me.cw[Chr(195)]=537;
	me.cw[Chr(196)]=605;
	me.cw[Chr(197)]=613;
	me.cw[Chr(198)]=659;
	me.cw[Chr(199)]=562;
	me.cw[Chr(200)]=635;
	me.cw[Chr(201)]=659;
	me.cw[Chr(202)]=610;
	me.cw[Chr(203)]=659;
	me.cw[Chr(204)]=684;
	me.cw[Chr(205)]=601;
	me.cw[Chr(206)]=610;
	me.cw[Chr(207)]=562;
	me.cw[Chr(208)]=537;
	me.cw[Chr(209)]=0;
	me.cw[Chr(210)]=562;
	me.cw[Chr(211)]=562;
	me.cw[Chr(212)]=0;
	me.cw[Chr(213)]=0;
	me.cw[Chr(214)]=0;
	me.cw[Chr(215)]=0;
	me.cw[Chr(216)]=0;
	me.cw[Chr(217)]=0;
	me.cw[Chr(218)]=0;
	me.cw[Chr(219)]=750;
	me.cw[Chr(220)]=750;
	me.cw[Chr(221)]=750;
	me.cw[Chr(222)]=750;
	me.cw[Chr(223)]=610;
	me.cw[Chr(224)]=342;
	me.cw[Chr(225)]=645;
	me.cw[Chr(226)]=537;
	me.cw[Chr(227)]=488;
	me.cw[Chr(228)]=503;
	me.cw[Chr(229)]=488;
	me.cw[Chr(230)]=537;
	me.cw[Chr(231)]=0;
	me.cw[Chr(232)]=0;
	me.cw[Chr(233)]=0;
	me.cw[Chr(234)]=0;
	me.cw[Chr(235)]=0;
	me.cw[Chr(236)]=0;
	me.cw[Chr(237)]=0;
	me.cw[Chr(238)]=0;
	me.cw[Chr(239)]=610;
	me.cw[Chr(240)]=610;
	me.cw[Chr(241)]=635;
	me.cw[Chr(242)]=659;
	me.cw[Chr(243)]=684;
	me.cw[Chr(244)]=757;
	me.cw[Chr(245)]=757;
	me.cw[Chr(246)]=635;
	me.cw[Chr(247)]=752;
	me.cw[Chr(248)]=771;
	me.cw[Chr(249)]=732;
	me.cw[Chr(250)]=684;
	me.cw[Chr(251)]=1157;
	me.cw[Chr(252)]=750;
	me.cw[Chr(253)]=750;
	me.cw[Chr(254)]=750;
	me.cw[Chr(255)]=750;
	
	
	me.diff =  "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef /.notdef"
	
	

}
func (me * Loma)GetType() string{
	return me.fonttype
}
func (me * Loma)GetName() string{
	return me.name
}	
func (me * Loma)GetDesc() map[string]string{
	return me.desc
}
func (me * Loma)GetUp() int{
	return me.up
}
func (me * Loma)GetUt()  int{
	return me.ut
}
func (me * Loma)GetCw() map[string]int{
	return me.cw
}
func (me * Loma)GetEnc() string{
	return me.enc
}
func (me * Loma)GetDiff() string {
	return me.diff
}
func (me * Loma) GetOriginalsize() int{
	return 0
}
