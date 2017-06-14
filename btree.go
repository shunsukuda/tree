package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

const (
	DEGREE  = 4
	OUT_DOT = "btree.dot"
	OUT_PNG = "btree.png"
)

type BNode struct {
	n    int
	leaf bool
	key  []int
	c    []*BNode
}

func (x *BNode) NumKey() int {
	x.n = len(x.key)
	return x.n
}

func (n *BNode) IsLeaf() bool {
	return n.leaf
}

func (n *BNode) IsFull() bool {
	if n.NumKey() >= DEGREE*2-1 {
		return true
	} else {
		return false
	}
}

func (n *BNode) GetKey(index int) int {
	return n.key[index]
}

func (n *BNode) GetChild(index int) *BNode {
	return n.c[index]
}

func (x *BNode) insertionKey(value int, index int) {
	right := make([]int, len(x.key)-index)
	copy(right, x.key[index:])
	x.key = x.key[:index]
	x.key = append(x.key, value)
	x.key = append(x.key, right...)
}

func (x *BNode) insertionChild(value *BNode, index int) {
	right := make([]*BNode, len(x.c)-index)
	copy(right, x.c[index:])
	x.c = x.c[:index]
	x.c = append(x.c, value)
	x.c = append(x.c, right...)
}

func (x *BNode) Init() {
	x.n = 0
	x.leaf = true
	x.key = make([]int, 0, DEGREE*2-1)
	x.c = make([]*BNode, 0, DEGREE*2)
}

func (x BNode) Get() BNode {
	return x
}

func (x *BNode) Seak(key int) (index int, equal bool) {
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

func (x *BNode) Search(key int) (*BNode, int) {
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

func (x *BNode) Split(index int) {
	y := x.GetChild(index)
	var z BNode
	z.Init()
	z.key = y.key[DEGREE:]
	if !x.IsLeaf() {
		z.c = y.c[DEGREE:]
	}
	x.key[index] = y.key[DEGREE-1]
	x.insertionChild(&z, index+1)
	y.key = y.key[:DEGREE-1]
	y.c = y.c[:DEGREE]
}

func (t *BNode) Insert(key int) {
	if t.IsFull() {
		x := t.Get()
		t.Init()
		t.leaf = false
		t.c = append(t.c, &x)
		t.Split(0)
	}
	t.InsertNonFull(key)
}

func (x *BNode) InsertNonFull(key int) {
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

func (t *BNode) PrintRoot() {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		log.Fatal(err)
	}
	var rootName string = "|"
	for _, e := range t.key {
		rootName += strconv.FormatInt(int64(e), 10) + "|"
	}
	t.printNode(graph, rootName)

	if err := ioutil.WriteFile("./"+OUT_DOT, []byte(graph.String()), os.ModePerm); err != nil {
		log.Print("Dot File Write Error")
		log.Fatal(err)
	}
	if err := exec.Command("dot", "-Tpng", OUT_DOT, "-o", OUT_PNG).Run(); err != nil {
		log.Print("Dot Command Error")
		log.Fatal(err)
	}
}

func (x *BNode) printNode(graph *gographviz.Graph, parentName string) {
	var nodeName string = "|"
	for _, e := range t.key {
		nodeName += strconv.FormatInt(int64(e), 10) + "|"
	}
	graph.AddNode("G", nodeName, nil)
	if parentName != nodeName {
		graph.AddEdge(parentName, nodeName, true, nil)
	}
	for _, e := range x.c {
		if e != nil {
			e.printNode(graph, nodeName)
		}
	}
}
