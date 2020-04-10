package table

type Table struct {
	Name   string
	RowsName []string
	Rows   map[string]Row
	Indexs map[string]Index
}

func (t *Table) Init()  {
	t.Rows = make(map[string]Row)
	t.Indexs = make(map[string]Index)
}

func (t *Table) UpdateRowList (oth *Table)  {
	for _, elm := range oth.Rows{
		t.RowsName = append(t.RowsName, elm.Name)
	}
}

func IsContain(a *Table, b *Table) bool {
	if a.Name != b.Name {return false}
	if len(a.Rows) > len(b.Rows) {return false}
	for _, elm := range a.Rows {
		v, ok := b.Rows[elm.Name]
		if !ok || !CompareRow(&elm, &v) {return false}
	}
	if len(a.Indexs) > len(b.Indexs) {return false}
	for _, elm := range a.Indexs {
		v, ok := b.Indexs[elm.Name]
		if !ok || !CompareIndex(&elm, &v) {return false}
	}
	return true
}

func Sub(a *Table, b *Table) Table {
	ret := Table{Name:a.Name}
	ret.Init()

	for _, elm := range a.Rows {
		if _, ok := b.Rows[elm.Name]; !ok {
			ret.Rows[elm.Name] = elm
		}
	}
	for _, elm := range a.Indexs {
		if _, ok := b.Indexs[elm.Name]; !ok {
			ret.Indexs[elm.Name] = elm
		}
	}

	return ret
}