package table

type Index struct {
	Name string
	Unique bool
	ColName string
}

func MakeIndex(s []string) Index {
	var ret = Index{ColName:s[4]}
	if s[1] == "1" {ret.Unique = false} else {ret.Unique = true}
	if ret.Unique {ret.Name = "unique_" + ret.ColName} else{
		ret.Name = "index_" + ret.ColName
	}
	return ret
}

func CompareIndex(a *Index, b *Index) bool {
	if a.Unique == b.Unique && a.ColName == b.ColName {
		return true
	} else {return false}
}