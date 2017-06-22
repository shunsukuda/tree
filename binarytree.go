package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/awalterschulze/gographviz"
)

type BinaryNode struct {
	empty  bool
	key    int
	parent *BinaryNode
	left   *BinaryNode
	right  *BinaryNode
}

var _index []int

func (x *BinaryNode) IsEmpty() bool { // {{{
	return x.empty
} // }}}

func (x *BinaryNode) IsLeaf() bool { // {{{
	if x.left == nil && x.right == nil {
		return true
	} else {
		return false
	}
} // }}}

func (x *BinaryNode) New(parent *BinaryNode) { // {{{
	x.empty = true
	x.key = -1
	x.parent = parent
	x.left = nil
	x.right = nil
} // }}}

func (x *BinaryNode) NewLeft(key int) { // {{{
	x.left = new(BinaryNode)
	x.left.New(x)
	x.left.key = key
	x.left.empty = false
} // }}}

func (x *BinaryNode) NewRight(key int) { // {{{
	x.right = new(BinaryNode)
	x.right.New(x)
	x.right.key = key
	x.right.empty = false
} // }}}

func (x *BinaryNode) Search(key int) *BinaryNode { // {{{
	if x.key == key {
		return x
	} else if x.left != nil && x.key > key {
		return x.left.Search(key)
	} else if x.right != nil && x.key < key {
		return x.right.Search(key)
	} else {
		return nil
	}
} // }}}

func (x *BinaryNode) SearchMin() *BinaryNode { // {{{
	if x.left != nil {
		return x.left.SearchMin()
	} else {
		return x
	}
} // }}}

func (x *BinaryNode) SearchMax() *BinaryNode { // {{{
	if x.right != nil {
		return x.right.SearchMax()
	} else {
		return x
	}
} // }}}

func (t *BinaryNode) Insert(key int) { // {{{
	if t.IsEmpty() {
		t.key = key
		t.empty = false
		return
	}
	if t.key > key {
		if t.left == nil {
			t.NewLeft(key)
			return
		}
		t.left.Insert(key)
	} else if t.key < key {
		if t.right == nil {
			t.NewRight(key)
			return
		}
		t.right.Insert(key)
	}
} // }}}

func (t *BinaryNode) Remove(key int) bool { // {{{
	if t.IsEmpty() {
		return false
	}
	x := t.Search(key)
	if x == nil {
		return false
	}

	p := x.parent
	if x.left == nil && x.right == nil {
		if x == p.left {
			p.left = nil
		} else { // if x == p.right {
			p.right = nil
		}
	} else if x.left != nil && x.right == nil {
		if x == p.left {
			p.left = x.left
			x.left.parent = p
		} else { // if x == p.right {
			p.right = x.left
			x.left.parent = p
		}
	} else if x.left == nil && x.right != nil {
		if x == p.left {
			p.left = x.right
			x.right.parent = p
		} else { // if x == p.right {
			p.right = x.right
			x.right.parent = p
		}
	} else { // if x.left != nil && x.right != nil {
		m := x.right.SearchMin()
		x.Remove(m.key)
		x.key = m.key
	}
	if x.parent == nil && x.IsLeaf() {
		x.empty = true
	}
	return true
} // }}}

func (t *BinaryNode) PrintRoot(filename string, extension string) { // {{{
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		log.Fatal(err)
	}
	rootName := strconv.FormatInt(int64(t.key), 10)
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

func (x *BinaryNode) printNode(graph *gographviz.Graph, parentName string) { // {{{
	nodeName := strconv.FormatInt(int64(x.key), 10)
	graph.AddNode("G", "N"+nodeName, nil)
	if parentName != nodeName {
		graph.AddEdge("N"+parentName, "N"+nodeName, true, nil)
	}
	if x.left != nil {
		x.left.printNode(graph, nodeName)
	} else {
		graph.AddNode("G", "L"+nodeName, nil)
		graph.AddEdge("N"+nodeName, "L"+nodeName, true, nil)
	}
	if x.right != nil {
		x.right.printNode(graph, nodeName)
	} else {
		graph.AddNode("G", "R"+nodeName, nil)
		graph.AddEdge("N"+nodeName, "R"+nodeName, true, nil)
	}
} // }}}

func (t *BinaryNode) ListingNode() { // {{{
	var list []int
	if t.IsEmpty() {
		return
	}
	t.listingNode(list)
	return
} // }}}

func (x *BinaryNode) listingNode(list []int) { // {{{
	list = append(list, x.key)
	if x.left != nil {
		x.left.listingNode(list)
	}
	if x.right != nil {
		x.right.listingNode(list)
	}
	if x.IsLeaf() {
		for _, el := range list {
			contains := false
			for _, ei := range _index {
				if ei == el {
					contains = true
					break
				}
			}
			if !contains {
				_index = append(_index, el)
			}
		}
	}
} // }}}

func (t *BinaryNode) RandomInsert(max int, times int) { // {{{
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < times; i++ {
		key := rand.Intn(max)
		if t.Search(key) == nil {
			t.Insert(key)
		}
	}
} // }}}
