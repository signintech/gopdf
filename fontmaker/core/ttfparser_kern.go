package core

import (
	"bytes"
	"fmt"
)

//Parsekern parse kerning table  https://www.microsoft.com/typography/otspec/kern.htm
func (t *TTFParser) Parsekern(fd *bytes.Reader) error {

	t.kern = nil //clear
	err := t.Seek(fd, "kern")
	if err == ErrTableNotFound {
		return nil
	} else if err != nil {
		return err
	}

	t.kern = new(KernTable) //init

	version, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.kern.Version = version

	nTables, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.kern.NTables = nTables

	i := uint(0)
	for i < nTables {
		err = t.parsekernSubTable(fd)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}

func (t *TTFParser) parsekernSubTable(fd *bytes.Reader) error {

	t.Skip(fd, 2+2) //skip version and length

	coverage, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	format := coverage & 0xf0
	//fmt.Printf("format = %d\n", format) //debug
	t.kern.Kerning = make(KernMap) //init
	if format == 0 {
		t.parsekernSubTableFormat0(fd)
	} else {
		//not support other format yet
		return fmt.Errorf("not support kerning format %d", format)
	}

	return nil
}

func (t *TTFParser) parsekernSubTableFormat0(fd *bytes.Reader) error {
	nPairs, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.Skip(fd, 2+2+2) //skip searchRange , entrySelector , rangeShift

	i := uint(0)
	for i < nPairs {
		left, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		right, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		value, err := t.ReadShortInt16(fd)
		if err != nil {
			return err
		}

		if _, ok := t.kern.Kerning[left]; !ok {
			kval := make(KernValue)
			kval[right] = value
			t.kern.Kerning[left] = kval
		} else {
			(t.kern.Kerning[left])[right] = value
		}
		//_ = fmt.Sprintf("nPairs %d left %d right %d value %d\n", nPairs, left, right, value) //debug
		i++
	}
	return nil
}
