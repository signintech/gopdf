package core

import (
	"fmt"
	"os"
)

//Parsekern parse kerning table  https://www.microsoft.com/typography/otspec/kern.htm
func (t *TTFParser) Parsekern(fd *os.File) error {

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

	i := uint64(0)
	for i < nTables {
		var subtable KernSubTable
		err = t.parsekernSubTable(fd, &subtable)
		if err != nil {
			return err
		}
		t.kern.Subtables = append(t.kern.Subtables, subtable)
		i++
	}

	return nil
}

func (t *TTFParser) parsekernSubTable(fd *os.File, subtable *KernSubTable) error {

	ver, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	length, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	coverage, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	format := subtable.Coverage & 0xf0
	//debug
	fmt.Printf("format = %d\n", format)
	if format == 0 {
		t.parsekernSubTableFormat0(fd, subtable)
	} else {
		//not support other format yet
		return fmt.Errorf("not support kerning format %d", format)
	}

	subtable.Version = ver
	subtable.Length = length
	subtable.Coverage = coverage
	return nil
}

func (t *TTFParser) parsekernSubTableFormat0(fd *os.File, subtable *KernSubTable) error {
	nPairs, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.Skip(fd, 2+2+2)

	i := uint64(0)
	for i < nPairs {
		left, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		right, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		value, err := t.ReadShort(fd)
		if err != nil {
			return err
		}

		//debug
		_ = fmt.Sprintf("nPairs %d left %d right %d value %d\n", nPairs, left, right, value)
		i++
	}
	return nil
}
