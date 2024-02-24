package b_plus_tree

import (
	"testing"
)

func TestSearchNode(t *testing.T) {
	currentNode := NewNodeTest(1, 4, 21)

	t.Run("smaller than", func(t *testing.T) {
		key := NewKey(3)

		get := currentNode.Search(key)
		want := 1

		assertSearch(t, get, want)
	})

	t.Run("bigger than", func(t *testing.T) {
		key := NewKey(654)

		get := currentNode.Search(key)
		want := 3

		assertSearch(t, get, want)
	})
}

func assertSearch(t testing.TB, get, want int) {
	t.Helper()

	if get != want {
		t.Errorf("get: %d | want: %d", get, want)
	}
}

func BenchmarkSearch(b *testing.B) {
	b.StopTimer()
	currentNode := NewNodeTest(1, 4, 21)
	key := NewKey(3)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		currentNode.Search(key)
	}
}

func TestFindLeaf(t *testing.T) {
	t.Run("just single node", func(t *testing.T) {
		currentNode := NewNodeTest(1, 4, 21)

		key := NewKey(3)
		leaf, index, parentIndex := currentNode.SearchLeaf(key)
		indexWant, parentIndexWant := 1, NotExist

		assertFindLeaf(t, leaf, &currentNode, indexWant, index, parentIndex, parentIndexWant)
	})

	t.Run("more then one node (we have the parent)", func(t *testing.T) {
		rootNode := NewNodeTest(25, 104, 210)
		key := NewKey(110)

		nextNode := NewNodeTest(104, 105, 106, 150, 185)
		rootNode.Key[2].ConnectNodeTest(&nextNode)

		leaf, index, parentIndex := rootNode.SearchLeaf(key)
		foundLeaf, indexWant, parentIndexWant := &nextNode, 3, 2

		assertFindLeaf(t, leaf, foundLeaf, indexWant, index, parentIndex, parentIndexWant)
	})
}

func assertFindLeaf(t testing.TB, leaf, foundLeaf *Node, indexWant, index, parentIndex, parentIndexWant int) {
	t.Helper()

	if leaf != foundLeaf || indexWant != index || parentIndex != parentIndexWant {
		t.Errorf("get leaf: %-v | want leaf: %-v | get parent index: %d | want parent index: %d | get index: %d | want index: %d",
			leaf, foundLeaf, parentIndex, parentIndexWant, index, indexWant)
	}
}

// 2ns/op per node
func BenchmarkFindLeaf(b *testing.B) {
	b.StartTimer()
	rootNode := NewNodeTest(25, 104, 210)
	key := NewKey(110)

	nextNode := NewNodeTest(104, 105, 106, 150, 185)
	rootNode.Key[2].ConnectNodeTest(&nextNode)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		rootNode.SearchLeaf(key)
	}
}

func TestInsertKey(t *testing.T) {
	leaf := NewNodeTest(1, 10, 125, 1520, 0)
	key := NewKey(59)

	leaf.InsertKey(key, 2)
	want := NewNodeTest(1, 10, 59, 125, 1520)

	if CompareNodeKeyTest(leaf, want) {
		t.Errorf("get: %-v want: %-v", leaf.Key, want.Key)
	}
}

func BenchmarkInsertKey(b *testing.B) {
	b.StopTimer()

	leaf := NewNodeTest(1, 10, 125, 1520, 0)
	key := NewKey(59)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		leaf.InsertKey(key, 2)
	}
}
