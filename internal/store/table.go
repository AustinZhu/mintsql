package store

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type DataType uint

const (
	Int DataType = iota
	Text
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

type Column struct {
	Name string
	Type DataType
}

func NewColumn(name string, dt string) (c *Column) {
	c = new(Column)
	c.Name = name
	switch strings.ToLower(dt) {
	case "int":
		c.Type = Int
	case "text":
		c.Type = Text
	}
	return
}

type Table struct {
	Columns []*Column
	Rows    [][]Cell
}

func NewTable(cols []*Column) (tb *Table) {
	tb = &Table{
		Columns: cols,
		Rows:    make([][]Cell, 0),
	}
	return
}
