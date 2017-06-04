package main

import "fmt"

type S struct {
	i []int
}

func (s S) Get() S {
	return s
}

func (s *S) New() {
	s.i = make([]int, 0, 10)
}

func (s *S) f() {
	x := s.Get()
	s.New()
	s.i = append(s.i, 10)
	fmt.Println(x, s)
}

func insertion(in []int, value int, index int) (out []int) {
	right := make([]int, len(in)-index)
	copy(right, in[index:])
	out = in[:index]
	out = append(out, value)
	out = append(out, right...)
	return
}

func main() {
	var 
	var x int
	fmt.Scanf("%s
}
