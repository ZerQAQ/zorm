package set

type void struct {}
type elmType interface{}

type Set struct {
	M map[elmType]void
}

func MakeSet() *Set {
	s := new(Set)
	s.M	= make(map[elmType]void)
	return s
}

func (s *Set) Insert (v elmType) bool {
	_, ok := s.M[v]
	if ok {return false} else {
		s.M[v] = void{}
		return true
	}
}

func (s *Set) Erase (v elmType) bool {
	_, ok := s.M[v]
	if ok {
		delete(s.M, v)
		return true
	} else {
		return false
	}
}

func (s *Set) Find (v elmType) bool{
	_, ok := s.M[v]
	return ok
}