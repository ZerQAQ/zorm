package table

type Index struct {
	Name string
	Unique bool
	ColName string
}

func MakeIndex(s []string) Index {
	var ret = Index{Name: s[2], ColName:s[4]}
	if s[1] == "1" {ret.Unique = true} else {ret.Unique = false}
	return ret
}