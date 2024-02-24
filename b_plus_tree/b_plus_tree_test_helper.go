package b_plus_tree

func NewNodeTest(keys ...int) Node {
	keyArray := make([]Key, len(keys))

	for i, key := range keys {
		keyArray[i] = Key{key, nil}
	}

	return Node{
		Key:     keyArray,
		Pointer: len(keys),
	}
}

func (k *Key) ConnectNodeTest(node *Node) {
	k.NextNode = node
}

func CompareNodeKeyTest(n1, n2 Node) bool {
	for index := range n1.Key {
		if n1.Key[index] != n2.Key[index] {
			return true
		}
	}

	return false
}
