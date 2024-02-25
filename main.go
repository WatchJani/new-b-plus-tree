package main

import (
	"fmt"
	b "root/b_plus_tree"
)

func main() {
	BPTree := b.NewBPTree(40_000, 5)

	for value := 0; value < 10; value++ {
		BPTree.Insert(value)
	}

	fmt.Println(BPTree.Root)
}
