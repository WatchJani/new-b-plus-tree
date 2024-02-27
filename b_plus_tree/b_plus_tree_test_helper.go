package b_plus_tree

func NewNodeTest(keys ...int) Node[int, int] {
	keyArray := make([]Key[int, int], len(keys))

	for i, key := range keys {
		keyArray[i] = Key[int, int]{key, 0, nil}
	}

	return Node[int, int]{
		Key:     keyArray,
		Pointer: len(keys),
	}
}

func (k *Key[K, V]) ConnectNodeTest(node *Node[K, V]) {
	k.NextNode = node
}

func CompareNodeKeyTest(n1, n2 Node[int, int]) bool {
	for index := range n1.Key {
		if n1.Key[index] != n2.Key[index] {
			return true
		}
	}

	return false
}
