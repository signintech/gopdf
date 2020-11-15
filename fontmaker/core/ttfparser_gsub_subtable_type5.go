package core

/*
NOT SUPPORT YET
//5.1 Context Substitution Format 1: Simple Glyph Contexts
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType5Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	GSUBLookupSubTableType5Format1,
	error,
) {
	var result GSUBLookupSubTableType5Format1
	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType5Format1{}, err
	}
	result.coverageOffset = int64(coverageOffset) + offset //set result

	subRuleSetCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType5Format1{}, err
	}
	result.subRuleSetCount = subRuleSetCount //set result

	var subRuleSetOffsets []int64
	for i := uint(0); i < result.subRuleSetCount; i++ {
		subRuleSetOffset, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType5Format1{}, err
		}
		subRuleSetOffsets = append(subRuleSetOffsets, offset+int64(subRuleSetOffset))
	}
	result.subRuleSetOffsets = subRuleSetOffsets //set result

	var subRuleSets []SubRuleSetTable
	for _, subRuleSetOffset := range result.subRuleSetOffsets {
		_, err := fd.Seek(subRuleSetOffset, 0)
		if err != nil {
			return GSUBLookupSubTableType5Format1{}, err
		}
		subRuleCount, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType5Format1{}, err
		}
		var subRuleOffsets []int64
		for j := uint(0); j < subRuleCount; j++ {
			subRuleOffset, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBLookupSubTableType5Format1{}, err
			}
			subRuleOffsets = append(subRuleOffsets, subRuleSetOffset+int64(subRuleOffset))
		}

		subRuleSets = append(subRuleSets, SubRuleSetTable{
			subRuleCount:   subRuleCount,
			subRuleOffsets: subRuleOffsets,
		})
	}

	for i, subRuleSet := range subRuleSets {
		var subRules []SubRuleTable
		for _, subRuleOffset := range subRuleSet.subRuleOffsets {
			_, err := fd.Seek(subRuleOffset, 0)
			if err != nil {
				return GSUBLookupSubTableType5Format1{}, err
			}
			glyphCount, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBLookupSubTableType5Format1{}, err
			}
			substitutionCount, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBLookupSubTableType5Format1{}, err
			}

			//parse inputSequence
			var inputSequence []uint
			for j := uint(1); j < glyphCount; j++ {
				input, err := t.ReadUShort(fd)
				if err != nil {
					return GSUBLookupSubTableType5Format1{}, err
				}
				inputSequence = append(inputSequence, input)
			}

			//parse substLookupRecords
			var substLookupRecords []SubstLookupRecord
			for j := uint(0); j < substitutionCount; j++ {
				glyphSequenceIndex, err := t.ReadUShort(fd)
				if err != nil {
					return GSUBLookupSubTableType5Format1{}, err
				}
				lookupListIndex, err := t.ReadUShort(fd)
				if err != nil {
					return GSUBLookupSubTableType5Format1{}, err
				}
				substLookupRecords = append(substLookupRecords, SubstLookupRecord{
					glyphSequenceIndex: glyphSequenceIndex,
					lookupListIndex:    lookupListIndex,
				})
			}

			subRules = append(subRules, SubRuleTable{
				glyphCount:         glyphCount,
				substitutionCount:  substitutionCount,
				inputSequence:      inputSequence,
				substLookupRecords: substLookupRecords,
			})

		} //for for
		subRuleSets[i].subRules = subRules //set subRuleSets
	} //end for

	result.subRuleSets = subRuleSets //set result
	return result, nil
}

func (t *TTFParser) processGSUBLookupListTableSubTableLookupType5Format1(
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType5Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableResult, error) {

	var result GSubLookupSubtableResult
	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableResult{}, err
	}
	glyphIDs := coverage.glyphIDs
	_ = glyphIDs

	for i, subRuleSet := range subtable.subRuleSets {
		firstGlyph := glyphIDs[i]
		fmt.Printf("firstGlyph = %d\n", firstGlyph)
		for _, subRule := range subRuleSet.subRules {
			for _, glyphID := range subRule.inputSequence {
				fmt.Printf("\tglyphID =%d\n", glyphID)
			}
		}
	}

	return result, nil
}
*/
