package fontmaker

import (
	//"encoding/binary"
	//"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
)

type TTFParser struct {
	tables map[string]uint64
}

func (me *TTFParser) Parse(fontpath string) error {
	fmt.Printf("start parse\n")
	fd, err := os.Open(fontpath)
	if err != nil {
		return err
	}
	defer fd.Close()
	version, err := me.Read(fd, 4)
	if err != nil {
		return err
	}

	if !me.CompareBytes(version, []byte{0x00, 0x01, 0x00, 0x00}) {
		return errors.New("Unrecognized file (font) format")
	}

	i := uint64(0)
	numTables, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}
	me.Skip(fd, 3*2) //searchRange, entrySelector, rangeShift
	me.tables = make(map[string]uint64)
	for i < numTables {

		tag, err := me.Read(fd, 4)
		if err != nil {
			return err
		}

		err = me.Skip(fd, 4)
		if err != nil {
			return err
		}

		offset, err := me.ReadULong(fd)
		if err != nil {
			return err
		}

		err = me.Skip(fd, 4)
		if err != nil {
			return err
		}
		//fmt.Printf("%s\n", me.BytesToString(tag))
		me.tables[me.BytesToString(tag)] = offset
		i++
	}

	fmt.Printf("%+v\n", me.tables)

	me.ParseHead(fd)

	return nil
}

func (me *TTFParser) ParseHead(fd *os.File) error {
	me.Seek(fd, "head")
	me.Skip(fd, 3*4) // version, fontRevision, checkSumAdjustment
	magicNumber, err := me.ReadULong(fd)
	if err != nil {
		return err
	}
	fmt.Printf("%d", magicNumber)
	return nil
}

func (me *TTFParser) Seek(fd *os.File, tag string) error {
	val, ok := me.tables[tag]
	if !ok {
		return errors.New("me.tables not contain key=" + tag)
	}
	_, err := fd.Seek(int64(val), 1)
	if err != nil {
		return err
	}
	return nil
}

func (me *TTFParser) BytesToString(b []byte) string {
	return string(b)
}

func (me *TTFParser) ReadUShort(fd *os.File) (uint64, error) {
	buff, err := me.Read(fd, 2)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), nil
}

func (me *TTFParser) ReadULong(fd *os.File) (uint64, error) {
	buff, err := me.Read(fd, 4)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), nil
}

func (me *TTFParser) Skip(fd *os.File, length int64) error {
	_, err := fd.Seek(int64(length), 1)
	if err != nil {
		return err
	}
	return nil
}

func (me *TTFParser) Read(fd *os.File, length int) ([]byte, error) {
	buff := make([]byte, length)
	readlength, err := fd.Read(buff)
	if err != nil {
		return nil, err
	}
	if readlength != length {
		return nil, errors.New("file out of length")
	}
	//fmt.Printf("%d,%s\n", readlength, string(buff))
	return buff, nil
}

func (me *TTFParser) CompareBytes(a []byte, b []byte) bool {

	if a == nil && b == nil {
		return true
	} else if a == nil && b != nil {
		return false
	} else if a != nil && b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	i := 0
	length := len(a)
	for i < length {
		if a[i] != b[i] {
			return false
		}
		i++
	}
	return true
}
