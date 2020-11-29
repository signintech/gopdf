package core

import (
	"fmt"
	"sort"
)

func (t *TTFParser) GSUBProcessGlyphs(glyphindexs []uint) ([]uint, error) {
	fmt.Printf("warn: fake mlym")
	err := t.gsubPreprocessGlyphs("mlym") //JUST for test
	if err != nil {
		return nil, err
	}

	return glyphindexs, nil
}

func (t *TTFParser) gsubPreprocessGlyphs(script string) error {

	featureRecords, err := findFeatureRecordsForScript(script, t.gsubScriptList, t.gsubFeatureList)
	if err != nil {
		return err
	}

	var lookupIndexes []preprocessLookupIndex
	for _, featureRecord := range featureRecords {
		for _, lookupListIndex := range featureRecord.featureTable.lookupListIndices {
			lookupIndexes = append(lookupIndexes, preprocessLookupIndex{
				featureTag:      featureRecord.featureTag,
				lookupListIndex: lookupListIndex,
			})
		}
	}

	sort.Slice(lookupIndexes, func(i, j int) bool {
		if lookupIndexes[i].lookupListIndex <= lookupIndexes[j].lookupListIndex {
			return true
		}
		return false
	})

	lookups := t.gsubLookups
	for _, lookupIndex := range lookupIndexes {
		_ = lookupIndex
	}
	_ = lookups
	return nil
}

func findFeatureRecordsForScript(script string, scriptList GSUBParseScriptListResult, featureList GSUBParseFeatureListResult) ([]FeatureRecord, error) {
	var featureRecords []FeatureRecord
	for scriptTag, scriptItem := range scriptList.scripts {
		if scriptTag == script {
			var indexs = scriptItem.defaultLangSys.featureIndices
			sort.Slice(indexs, func(i, j int) bool {
				if scriptItem.defaultLangSys.featureIndices[i] <= scriptItem.defaultLangSys.featureIndices[j] {
					return true
				}
				return false
			})
			for _, index := range indexs {
				featureRecords = append(featureRecords, featureList.featureRecords[index])
			}
			break
		}
	}
	return featureRecords, nil
}

type preprocessLookupIndex struct {
	featureTag      string
	lookupListIndex uint
}
