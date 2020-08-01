package def

// HashSet 简易实现
type HashSet struct {
	maps map[interface{}]bool
}

func NewHashSet() *HashSet {
	return &HashSet{map[interface{}]bool{}}
}

func (s *HashSet) Add(i interface{}) bool {
	// 索引不存在情况返回类型默认初始值
	isExist := s.maps[i]
	if !isExist {
		s.maps[i] = true
	}
	return !isExist
}

func (s *HashSet) Contains(i interface{}) bool {
	isExist := s.maps[i]
	return isExist
}

func (s *HashSet) Remove(i interface{}) {
	delete(s.maps, i)
}
