package set

// 集合数据结构实现

type Set[T int | string] struct {
	length int
	m      map[T]bool
}

func NewSet[T int | string]() *Set[T] {
	return &Set[T]{m: map[T]bool{}}
}

func (s *Set[T]) Add(item T) {
	if s.Contains(item) {
		return
	}

	s.m[item] = true
	s.length++
}

func (s *Set[T]) Remove(item T) {
	if !s.Contains(item) {
		return
	}

	delete(s.m, item)
	s.length--
}

func (s *Set[T]) Contains(item T) bool {
	return s.m[item]
}

func (s *Set[T]) ToList() (list []T) {
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *Set[T]) Len() int {
	return s.length
}

func (s *Set[T]) IsEmpty() bool {
	return s.length == 0
}

func (s *Set[T]) Clear() {
	s.m = map[T]bool{}
	s.length = 0
}
