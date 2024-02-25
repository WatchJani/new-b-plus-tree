package main

import (
	"fmt"
	"math/rand"
	b "root/b_plus_tree"
)

func main() {
	BPTree := b.NewBPTree(40_000, 10)

	for value := 0; value < 40000; value++ {
		BPTree.Insert(rand.Intn(100_000_000))
	}

	BPTree.All()

	fmt.Println(BPTree.Root.Key)
}
