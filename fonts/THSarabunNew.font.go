package fonts

import (
	"github.com/signintech/gopdf"
)

type THSarabunNew struct {
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

func (me *THSarabunNew) Init() {
	me.cw = make(gopdf.FontCw)
	me.cw[gopdf.Chr(0)] = 692
	me.cw[gopdf.Chr(1)] = 692
	me.cw[gopdf.Chr(2)] = 692
	me.cw[gopdf.Chr(3)] = 692
	me.cw[gopdf.Chr(4)] = 692
	me.cw[gopdf.Chr(5)] = 692
	me.cw[gopdf.Chr(6)] = 692
	me.cw[gopdf.Chr(7)] = 692
	me.cw[gopdf.Chr(8)] = 692
	me.cw[gopdf.Chr(9)] = 692
	me.cw[gopdf.Chr(10)] = 692
	me.cw[gopdf.Chr(11)] = 692
	me.cw[gopdf.Chr(12)] = 692
	me.cw[gopdf.Chr(13)] = 692
	me.cw[gopdf.Chr(14)] = 692
	me.cw[gopdf.Chr(15)] = 692
	me.cw[gopdf.Chr(16)] = 692
	me.cw[gopdf.Chr(17)] = 692
	me.cw[gopdf.Chr(18)] = 692
	me.cw[gopdf.Chr(19)] = 692
	me.cw[gopdf.Chr(20)] = 692
	me.cw[gopdf.Chr(21)] = 692
	me.cw[gopdf.Chr(22)] = 692
	me.cw[gopdf.Chr(23)] = 692
	me.cw[gopdf.Chr(24)] = 692
	me.cw[gopdf.Chr(25)] = 692
	me.cw[gopdf.Chr(26)] = 692
	me.cw[gopdf.Chr(27)] = 692
	me.cw[gopdf.Chr(28)] = 692
	me.cw[gopdf.Chr(29)] = 692
	me.cw[gopdf.Chr(30)] = 692
	me.cw[gopdf.Chr(31)] = 692
	me.cw[gopdf.ToByte(" ")] = 216
	me.cw[gopdf.ToByte("!")] = 147
	me.cw[gopdf.ToByte("\"")] = 208
	me.cw[gopdf.ToByte("#")] = 403
	me.cw[gopdf.ToByte("$")] = 361
	me.cw[gopdf.ToByte("%")] = 585
	me.cw[gopdf.ToByte("&")] = 423
	me.cw[gopdf.ToByte("'")] = 120
	me.cw[gopdf.ToByte("(")] = 190
	me.cw[gopdf.ToByte(")")] = 190
	me.cw[gopdf.ToByte("*")] = 285
	me.cw[gopdf.ToByte("+")] = 411
	me.cw[gopdf.ToByte(",")] = 162
	me.cw[gopdf.ToByte("-")] = 216
	me.cw[gopdf.ToByte(".")] = 162
	me.cw[gopdf.ToByte("/")] = 270
	me.cw[gopdf.ToByte("0")] = 362
	me.cw[gopdf.ToByte("1")] = 362
	me.cw[gopdf.ToByte("2")] = 362
	me.cw[gopdf.ToByte("3")] = 362
	me.cw[gopdf.ToByte("4")] = 362
	me.cw[gopdf.ToByte("5")] = 362
	me.cw[gopdf.ToByte("6")] = 362
	me.cw[gopdf.ToByte("7")] = 362
	me.cw[gopdf.ToByte("8")] = 362
	me.cw[gopdf.ToByte("9")] = 362
	me.cw[gopdf.ToByte(":")] = 162
	me.cw[gopdf.ToByte(";")] = 162
	me.cw[gopdf.ToByte("<")] = 411
	me.cw[gopdf.ToByte("=")] = 411
	me.cw[gopdf.ToByte(">")] = 411
	me.cw[gopdf.ToByte("?")] = 283
	me.cw[gopdf.ToByte("@")] = 536
	me.cw[gopdf.ToByte("A")] = 400
	me.cw[gopdf.ToByte("B")] = 378
	me.cw[gopdf.ToByte("C")] = 406
	me.cw[gopdf.ToByte("D")] = 431
	me.cw[gopdf.ToByte("E")] = 351
	me.cw[gopdf.ToByte("F")] = 351
	me.cw[gopdf.ToByte("G")] = 425
	me.cw[gopdf.ToByte("H")] = 441
	me.cw[gopdf.ToByte("I")] = 147
	me.cw[gopdf.ToByte("J")] = 264
	me.cw[gopdf.ToByte("K")] = 376
	me.cw[gopdf.ToByte("L")] = 353
	me.cw[gopdf.ToByte("M")] = 548
	me.cw[gopdf.ToByte("N")] = 441
	me.cw[gopdf.ToByte("O")] = 486
	me.cw[gopdf.ToByte("P")] = 378
	me.cw[gopdf.ToByte("Q")] = 487
	me.cw[gopdf.ToByte("R")] = 379
	me.cw[gopdf.ToByte("S")] = 352
	me.cw[gopdf.ToByte("T")] = 379
	me.cw[gopdf.ToByte("U")] = 466
	me.cw[gopdf.ToByte("V")] = 390
	me.cw[gopdf.ToByte("W")] = 588
	me.cw[gopdf.ToByte("X")] = 418
	me.cw[gopdf.ToByte("Y")] = 366
	me.cw[gopdf.ToByte("Z")] = 424
	me.cw[gopdf.ToByte("[")] = 196
	me.cw[gopdf.ToByte("\\")] = 262
	me.cw[gopdf.ToByte("]")] = 196
	me.cw[gopdf.ToByte("^")] = 412
	me.cw[gopdf.ToByte("_")] = 352
	me.cw[gopdf.ToByte("`")] = 204
	me.cw[gopdf.ToByte("a")] = 344
	me.cw[gopdf.ToByte("b")] = 401
	me.cw[gopdf.ToByte("c")] = 331
	me.cw[gopdf.ToByte("d")] = 401
	me.cw[gopdf.ToByte("e")] = 374
	me.cw[gopdf.ToByte("f")] = 206
	me.cw[gopdf.ToByte("g")] = 311
	me.cw[gopdf.ToByte("h")] = 390
	me.cw[gopdf.ToByte("i")] = 143
	me.cw[gopdf.ToByte("j")] = 155
	me.cw[gopdf.ToByte("k")] = 316
	me.cw[gopdf.ToByte("l")] = 200
	me.cw[gopdf.ToByte("m")] = 601
	me.cw[gopdf.ToByte("n")] = 390
	me.cw[gopdf.ToByte("o")] = 398
	me.cw[gopdf.ToByte("p")] = 401
	me.cw[gopdf.ToByte("q")] = 401
	me.cw[gopdf.ToByte("r")] = 217
	me.cw[gopdf.ToByte("s")] = 282
	me.cw[gopdf.ToByte("t")] = 238
	me.cw[gopdf.ToByte("u")] = 390
	me.cw[gopdf.ToByte("v")] = 341
	me.cw[gopdf.ToByte("w")] = 507
	me.cw[gopdf.ToByte("x")] = 318
	me.cw[gopdf.ToByte("y")] = 337
	me.cw[gopdf.ToByte("z")] = 321
	me.cw[gopdf.ToByte("{")] = 208
	me.cw[gopdf.ToByte("|")] = 153
	me.cw[gopdf.ToByte("}")] = 208
	me.cw[gopdf.ToByte("~")] = 416
	me.cw[gopdf.Chr(127)] = 692
	me.cw[gopdf.Chr(128)] = 406
	me.cw[gopdf.Chr(129)] = 692
	me.cw[gopdf.Chr(130)] = 692
	me.cw[gopdf.Chr(131)] = 692
	me.cw[gopdf.Chr(132)] = 692
	me.cw[gopdf.Chr(133)] = 479
	me.cw[gopdf.Chr(134)] = 692
	me.cw[gopdf.Chr(135)] = 692
	me.cw[gopdf.Chr(136)] = 692
	me.cw[gopdf.Chr(137)] = 692
	me.cw[gopdf.Chr(138)] = 692
	me.cw[gopdf.Chr(139)] = 692
	me.cw[gopdf.Chr(140)] = 692
	me.cw[gopdf.Chr(141)] = 692
	me.cw[gopdf.Chr(142)] = 692
	me.cw[gopdf.Chr(143)] = 692
	me.cw[gopdf.Chr(144)] = 692
	me.cw[gopdf.Chr(145)] = 247
	me.cw[gopdf.Chr(146)] = 247
	me.cw[gopdf.Chr(147)] = 370
	me.cw[gopdf.Chr(148)] = 370
	me.cw[gopdf.Chr(149)] = 216
	me.cw[gopdf.Chr(150)] = 360
	me.cw[gopdf.Chr(151)] = 720
	me.cw[gopdf.Chr(152)] = 692
	me.cw[gopdf.Chr(153)] = 692
	me.cw[gopdf.Chr(154)] = 692
	me.cw[gopdf.Chr(155)] = 692
	me.cw[gopdf.Chr(156)] = 692
	me.cw[gopdf.Chr(157)] = 692
	me.cw[gopdf.Chr(158)] = 692
	me.cw[gopdf.Chr(159)] = 692
	me.cw[gopdf.Chr(160)] = 216
	me.cw[gopdf.Chr(161)] = 386
	me.cw[gopdf.Chr(162)] = 378
	me.cw[gopdf.Chr(163)] = 382
	me.cw[gopdf.Chr(164)] = 393
	me.cw[gopdf.Chr(165)] = 393
	me.cw[gopdf.Chr(166)] = 408
	me.cw[gopdf.Chr(167)] = 294
	me.cw[gopdf.Chr(168)] = 367
	me.cw[gopdf.Chr(169)] = 377
	me.cw[gopdf.Chr(170)] = 380
	me.cw[gopdf.Chr(171)] = 384
	me.cw[gopdf.Chr(172)] = 519
	me.cw[gopdf.Chr(173)] = 519
	me.cw[gopdf.Chr(174)] = 425
	me.cw[gopdf.Chr(175)] = 425
	me.cw[gopdf.Chr(176)] = 343
	me.cw[gopdf.Chr(177)] = 461
	me.cw[gopdf.Chr(178)] = 532
	me.cw[gopdf.Chr(179)] = 543
	me.cw[gopdf.Chr(180)] = 391
	me.cw[gopdf.Chr(181)] = 391
	me.cw[gopdf.Chr(182)] = 378
	me.cw[gopdf.Chr(183)] = 430
	me.cw[gopdf.Chr(184)] = 335
	me.cw[gopdf.Chr(185)] = 420
	me.cw[gopdf.Chr(186)] = 428
	me.cw[gopdf.Chr(187)] = 428
	me.cw[gopdf.Chr(188)] = 381
	me.cw[gopdf.Chr(189)] = 381
	me.cw[gopdf.Chr(190)] = 447
	me.cw[gopdf.Chr(191)] = 447
	me.cw[gopdf.Chr(192)] = 425
	me.cw[gopdf.Chr(193)] = 400
	me.cw[gopdf.Chr(194)] = 375
	me.cw[gopdf.Chr(195)] = 322
	me.cw[gopdf.Chr(196)] = 378
	me.cw[gopdf.Chr(197)] = 381
	me.cw[gopdf.Chr(198)] = 425
	me.cw[gopdf.Chr(199)] = 335
	me.cw[gopdf.Chr(200)] = 393
	me.cw[gopdf.Chr(201)] = 438
	me.cw[gopdf.Chr(202)] = 381
	me.cw[gopdf.Chr(203)] = 427
	me.cw[gopdf.Chr(204)] = 454
	me.cw[gopdf.Chr(205)] = 387
	me.cw[gopdf.Chr(206)] = 372
	me.cw[gopdf.Chr(207)] = 391
	me.cw[gopdf.Chr(208)] = 357
	me.cw[gopdf.Chr(209)] = 0
	me.cw[gopdf.Chr(210)] = 316
	me.cw[gopdf.Chr(211)] = 316
	me.cw[gopdf.Chr(212)] = 0
	me.cw[gopdf.Chr(213)] = 0
	me.cw[gopdf.Chr(214)] = 0
	me.cw[gopdf.Chr(215)] = 0
	me.cw[gopdf.Chr(216)] = 0
	me.cw[gopdf.Chr(217)] = 0
	me.cw[gopdf.Chr(218)] = 0
	me.cw[gopdf.Chr(219)] = 692
	me.cw[gopdf.Chr(220)] = 692
	me.cw[gopdf.Chr(221)] = 692
	me.cw[gopdf.Chr(222)] = 692
	me.cw[gopdf.Chr(223)] = 411
	me.cw[gopdf.Chr(224)] = 203
	me.cw[gopdf.Chr(225)] = 377
	me.cw[gopdf.Chr(226)] = 237
	me.cw[gopdf.Chr(227)] = 242
	me.cw[gopdf.Chr(228)] = 244
	me.cw[gopdf.Chr(229)] = 205
	me.cw[gopdf.Chr(230)] = 399
	me.cw[gopdf.Chr(231)] = 0
	me.cw[gopdf.Chr(232)] = 0
	me.cw[gopdf.Chr(233)] = 0
	me.cw[gopdf.Chr(234)] = 0
	me.cw[gopdf.Chr(235)] = 0
	me.cw[gopdf.Chr(236)] = 0
	me.cw[gopdf.Chr(237)] = 0
	me.cw[gopdf.Chr(238)] = 0
	me.cw[gopdf.Chr(239)] = 450
	me.cw[gopdf.Chr(240)] = 449
	me.cw[gopdf.Chr(241)] = 449
	me.cw[gopdf.Chr(242)] = 449
	me.cw[gopdf.Chr(243)] = 449
	me.cw[gopdf.Chr(244)] = 449
	me.cw[gopdf.Chr(245)] = 449
	me.cw[gopdf.Chr(246)] = 449
	me.cw[gopdf.Chr(247)] = 449
	me.cw[gopdf.Chr(248)] = 449
	me.cw[gopdf.Chr(249)] = 449
	me.cw[gopdf.Chr(250)] = 522
	me.cw[gopdf.Chr(251)] = 697
	me.cw[gopdf.Chr(252)] = 692
	me.cw[gopdf.Chr(253)] = 692
	me.cw[gopdf.Chr(254)] = 692
	me.cw[gopdf.Chr(255)] = 692
	me.up = -35
	me.ut = 30
	me.fonttype = "TrueType"
	me.name = "THSarabunNew"
	me.enc = "cp874"
	me.diff = "130 /.notdef /.notdef /.notdef 134 /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef /.notdef 142 /.notdef 152 /.notdef /.notdef /.notdef /.notdef /.notdef 158 /.notdef /.notdef 161 /kokaithai /khokhaithai /khokhuatthai /khokhwaithai /khokhonthai /khorakhangthai /ngonguthai /chochanthai /chochingthai /chochangthai /sosothai /chochoethai /yoyingthai /dochadathai /topatakthai /thothanthai /thonangmonthothai /thophuthaothai /nonenthai /dodekthai /totaothai /thothungthai /thothahanthai /thothongthai /nonuthai /bobaimaithai /poplathai /phophungthai /fofathai /phophanthai /fofanthai /phosamphaothai /momathai /yoyakthai /roruathai /ruthai /lolingthai /luthai /wowaenthai /sosalathai /sorusithai /sosuathai /hohipthai /lochulathai /oangthai /honokhukthai /paiyannoithai /saraathai /maihanakatthai /saraaathai /saraamthai /saraithai /saraiithai /sarauethai /saraueethai /sarauthai /sarauuthai /phinthuthai /.notdef /.notdef /.notdef /.notdef /bahtthai /saraethai /saraaethai /saraothai /saraaimaimuanthai /saraaimaimalaithai /lakkhangyaothai /maiyamokthai /maitaikhuthai /maiekthai /maithothai /maitrithai /maichattawathai /thanthakhatthai /nikhahitthai /yamakkanthai /fongmanthai /zerothai /onethai /twothai /threethai /fourthai /fivethai /sixthai /seventhai /eightthai /ninethai /angkhankhuthai /khomutthai /.notdef /.notdef /.notdef /.notdef"
	me.desc = make([]gopdf.FontDescItem, 8)
	me.desc[0] = gopdf.FontDescItem{Key: "Ascent", Val: "850"}
	me.desc[1] = gopdf.FontDescItem{Key: "Descent", Val: "-250"}
	me.desc[2] = gopdf.FontDescItem{Key: "CapHeight", Val: "476"}
	me.desc[3] = gopdf.FontDescItem{Key: "Flags", Val: "32"}
	me.desc[4] = gopdf.FontDescItem{Key: "FontBBox", Val: "[-427 -421 947 836]"}
	me.desc[5] = gopdf.FontDescItem{Key: "ItalicAngle", Val: "0"}
	me.desc[6] = gopdf.FontDescItem{Key: "StemV", Val: "70"}
	me.desc[7] = gopdf.FontDescItem{Key: "MissingWidth", Val: "692"}
}
func (me *THSarabunNew) GetType() string {
	return me.fonttype
}
func (me *THSarabunNew) GetName() string {
	return me.name
}
func (me *THSarabunNew) GetDesc() []gopdf.FontDescItem {
	return me.desc
}
func (me *THSarabunNew) GetUp() int {
	return me.up
}
func (me *THSarabunNew) GetUt() int {
	return me.ut
}
func (me *THSarabunNew) GetCw() gopdf.FontCw {
	return me.cw
}
func (me *THSarabunNew) GetEnc() string {
	return me.enc
}
func (me *THSarabunNew) GetDiff() string {
	return me.diff
}
func (me *THSarabunNew) GetOriginalsize() int {
	return 98764
}
func (me *THSarabunNew) SetFamily(family string) {
	me.family = family
}
func (me *THSarabunNew) GetFamily() string {
	return me.family
}
