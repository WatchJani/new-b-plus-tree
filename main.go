package main

import (
	"fmt"
	"log"
	"math/rand"

	b "github.com/WatchJani/new-b-plus-tree/b_plus_tree"
)

func main() {
	tree := b.NewBPTree[int, int](50, 5)

	for range 10 {
		tree.Insert(rand.Intn(50000), 5)
	}

	tree.PositionSearch(51100)

	fmt.Println(tree.AllRight())

	key, err := tree.GetCurrentKey()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(key)
}
