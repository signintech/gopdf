package core

import "bytes"

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

	return nil
}
