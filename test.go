package main

import "fmt"

func main() {
	s := make([]int, 0, 10)
	s = []int{1, 2, 3, 4, 5}
	index := 3
	r := make([]int, len(s)-index)
	copy(r, s[index:])
	fmt.Println("s:", s, len(s), "/", cap(s), r, len(r), "/", cap(r))
	s = s[:index]
	fmt.Println("s:", s, len(s), "/", cap(s), r, len(r), "/", cap(r))
	s = append(s, 10)
	fmt.Println("s:", s, len(s), "/", cap(s), r, len(r), "/", cap(r))
	s = append(s, r...)
	fmt.Println("s:", s, len(s), "/", cap(s), r, len(r), "/", cap(r))
}
