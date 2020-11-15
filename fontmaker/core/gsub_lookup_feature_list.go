package core

//FeatureRecord FeatureRecord
type FeatureRecord struct {
	featureTag    []byte //4-byte feature identification tag
	featureOffset int64  //Offset to Feature table, from beginning of FeatureList
}
