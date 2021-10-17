package table

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

type Cell []byte

func (c Cell) AsInt() (i int32) {
	err := binary.Read(bytes.NewBuffer(c), binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}
	return
}

func (c Cell) AsText() string {
	return string(c)
}

func FromString(s string, dataType DataType) Cell {
	switch dataType {
	case Int:
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	case Text:
		return Cell(s)
	default:
		return nil
	}
}
