package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/awalterschulze/gographviz"
)

type BNode struct {
	n      int
	leaf   bool
	key    []int
	c      []*BNode
	degree int
}

func (x *BNode) NumKey() int { // {{{
	x.n = len(x.key)
	return x.n
} // }}}

func (n *BNode) IsLeaf() bool { // {{{
	return n.leaf
} // }}}

func (n *BNode) IsFull() bool { // {{{
	if n.NumKey() == n.degree*2-1 {
		return true
	} else {
		return false
	}
} // }}}

func (n *BNode) GetKey(index int) int { // {{{
	if index < 0 || index >= len(n.key) {
		log.Fatal("out of index(key) ", n, "(index:", index, ")")
	}
	return n.key[index]
} // }}}

func (n *BNode) GetChild(index int) *BNode { // {{{
	if index < 0 || index >= len(n.c) {
		log.Fatal("out of index(c) ", n, "(index:", index, ")")
	}
	return n.c[index]
} // }}}

func (x *BNode) New(degree int) { // {{{
	min := 2
	if degree >= min {
		x.degree = degree
	} else {
		log.Fatal("degree must >= " + strconv.FormatInt(int64(min), 10))
	}
	x.n = 0
	x.leaf = true
	x.key = make([]int, 0, x.degree*2-1)
	x.c = make([]*BNode, 0, x.degree*2)
} // }}}

func (x BNode) Get() BNode { // {{{
	return x
} // }}}

func (x *BNode) Search(key int) (*BNode, int) { // {{{
	if x.NumKey() == 0 {
		return nil, 0
	}
	index := sort.SearchInts(x.key, key)
	if index < len(x.key) && key == x.GetKey(index) {
		return x, index
	} else {
		if x.IsLeaf() {
			return nil, 0
		} else {
			return x.GetChild(index).Search(key)
		}
	}
} // }}}

func (x *BNode) Split(index int) { // {{{
	y := x.GetChild(index)
	var z BNode
	z.New(x.degree)
	z.key = y.key[x.degree:]
	x.key = append(x.key, y.key[x.degree-1])
	y.key = y.key[:x.degree-1]
	sort.Ints(x.key)
	if len(y.c) != 0 {
		z.c = y.c[x.degree:]
		y.c = y.c[:x.degree]
	}
	x.insertionChild(&z, index+1)
	if len(y.c) != 0 {
		y.leaf = false
	}
	if len(z.c) != 0 {
		z.leaf = false
	}
} // }}}

func (x *BNode) insertionChild(c *BNode, index int) { // {{{
	right := make([]*BNode, len(x.c)-index)
	copy(right, x.c[index:])
	x.c = x.c[:index]
	x.c = append(x.c, c)
	x.c = append(x.c, right...)
} // }}}

func (t *BNode) Insert(key int) { // {{{
	if t.IsFull() {
		x := t.Get()
		t.New(t.degree)
		t.leaf = false
		t.c = append(t.c, &x)
		t.Split(0)
	}
	t.insertNonFull(key)
} // }}}

func (x *BNode) insertNonFull(key int) { // {{{
	index := sort.SearchInts(x.key, key)
	if x.IsLeaf() {
		x.key = append(x.key, key)
		sort.Ints(x.key)
		return
	}
	if x.GetChild(index).IsFull() {
		x.Split(index)
		index = sort.SearchInts(x.key, key)
	}
	x.GetChild(index).insertNonFull(key)

} // }}}

func (t *BNode) PrintRoot(filename string, extension string) { // {{{
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		log.Fatal(err)
	}
	var rootName string = `"|`
	for _, e := range t.key {
		rootName += strconv.FormatInt(int64(e), 10) + "|"
	}
	rootName += `"`
	t.printNode(graph, rootName)

	dotFilepath := filepath.Join("graph", filename+".dot")
	imgFilepath := filepath.Join("graph", filename+"."+extension)
	if err := ioutil.WriteFile(dotFilepath, []byte(graph.String()), os.ModePerm); err != nil {
		log.Print("Dot File Write Error")
		log.Fatal(err)
	}
	if err := exec.Command("dot", "-T"+extension, dotFilepath, "-o", imgFilepath).Run(); err != nil {
		log.Print("Dot Command Error")
		log.Fatal(err)
	}
} // }}}

func (x *BNode) printNode(graph *gographviz.Graph, parentName string) { // {{{
	var nodeName string = `"|`
	for _, e := range x.key {
		nodeName += strconv.FormatInt(int64(e), 10) + "|"
	}
	nodeName += `"`
	graph.AddNode("G", nodeName, nil)
	if parentName != nodeName {
		graph.AddEdge(parentName, nodeName, true, nil)
	}
	for _, e := range x.c {
		if e != nil {
			e.printNode(graph, nodeName)
		}
	}
} // }}}

func (t *BNode) RandomInsert(max int, times int) { // {{{
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < times; i++ {
		key := rand.Intn(max)
		if p, _ := t.Search(key); p == nil {
			t.Insert(key)
		}
	}
} // }}}
