package fontmaker

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type TTFParser struct {
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

	for i < numTables {

		i++
	}
	return nil
}

func (me *TTFParser) ReadUShort(fd *os.File) (uint64, error) {
	buff, err := me.Read(fd, 2)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("%#v\n", buff)
	num, length := binary.Uvarint(buff)
	if length == 0 {
		return 0, errors.New("buf too small")
	} else if length < 0 {
		return 0, errors.New("value larger than 64 bits (overflow)")
	}
	return num, nil
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
