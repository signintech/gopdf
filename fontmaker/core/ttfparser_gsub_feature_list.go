package core

import (
	"bytes"
	"fmt"
)

func (t *TTFParser) parseFeatureList(fd *bytes.Reader, featureListOffset int64) error {

	_, err := fd.Seek(featureListOffset, 0)
	if err != nil {
		return err
	}

	//Number of FeatureRecords in this table
	featureCount, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	var featureRecords []FeatureRecord
	for i := uint(0); i < featureCount; i++ {
		featureTag, err := t.Read(fd, 4)
		if err != nil {
			return err
		}
		featureOffset, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		featureRecords = append(featureRecords, FeatureRecord{
			featureTag:    featureTag,
			featureOffset: featureListOffset + int64(featureOffset),
		})
	}

	var featureTables []FeatureTable
	for _, fr := range featureRecords {
		_, err := fd.Seek(fr.featureOffset, 0)
		if err != nil {
			return err
		}
		featureParamsOffset, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		lookupIndexCount, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		var lookupListIndices []uint
		for j := uint(0); j < lookupIndexCount; j++ {
			lookupListIndex, err := t.ReadUShort(fd)
			if err != nil {
				return err
			}
			lookupListIndices = append(lookupListIndices, lookupListIndex)
		}

		featureTables = append(featureTables, FeatureTable{
			featureParamsOffset: fr.featureOffset + int64(featureParamsOffset),
			lookupIndexCount:    lookupIndexCount,
			lookupListIndices:   lookupListIndices,
		})
		//Offset16	featureParamsOffset
	}

	fmt.Printf("featureCount = %d\n", featureCount)

	return nil
}
