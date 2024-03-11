package b_plus_tree

import "fmt"

func (t *BPTree[K, V]) AllRight() (int, int) {
	current := t.root

	make := make(map[K]struct{})

	//go to left first key
	for current.key[0].nextNode != nil {
		current = current.key[0].nextNode
	}

	var counter int

	var less K

	for current != nil {
		for i := 0; i < current.pointer; i++ {
			make[current.key[i].key] = struct{}{}
			if less <= current.key[i].key {
				counter++
				less = current.key[i].key
			} else {
				break
			}

			fmt.Println(current.key[i])
		}

		fmt.Println()

		current = current.linkNodeRight
	}

	return counter, len(make)
}

func (t *BPTree[K, V]) AllLeft() (int, int) {
	current := t.root

	make := make(map[K]struct{})

	//go to right first key
	for current.key[current.pointer].nextNode != nil {
		current = current.key[current.pointer].nextNode
	}

	var counter int

	var less K = current.key[current.pointer-1].key // well

	for current != nil {
		for i := current.pointer - 1; i >= 0; i-- {
			make[current.key[i].key] = struct{}{}
			if less >= current.key[i].key {
				counter++
				less = current.key[i].key
			} else {
				break
			}

			fmt.Println(current.key[i])
		}

		fmt.Println()

		current = current.linkNodeLeft
	}

	return counter, len(make)
}
