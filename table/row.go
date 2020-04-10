package table

type Row struct {
	NameRaw string
	Name string
	Type string
	Pk bool
	AutoIncrement bool
	Null bool
	Default string
}

func MakeRow(s []string) Row {
	ret := Row{Name: s[0], Type:s[1], Default:s[4]}
	if s[2] == "YES" {ret.Null = true} else {ret.Null = false}
	if s[3] == "PRI" {ret.Pk = true} else {ret.Pk = false}
	if s[5] == "auto_increment" {ret.AutoIncrement = true
	} else  {ret.AutoIncrement = false}
	return ret
}

func CompareRow(a *Row, b *Row) bool {
	if a.Name == b.Name && a.Type == b.Type &&
		a.Pk == b.Pk && a.AutoIncrement == b.AutoIncrement &&
		a.Null == b.Null && a.Default == b.Default{
		return true
	} else {return false}
}