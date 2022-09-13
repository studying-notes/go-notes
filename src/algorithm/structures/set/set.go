package set

type Set struct {
	maps map[interface{}]bool
}

func NewSet() *Set {
	return &Set{map[interface{}]bool{}}
}

func (s *Set) Add(item interface{}) {
	s.maps[item] = true
}

func (s *Set) Remove(item interface{}) {
	delete(s.maps, item)
}

func (s *Set) Contains(item interface{}) bool {
	return s.maps[item]
}

// List 转换为无序列表
func (s *Set) List() (list []interface{}) {
	for item := range s.maps {
		list = append(list, item)
	}
	return list
}

func (s *Set) Len() int {
	return len(s.List())
}

// IsEmpty 判断是否为空
func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set) Clear() {
	s.maps = map[interface{}]bool{}
}
