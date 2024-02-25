package main

import (
	"fmt"
	b "root/b_plus_tree"
)

func main() {
	BPTree := b.NewBPTree(40_000, 10)

	for value := 0; value < 40000; value++ {
		BPTree.Insert(value)
	}

	BPTree.All()

	fmt.Println(BPTree.Root.Key)
}
