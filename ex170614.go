package main

import "fmt"

func main() {
	var binaryRoot BinaryNode
	var nodeList []int
	binaryRoot.Init()
	binaryRoot.RandomInsert(25, 10)
	binaryRoot.PrintRoot("png")
	binaryRoot.FullSearch(nodeList)
	fmt.Println(nodeList)
}
