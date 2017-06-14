package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/awalterschulze/gographviz"
)

const (
	OUTFILE = "binary"
)

type BinaryNode struct {
	n     bool
	key   int
	left  *BinaryNode
	right *BinaryNode
}

func (x *BinaryNode) IsEmpty() bool {
	return x.n
}

func (x *BinaryNode) Init() {
	x.n = true
}

func (x *BinaryNode) Update(key int) {
	x.n = false
	x.key = key
}

func (x *BinaryNode) Search(key int) *BinaryNode {
	if x.key == key {
		return x
	} else if x.left != nil && x.key > key {
		return x.left.Search(key)
	} else if x.right != nil && x.key < key {
		return x.right.Search(key)
	} else {
		return nil
	}
}

func (t *BinaryNode) Insert(key int) {
	if t.IsEmpty() {
		t.Update(key)
		return
	}
	if t.key > key {
		if t.left == nil {
			t.left = new(BinaryNode)
			t.left.Init()
		}
		t.left.Insert(key)
	} else if t.key < key {
		if t.right == nil {
			t.right = new(BinaryNode)
			t.right.Init()
		}
		t.right.Insert(key)
	}
}

func (t *BinaryNode) PrintRoot(extension string) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		log.Fatal(err)
	}
	rootName := strconv.FormatInt(int64(t.key), 10)
	t.printNode(graph, rootName)

	if err := ioutil.WriteFile("./"+OUTFILE+".dot", []byte(graph.String()), os.ModePerm); err != nil {
		log.Print("Dot File Write Error")
		log.Fatal(err)
	}
	if err := exec.Command("dot", "-T"+extension, OUTFILE+".dot", "-o", OUTFILE+"."+extension).Run(); err != nil {
		log.Print("Dot Command Error")
		log.Fatal(err)
	}
}

func (x *BinaryNode) printNode(graph *gographviz.Graph, parentName string) {
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
}

func (t *BinaryNode) FullSearch(list []int) {
	fmt.Println(list)
	list = append(list, t.key)
	if t.left != nil {
		t.left.FullSearch(list)
	}
	if t.right != nil {
		t.right.FullSearch(list)
	}
}

func (t *BinaryNode) RandomInsert(max int, times int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < times; i++ {
		key := rand.Intn(max)
		if t.Search(key) == nil {
			t.Insert(key)
		}
	}
}
