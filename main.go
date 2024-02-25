package main

import (
	"math/rand"
	b "root/b_plus_tree"
)

func main() {
	BPTree := b.NewBPTree(40_000, 199)

	for value := 0; value < 40000; value++ {
		BPTree.Insert(rand.Intn(40_000), rand.Intn(40_000))
	}

	BPTree.All()

	// fmt.Println(BPTree.Root.Key)
}
