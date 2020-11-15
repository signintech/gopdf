package core

import (
	"bytes"
)

func (t *TTFParser) parseScriptList(fd *bytes.Reader, scriptListOffset int64) (GSUBParseScriptListResult, error) {

	_, err := fd.Seek(scriptListOffset, 0)
	if err != nil {
		return GSUBParseScriptListResult{}, err
	}

	scriptCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBParseScriptListResult{}, err
	}

	var scriptRecords []GSUBScriptRecord
	for i := uint(0); i < scriptCount; i++ {
		scriptTag, err := t.Read(fd, 4)
		if err != nil {
			return GSUBParseScriptListResult{}, err
		}
		scriptOffset, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBParseScriptListResult{}, err
		}

		scriptRecords = append(scriptRecords, GSUBScriptRecord{
			scriptTag:    scriptTag,
			scriptOffset: scriptListOffset + int64(scriptOffset),
		})
	}

	//parse ScriptTable
	scriptTables := make(map[string]GSUBScriptTable)
	for _, scriptRecord := range scriptRecords {
		_, err := fd.Seek(scriptRecord.scriptOffset, 0)
		if err != nil {
			return GSUBParseScriptListResult{}, err
		}
		defaultLangSys, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBParseScriptListResult{}, err
		}
		langSysCount, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBParseScriptListResult{}, err
		}

		//parse LangSysRecord
		var langSysRecords []LangSysRecord
		for j := uint(0); j < langSysCount; j++ {
			langSysTag, err := t.Read(fd, 4)
			if err != nil {
				return GSUBParseScriptListResult{}, err
			}
			langSysOffset, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBParseScriptListResult{}, err
			}

			langSysRecords = append(langSysRecords, LangSysRecord{
				langSysTag:    langSysTag,
				langSysOffset: scriptRecord.scriptOffset + int64(langSysOffset),
			})
		}

		defaultLangSysOffset := int64(0)
		if defaultLangSys > 0 {
			defaultLangSysOffset = int64(defaultLangSys) + scriptRecord.scriptOffset
		}

		scriptTable := GSUBScriptTable{
			defaultLangSys: defaultLangSysOffset,
			langSysCount:   langSysCount,
			langSysRecords: langSysRecords,
		}
		scriptTables[scriptRecord.scriptTagString()] = scriptTable
	}

	var result GSUBParseScriptListResult
	for scriptTag, scriptTable := range scriptTables {
		mm := scriptTable.convertToMap()
		var languageSystemTable LanguageSystemTable
		for tag, langsystableOffset := range mm {

			_, err := fd.Seek(langsystableOffset, 0)
			if err != nil {
				return GSUBParseScriptListResult{}, err
			}

			err = t.Skip(fd, 2) //lookupOrder	= NULL (reserved for an offset to a reordering table)
			if err != nil {
				return GSUBParseScriptListResult{}, err
			}

			requiredFeatureIndex, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBParseScriptListResult{}, err
			}

			featureIndexCount, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBParseScriptListResult{}, err
			}

			var featureIndices []uint
			for j := uint(0); j < featureIndexCount; j++ {
				featureIndice, err := t.ReadUShort(fd)
				if err != nil {
					return GSUBParseScriptListResult{}, err
				}
				featureIndices = append(featureIndices, featureIndice)
			}
			//set languageSystemTables
			languageSystemTable.requiredFeatureIndex = requiredFeatureIndex
			languageSystemTable.featureIndexCount = featureIndexCount
			languageSystemTable.featureIndices = featureIndices
			//end set languageSystemTables
			result.append(scriptTag, tag, languageSystemTable)
		}
	}

	return result, nil
}
