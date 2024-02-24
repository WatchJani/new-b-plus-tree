package main

import (
	"fmt"
	b "root/b_plus_tree"
)

func main() {
	leaf := b.NewNodeTest(1, 10, 125, 1520, 0)
	key := b.NewKey(59)

	leaf.InsertKey(key, 3)
	fmt.Println(leaf.Key)
}
