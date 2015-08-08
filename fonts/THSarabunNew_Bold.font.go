package fonts

import (
	"github.com/signintech/gopdf"
)

type THSarabunNewBold struct {
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

func (me *THSarabunNewBold) Init() {
	me.cw = make(gopdf.FontCw)
	me.cw[gopdf.Chr(0)] = 698
	me.cw[gopdf.Chr(1)] = 698
	me.cw[gopdf.Chr(2)] = 698
	me.cw[gopdf.Chr(3)] = 698
	me.cw[gopdf.Chr(4)] = 698
	me.cw[gopdf.Chr(5)] = 698
	me.cw[gopdf.Chr(6)] = 698
	me.cw[gopdf.Chr(7)] = 698
	me.cw[gopdf.Chr(8)] = 698
	me.cw[gopdf.Chr(9)] = 698
	me.cw[gopdf.Chr(10)] = 698
	me.cw[gopdf.Chr(11)] = 698
	me.cw[gopdf.Chr(12)] = 698
	me.cw[gopdf.Chr(13)] = 698
	me.cw[gopdf.Chr(14)] = 698
	me.cw[gopdf.Chr(15)] = 698
	me.cw[gopdf.Chr(16)] = 698
	me.cw[gopdf.Chr(17)] = 698
	me.cw[gopdf.Chr(18)] = 698
	me.cw[gopdf.Chr(19)] = 698
	me.cw[gopdf.Chr(20)] = 698
	me.cw[gopdf.Chr(21)] = 698
	me.cw[gopdf.Chr(22)] = 698
	me.cw[gopdf.Chr(23)] = 698
	me.cw[gopdf.Chr(24)] = 698
	me.cw[gopdf.Chr(25)] = 698
	me.cw[gopdf.Chr(26)] = 698
	me.cw[gopdf.Chr(27)] = 698
	me.cw[gopdf.Chr(28)] = 698
	me.cw[gopdf.Chr(29)] = 698
	me.cw[gopdf.Chr(30)] = 698
	me.cw[gopdf.Chr(31)] = 698
	me.cw[gopdf.ToByte(" ")] = 226
	me.cw[gopdf.ToByte("!")] = 168
	me.cw[gopdf.ToByte("\"")] = 291
	me.cw[gopdf.ToByte("#")] = 462
	me.cw[gopdf.ToByte("$")] = 367
	me.cw[gopdf.ToByte("%")] = 688
	me.cw[gopdf.ToByte("&")] = 486
	me.cw[gopdf.ToByte("'")] = 169
	me.cw[gopdf.ToByte("(")] = 244
	me.cw[gopdf.ToByte(")")] = 244
	me.cw[gopdf.ToByte("*")] = 343
	me.cw[gopdf.ToByte("+")] = 417
	me.cw[gopdf.ToByte(",")] = 172
	me.cw[gopdf.ToByte("-")] = 226
	me.cw[gopdf.ToByte(".")] = 172
	me.cw[gopdf.ToByte("/")] = 276
	me.cw[gopdf.ToByte("0")] = 378
	me.cw[gopdf.ToByte("1")] = 378
	me.cw[gopdf.ToByte("2")] = 378
	me.cw[gopdf.ToByte("3")] = 378
	me.cw[gopdf.ToByte("4")] = 378
	me.cw[gopdf.ToByte("5")] = 378
	me.cw[gopdf.ToByte("6")] = 378
	me.cw[gopdf.ToByte("7")] = 378
	me.cw[gopdf.ToByte("8")] = 378
	me.cw[gopdf.ToByte("9")] = 378
	me.cw[gopdf.ToByte(":")] = 172
	me.cw[gopdf.ToByte(";")] = 172
	me.cw[gopdf.ToByte("<")] = 417
	me.cw[gopdf.ToByte("=")] = 417
	me.cw[gopdf.ToByte(">")] = 417
	me.cw[gopdf.ToByte("?")] = 289
	me.cw[gopdf.ToByte("@")] = 572
	me.cw[gopdf.ToByte("A")] = 451
	me.cw[gopdf.ToByte("B")] = 396
	me.cw[gopdf.ToByte("C")] = 432
	me.cw[gopdf.ToByte("D")] = 459
	me.cw[gopdf.ToByte("E")] = 358
	me.cw[gopdf.ToByte("F")] = 358
	me.cw[gopdf.ToByte("G")] = 454
	me.cw[gopdf.ToByte("H")] = 472
	me.cw[gopdf.ToByte("I")] = 163
	me.cw[gopdf.ToByte("J")] = 282
	me.cw[gopdf.ToByte("K")] = 411
	me.cw[gopdf.ToByte("L")] = 359
	me.cw[gopdf.ToByte("M")] = 574
	me.cw[gopdf.ToByte("N")] = 472
	me.cw[gopdf.ToByte("O")] = 518
	me.cw[gopdf.ToByte("P")] = 390
	me.cw[gopdf.ToByte("Q")] = 518
	me.cw[gopdf.ToByte("R")] = 395
	me.cw[gopdf.ToByte("S")] = 365
	me.cw[gopdf.ToByte("T")] = 421
	me.cw[gopdf.ToByte("U")] = 482
	me.cw[gopdf.ToByte("V")] = 422
	me.cw[gopdf.ToByte("W")] = 612
	me.cw[gopdf.ToByte("X")] = 426
	me.cw[gopdf.ToByte("Y")] = 396
	me.cw[gopdf.ToByte("Z")] = 451
	me.cw[gopdf.ToByte("[")] = 234
	me.cw[gopdf.ToByte("\\")] = 268
	me.cw[gopdf.ToByte("]")] = 234
	me.cw[gopdf.ToByte("^")] = 418
	me.cw[gopdf.ToByte("_")] = 376
	me.cw[gopdf.ToByte("`")] = 210
	me.cw[gopdf.ToByte("a")] = 367
	me.cw[gopdf.ToByte("b")] = 435
	me.cw[gopdf.ToByte("c")] = 347
	me.cw[gopdf.ToByte("d")] = 435
	me.cw[gopdf.ToByte("e")] = 391
	me.cw[gopdf.ToByte("f")] = 254
	me.cw[gopdf.ToByte("g")] = 332
	me.cw[gopdf.ToByte("h")] = 423
	me.cw[gopdf.ToByte("i")] = 163
	me.cw[gopdf.ToByte("j")] = 171
	me.cw[gopdf.ToByte("k")] = 340
	me.cw[gopdf.ToByte("l")] = 213
	me.cw[gopdf.ToByte("m")] = 636
	me.cw[gopdf.ToByte("n")] = 423
	me.cw[gopdf.ToByte("o")] = 412
	me.cw[gopdf.ToByte("p")] = 431
	me.cw[gopdf.ToByte("q")] = 431
	me.cw[gopdf.ToByte("r")] = 246
	me.cw[gopdf.ToByte("s")] = 296
	me.cw[gopdf.ToByte("t")] = 258
	me.cw[gopdf.ToByte("u")] = 418
	me.cw[gopdf.ToByte("v")] = 358
	me.cw[gopdf.ToByte("w")] = 480
	me.cw[gopdf.ToByte("x")] = 349
	me.cw[gopdf.ToByte("y")] = 361
	me.cw[gopdf.ToByte("z")] = 358
	me.cw[gopdf.ToByte("{")] = 255
	me.cw[gopdf.ToByte("|")] = 177
	me.cw[gopdf.ToByte("}")] = 255
	me.cw[gopdf.ToByte("~")] = 422
	me.cw[gopdf.Chr(127)] = 698
	me.cw[gopdf.Chr(128)] = 518
	me.cw[gopdf.Chr(129)] = 698
	me.cw[gopdf.Chr(130)] = 698
	me.cw[gopdf.Chr(131)] = 698
	me.cw[gopdf.Chr(132)] = 698
	me.cw[gopdf.Chr(133)] = 504
	me.cw[gopdf.Chr(134)] = 698
	me.cw[gopdf.Chr(135)] = 698
	me.cw[gopdf.Chr(136)] = 698
	me.cw[gopdf.Chr(137)] = 698
	me.cw[gopdf.Chr(138)] = 698
	me.cw[gopdf.Chr(139)] = 698
	me.cw[gopdf.Chr(140)] = 698
	me.cw[gopdf.Chr(141)] = 698
	me.cw[gopdf.Chr(142)] = 698
	me.cw[gopdf.Chr(143)] = 698
	me.cw[gopdf.Chr(144)] = 698
	me.cw[gopdf.Chr(145)] = 253
	me.cw[gopdf.Chr(146)] = 253
	me.cw[gopdf.Chr(147)] = 384
	me.cw[gopdf.Chr(148)] = 384
	me.cw[gopdf.Chr(149)] = 226
	me.cw[gopdf.Chr(150)] = 335
	me.cw[gopdf.Chr(151)] = 693
	me.cw[gopdf.Chr(152)] = 698
	me.cw[gopdf.Chr(153)] = 698
	me.cw[gopdf.Chr(154)] = 698
	me.cw[gopdf.Chr(155)] = 698
	me.cw[gopdf.Chr(156)] = 698
	me.cw[gopdf.Chr(157)] = 698
	me.cw[gopdf.Chr(158)] = 698
	me.cw[gopdf.Chr(159)] = 698
	me.cw[gopdf.Chr(160)] = 226
	me.cw[gopdf.Chr(161)] = 391
	me.cw[gopdf.Chr(162)] = 415
	me.cw[gopdf.Chr(163)] = 431
	me.cw[gopdf.Chr(164)] = 403
	me.cw[gopdf.Chr(165)] = 403
	me.cw[gopdf.Chr(166)] = 453
	me.cw[gopdf.Chr(167)] = 318
	me.cw[gopdf.Chr(168)] = 390
	me.cw[gopdf.Chr(169)] = 405
	me.cw[gopdf.Chr(170)] = 415
	me.cw[gopdf.Chr(171)] = 431
	me.cw[gopdf.Chr(172)] = 565
	me.cw[gopdf.Chr(173)] = 565
	me.cw[gopdf.Chr(174)] = 425
	me.cw[gopdf.Chr(175)] = 425
	me.cw[gopdf.Chr(176)] = 372
	me.cw[gopdf.Chr(177)] = 492
	me.cw[gopdf.Chr(178)] = 551
	me.cw[gopdf.Chr(179)] = 576
	me.cw[gopdf.Chr(180)] = 405
	me.cw[gopdf.Chr(181)] = 405
	me.cw[gopdf.Chr(182)] = 391
	me.cw[gopdf.Chr(183)] = 438
	me.cw[gopdf.Chr(184)] = 358
	me.cw[gopdf.Chr(185)] = 433
	me.cw[gopdf.Chr(186)] = 445
	me.cw[gopdf.Chr(187)] = 445
	me.cw[gopdf.Chr(188)] = 414
	me.cw[gopdf.Chr(189)] = 414
	me.cw[gopdf.Chr(190)] = 482
	me.cw[gopdf.Chr(191)] = 482
	me.cw[gopdf.Chr(192)] = 425
	me.cw[gopdf.Chr(193)] = 417
	me.cw[gopdf.Chr(194)] = 386
	me.cw[gopdf.Chr(195)] = 329
	me.cw[gopdf.Chr(196)] = 391
	me.cw[gopdf.Chr(197)] = 399
	me.cw[gopdf.Chr(198)] = 425
	me.cw[gopdf.Chr(199)] = 379
	me.cw[gopdf.Chr(200)] = 403
	me.cw[gopdf.Chr(201)] = 461
	me.cw[gopdf.Chr(202)] = 399
	me.cw[gopdf.Chr(203)] = 438
	me.cw[gopdf.Chr(204)] = 468
	me.cw[gopdf.Chr(205)] = 394
	me.cw[gopdf.Chr(206)] = 385
	me.cw[gopdf.Chr(207)] = 380
	me.cw[gopdf.Chr(208)] = 340
	me.cw[gopdf.Chr(209)] = 0
	me.cw[gopdf.Chr(210)] = 336
	me.cw[gopdf.Chr(211)] = 336
	me.cw[gopdf.Chr(212)] = 0
	me.cw[gopdf.Chr(213)] = 0
	me.cw[gopdf.Chr(214)] = 0
	me.cw[gopdf.Chr(215)] = 0
	me.cw[gopdf.Chr(216)] = 0
	me.cw[gopdf.Chr(217)] = 0
	me.cw[gopdf.Chr(218)] = 0
	me.cw[gopdf.Chr(219)] = 698
	me.cw[gopdf.Chr(220)] = 698
	me.cw[gopdf.Chr(221)] = 698
	me.cw[gopdf.Chr(222)] = 698
	me.cw[gopdf.Chr(223)] = 389
	me.cw[gopdf.Chr(224)] = 208
	me.cw[gopdf.Chr(225)] = 395
	me.cw[gopdf.Chr(226)] = 252
	me.cw[gopdf.Chr(227)] = 269
	me.cw[gopdf.Chr(228)] = 252
	me.cw[gopdf.Chr(229)] = 209
	me.cw[gopdf.Chr(230)] = 437
	me.cw[gopdf.Chr(231)] = 0
	me.cw[gopdf.Chr(232)] = 0
	me.cw[gopdf.Chr(233)] = 0
	me.cw[gopdf.Chr(234)] = 0
	me.cw[gopdf.Chr(235)] = 0
	me.cw[gopdf.Chr(236)] = 0
	me.cw[gopdf.Chr(237)] = 0
	me.cw[gopdf.Chr(238)] = 0
	me.cw[gopdf.Chr(239)] = 450
	me.cw[gopdf.Chr(240)] = 456
	me.cw[gopdf.Chr(241)] = 456
	me.cw[gopdf.Chr(242)] = 456
	me.cw[gopdf.Chr(243)] = 456
	me.cw[gopdf.Chr(244)] = 456
	me.cw[gopdf.Chr(245)] = 456
	me.cw[gopdf.Chr(246)] = 456
	me.cw[gopdf.Chr(247)] = 456
	me.cw[gopdf.Chr(248)] = 456
	me.cw[gopdf.Chr(249)] = 456
	me.cw[gopdf.Chr(250)] = 525
	me.cw[gopdf.Chr(251)] = 697
	me.cw[gopdf.Chr(252)] = 698
	me.cw[gopdf.Chr(253)] = 698
	me.cw[gopdf.Chr(254)] = 698
	me.cw[gopdf.Chr(255)] = 698
	me.up = -35
	me.ut = 30
	me.fonttype = "TrueType"
	me.name = "THSarabunNew-Bold"
	me.enc = "cp874"
	me.diff = "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef /.notdef"
	me.desc = make([]gopdf.FontDescItem, 8)
	me.desc[0] = gopdf.FontDescItem{Key: "Ascent", Val: "850"}
	me.desc[1] = gopdf.FontDescItem{Key: "Descent", Val: "-250"}
	me.desc[2] = gopdf.FontDescItem{Key: "CapHeight", Val: "476"}
	me.desc[3] = gopdf.FontDescItem{Key: "Flags", Val: "32"}
	me.desc[4] = gopdf.FontDescItem{Key: "FontBBox", Val: "[-466 -457 947 844]"}
	me.desc[5] = gopdf.FontDescItem{Key: "ItalicAngle", Val: "0"}
	me.desc[6] = gopdf.FontDescItem{Key: "StemV", Val: "120"}
	me.desc[7] = gopdf.FontDescItem{Key: "MissingWidth", Val: "698"}
}
func (me *THSarabunNewBold) GetType() string {
	return me.fonttype
}
func (me *THSarabunNewBold) GetName() string {
	return me.name
}
func (me *THSarabunNewBold) GetDesc() []gopdf.FontDescItem {
	return me.desc
}
func (me *THSarabunNewBold) GetUp() int {
	return me.up
}
func (me *THSarabunNewBold) GetUt() int {
	return me.ut
}
func (me *THSarabunNewBold) GetCw() gopdf.FontCw {
	return me.cw
}
func (me *THSarabunNewBold) GetEnc() string {
	return me.enc
}
func (me *THSarabunNewBold) GetDiff() string {
	return me.diff
}
func (me *THSarabunNewBold) GetOriginalsize() int {
	return 98764
}
func (me *THSarabunNewBold) SetFamily(family string) {
	me.family = family
}
func (me *THSarabunNewBold) GetFamily() string {
	return me.family
}
