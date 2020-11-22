package core

//FeatureRecord FeatureRecord
type FeatureRecord struct {
	featureTag    string //4-byte feature identification tag
	featureOffset int64  //Offset to Feature table, from beginning of FeatureList
	featureTable  FeatureTable
}

//FeatureTable A Feature table defines a feature with one or more lookups. The client uses the lookups to substitute or position glyphs.
type FeatureTable struct {
	featureParamsOffset int64  //Offset from start of Feature table to FeatureParams table, if defined for the feature and present, else NULL
	lookupIndexCount    uint   //Number of LookupList indices for this feature
	lookupListIndices   []uint //	Array of indices into the LookupList — zero-based (first lookup is LookupListIndex = 0)
}

//GSUBParseFeatureListResult result of parseFeatureList(...)
type GSUBParseFeatureListResult struct {
	featureRecords []FeatureRecord
}
