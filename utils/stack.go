package utils

type Stack struct {
	data []int
	sp   int // stack pointer
}

func NewStack(size int) *Stack {
	if size < 1 {
		size = 1
	}
	return &Stack{data: make([]int, size)}
}

func (s *Stack) Push(v int) {
	s.data[s.sp] = v
	s.sp++
}

func (s *Stack) Pop() int {
	s.sp--
	return s.data[s.sp]
}

func (s *Stack) IsEmpty() bool {
	return s.sp == 0
}

func (s *Stack) Cap() int {
	return len(s.data)
}

func (s *Stack) Reset() {
	s.sp = 0
}
