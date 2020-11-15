package core

//GSUBScriptListTable ScriptList table  https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#slTbl_sRec
type GSUBScriptListTable struct {
	scriptCount   uint               //Number of ScriptRecords
	scriptRecords []GSUBScriptRecord //Array of ScriptRecords, listed alphabetically by script tag
}

//GSUBScriptRecord Script Record https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#slTbl_sRec
type GSUBScriptRecord struct {
	scriptTag    []byte //4-byte script tag identifier
	scriptOffset int64  //Offset to Script table, from beginning of ScriptList
}

func (g GSUBScriptRecord) scriptTagString() string {
	return string(g.scriptTag)
}

//GSUBScriptTable Script Table https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#script-table-and-language-system-record
type GSUBScriptTable struct {
	defaultLangSys int64           //Offset to default LangSys table, from beginning of Script table — may be NULL
	langSysCount   uint            //Number of LangSysRecords for this script — excluding the default LangSys
	langSysRecords []LangSysRecord //Array of LangSysRecords, listed alphabetically by LangSys tag
}

func (g GSUBScriptTable) convertToMap() map[string]int64 {
	m := make(map[string]int64)
	m["DFLT"] = g.defaultLangSys
	for _, r := range g.langSysRecords {
		m[string(r.langSysTag)] = r.langSysOffset
	}
	return m
}

//LangSysRecord Language System Record https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#script-table-and-language-system-record
type LangSysRecord struct {
	langSysTag    []byte //4-byte LangSysTag identifier
	langSysOffset int64  //Offset to LangSys table, from beginning of Script table
}

//LanguageSystemTable https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#language-system-table
type LanguageSystemTable struct {
	//lookupOrder          int64  //= NULL (reserved for an offset to a reordering table)
	requiredFeatureIndex uint   //Index of a feature required for this language system; if no required features = 0xFFFF
	featureIndexCount    uint   //Number of feature index values for this language system — excludes the required feature
	featureIndices       []uint //Array of indices into the FeatureList, in arbitrary order
}

//GSUBParseScriptListResult result for parseScriptList
type GSUBParseScriptListResult struct {
	data map[string]map[string]([]uint)
}

func (g *GSUBParseScriptListResult) append(scriptTag, tag string, langSysTable LanguageSystemTable) {
	if g.data == nil {
		g.data = make(map[string]map[string]([]uint))
	}
	if g.data[scriptTag] == nil {
		g.data[scriptTag] = make(map[string]([]uint))
	}

	var indexs []uint
	if langSysTable.requiredFeatureIndex != 0xFFFF {
		indexs = append(indexs, langSysTable.requiredFeatureIndex)
	}
	indexs = append(indexs, langSysTable.featureIndices...)
	g.data[scriptTag][tag] = indexs
}
