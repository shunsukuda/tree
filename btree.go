package main

const (
	Degree = 4
)

type node struct {
	n    int
	leaf bool
	key  []int
	c    []*node
}

func (x *node) NumKey() int {
	x.n = len(x.key)
	return x.n
}

func (n *node) IsLeaf() bool {
	return n.leaf
}

func (n *node) IsFull() bool {
	if n.NumKey() >= Degree*2-1 {
		return true
	} else {
		return false
	}
}

func (n *node) GetKey(index int) int {
	return n.key[index]
}

func (n *node) GetChild(index int) *node {
	return n.c[index]
}

func (x *node) insertionKey(value int, index int) {
	right := make([]int, len(x.key)-index)
	copy(right, x.key[index:])
	x.key = x.key[:index]
	x.key = append(x.key, value)
	x.key = append(x.key, right...)
}

func (x *node) insertionChild(value *node, index int) {
	right := make([]*node, len(x.c)-index)
	copy(right, x.c[index:])
	x.c = x.c[:index]
	x.c = append(x.c, value)
	x.c = append(x.c, right...)
}

func (x *node) New() {
	x.n = 0
	x.leaf = true
	x.key = make([]int, 0, Degree*2-1)
	x.c = make([]*node, 0, Degree*2)
}

func (x node) Get() node {
	return x
}

func (x *node) Seak(key int) (index int, equal bool) {
	index = x.NumKey()
	equal = false
	for i := x.NumKey() - 1; i >= 0; i-- {
		if key > x.GetKey(i) {
			break
		}
		index = i
		if key == x.GetKey(i) {
			equal = true
			break
		}
	}
	return
}

func (x *node) Search(key int) (*node, int) {
	if x.NumKey() == 0 {
		return nil, 0
	}
	index, eq := x.Seak(key)
	if eq {
		return x, index
	}
	if x.IsLeaf() {
		return nil, 0
	} else {
		return x.GetChild(index).Search(key)
	}
}

func (x *node) Split(index int) {
	y := x.GetChild(index)
	var z node
	z.New()
	z.key = y.key[Degree:]
	if !x.IsLeaf() {
		z.c = y.c[Degree:]
	}
	x.key[index] = y.key[Degree-1]
	x.insertionChild(&z, index+1)
	y.key = y.key[:Degree-1]
	y.c = y.c[:Degree]
}

func (t *node) Insert(key int) {
	if t.IsFull() {
		x := t.Get()
		t.New()
		t.leaf = false
		t.c = append(t.c, &x)
		t.Split(0)
	}
	t.InsertNonFull(key)
}

func (x *node) InsertNonFull(key int) {
	index, _ := x.Seak(key)
	if x.IsLeaf() {
		x.insertionKey(key, index)
		return
	} else {
		if x.IsFull() {
			x.Split(index)
			if key > x.GetKey(index) {
				index++
			}
			x.GetChild(index).InsertNonFull(key)
		}
	}
}
