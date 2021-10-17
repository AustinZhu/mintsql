package table

import "strings"

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
