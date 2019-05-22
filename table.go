package goose

type ColumnType string

type Column struct {
	Name     string
	Type     ColumnType
	Length   int
	Default  string
	Null     bool
	OnUpdate string
	Collate  string
}

type Key struct {
	Name    string
	Column  Column
	Primary bool
}

type Constraint struct {
	Name            string
	LocalColumn     Column
	ReferenceTable  Table
	ReferenceColumn Column
}

type Table struct {
	Name    string
	Columns *ColumnMap
	Keys    []Key
	Engine  string
	Charset string
	Collate string
}

type ColumnMap struct {
	columns map[string]Column
	keys    []string
}

func (n *ColumnMap) Set(k string, v Column) {
	n.columns[k] = v
	n.keys = append(n.keys, k)
}

func CreateTable(name string) *Table {
	return &Table{
		Name: name,
	}
}

func (t *Table) SetColumns(columns []Column) *Table {
	var colMap = &ColumnMap{}
	for _, col := range columns {
		colMap.Set(col.Name, col)
	}
	t.Columns = colMap
	return t
}

func (t *Table) AddColumn(column Column, after Column) *Table {
	//t.Columns[]
	return t
}
