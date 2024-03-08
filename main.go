package main

import (
	"fmt"
	"math/rand"

	b "github.com/WatchJani/new-b-plus-tree/b_plus_tree"
)

func main() {
	tree := b.NewBPTree[int, int](50, 5)

	for range 10 {
		tree.Insert(rand.Intn(50000), 5)
	}

	tree.PositionSearch(b.NewKey(5151, 0))

	fmt.Println(tree.All())
	fmt.Println(tree.GetValueCurrentKey())
}
