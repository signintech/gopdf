package core

import (
	"bytes"
)

func (t *TTFParser) parseFeatureList(fd *bytes.Reader, featureListOffset int64, parseScriptListResult GSUBParseScriptListResult) (GSUBParseFeatureListResult, error) {

	_, err := fd.Seek(featureListOffset, 0)
	if err != nil {
		return GSUBParseFeatureListResult{}, err
	}

	featureCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBParseFeatureListResult{}, err
	}

	var featureRecords []FeatureRecord
	for i := uint(0); i < featureCount; i++ {
		featureTag, err := t.Read(fd, 4)
		if err != nil {
			return GSUBParseFeatureListResult{}, err
		}
		featureOffset, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBParseFeatureListResult{}, err
		}
		featureRecords = append(featureRecords, FeatureRecord{
			featureTag:    string(featureTag),
			featureOffset: featureListOffset + int64(featureOffset),
		})
	}

	//var featureTables []FeatureTable
	for i, fr := range featureRecords {
		_, err := fd.Seek(fr.featureOffset, 0)
		if err != nil {
			return GSUBParseFeatureListResult{}, err
		}
		featureParamsOffset, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBParseFeatureListResult{}, err
		}
		lookupIndexCount, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBParseFeatureListResult{}, err
		}
		var lookupListIndices []uint
		for j := uint(0); j < lookupIndexCount; j++ {
			lookupListIndex, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBParseFeatureListResult{}, err
			}
			lookupListIndices = append(lookupListIndices, lookupListIndex)
		}

		featureTable := FeatureTable{
			featureParamsOffset: fr.featureOffset + int64(featureParamsOffset),
			lookupIndexCount:    lookupIndexCount,
			lookupListIndices:   lookupListIndices,
		}
		featureRecords[i].featureTable = featureTable
	}

	return GSUBParseFeatureListResult{
		featureRecords: featureRecords,
	}, nil
}
