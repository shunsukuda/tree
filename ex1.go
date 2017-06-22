package main

import (
	"fmt"
	"strconv"
)

const (
	TRY = 25
)

func main() {
	var binaryRoot BinaryNode
	binaryRoot.New(nil)
	binaryRoot.RandomInsert(TRY, TRY)
	binaryRoot.PrintRoot("bintree", "png")
	binaryRoot.ListingNode()
	fmt.Println(_index, "len:", len(_index))
	var bRoot BNode
	bRoot.New(3)
	for i, e := range _index {
		bRoot.Insert(e)
		bRoot.PrintRoot("b"+strconv.FormatInt(int64(i), 10), "png")
	}
}
