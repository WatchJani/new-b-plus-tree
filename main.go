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

	tree.PositionSearch(b.NewKey(51100, 0))

	fmt.Println(tree.AllRight())

	key, err := tree.GetValueCurrentKey()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(*key)
}
