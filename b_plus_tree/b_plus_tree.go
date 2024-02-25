package b_plus_tree

import (
	"math"
)

const NotExist int = -1

type Stack struct {
	Stack   []int
	Current int
}

func newStack(capacity int) Stack {
	return Stack{
		Stack: make([]int, capacity),
	}
}

func (s *Stack) add(index int) {
	s.Stack[s.Current] = index
	s.Current++
}

func (s *Stack) next() int {
	if s.Current > 0 {
		s.Current--
		return s.Stack[s.Current]
	}

	return 0
}

func (s *Stack) clear() {
	s.Current = 0
}

type Key struct {
	Value    int
	NextNode *Node
}

func (k *Key) updateNextNode(n *Node) {
	k.NextNode = n
}

// *fix that (maybe need the pointer)
func newKey(value int) Key {
	return Key{
		Value: value,
	}
}

type Node struct {
	Pointer  int
	Parent   *Node
	LinkNode *Node
	Key      []Key
}

// search returns the index where the specified key should be inserted in the sorted keys array.
func (n *Node) search(key Key) int {
	for i, currentKey := range n.Key[:n.Pointer] {
		if key.Value < currentKey.Value {
			return i
		}
	}

	return n.Pointer
}

// bed version
// find the leaf, parent index if exist, and find position in node to insert new element
// func (n *Node) SearchLeaf(key Key) (*Node, int, int) {
// 	currentNode, nextIndex, prevuesIndex := n, n.Search(key), NotExist

// 	for currentNode.Key[nextIndex].NextNode != nil {
// 		prevuesIndex = nextIndex
// 		currentNode = currentNode.Key[nextIndex].NextNode
// 		nextIndex = currentNode.Search(key)
// 	}

// 	return currentNode, nextIndex, prevuesIndex
// }

// func (t BPTree) SearchNextIndex(key Key) int {
// 	return t.CurrentNode.Search(key)
// }

// func (t *BPTree) MoveToNextNode(nextIndex int) {
// 	t.Add(nextIndex)
// 	t.CurrentNode = t.NextNode(nextIndex)
// }

// func (t BPTree) NextNode(nextIndex int) *Node {
// 	return t.CurrentNode.Key[nextIndex].NextNode
// }

// func (t *BPTree) SearchLeaf(key Key) {
// 	t.CurrentNode = t.Root
// 	nextIndex := t.SearchNextIndex(key)

// 	for t.NextNode(nextIndex) != nil {
// 		t.MoveToNextNode(nextIndex)
// 		nextIndex = t.SearchNextIndex(key)
// 	}

// 	t.Add(nextIndex)
// }

// best version
// func (t *BPTree) SearchLeaf(key Key) {
// 	t.CurrentNode = t.Root
// 	nextIndex := t.CurrentNode.Search(key)

// 	for t.CurrentNode.Key[nextIndex].NextNode != nil {
// 		t.Add(nextIndex)
// 		t.CurrentNode = t.CurrentNode.Key[nextIndex].NextNode
// 		nextIndex = t.CurrentNode.Search(key)
// 	}

// 	t.Add(nextIndex)
// }

func (t *BPTree) searchLeaf(key Key) {
	t.CurrentNode = t.Root

	for {
		nextIndex := t.CurrentNode.search(key)
		t.add(nextIndex)

		if nextNode := t.CurrentNode.Key[nextIndex].NextNode; nextNode == nil {
			break
		} else {
			t.CurrentNode = nextNode
		}
	}
}

func newNode(degree int) *Node {
	return &Node{
		Key: make([]Key, degree+1),
	}
}

func (n *Node) link(node *Node) {
	if n.LinkNode != nil {
		node.LinkNode = n.LinkNode
	}

	n.LinkNode = node
}

func (t *BPTree) createNode() *Node {
	newNode := newNode(t.Degree)
	t.Memory = append(t.Memory, newNode)

	newNode.Key[0] = Key{}

	return newNode
}

func (t *BPTree) Insert(key int) {
	createKey := newKey(key) //make real key

	if t.Root == nil { //Check if root exist
		t.Root = t.createNode() //create new root
	}

	t.searchLeaf(createKey) //find leaf
	t.appendToLeaf(createKey)

	if t.CurrentNode == nil {
		// fmt.Println("nil")
	}

	t.clear()
}

func (t *BPTree) appendToLeaf(key Key) {
	positionInsert := t.next()

	if t.shouldUpdate(positionInsert, key) {
		t.updateExistingKey(positionInsert, key)
	} else {
		t.insertNewKey(positionInsert, key)
	}

	if t.CurrentNode.Pointer == t.Degree {
		t.newParent()

		middleKey := t.Degree / 2
		newNode := t.split(middleKey)

		t.updateState(middleKey, newNode)
	}

	t.CurrentNode = t.CurrentNode.Parent
}

func (t *BPTree) updateState(middleKey int, newNode *Node) {
	indexToUpdate := t.next()
	//add and update parent
	t.CurrentNode.Parent.insertKey(t.CurrentNode.Key[middleKey], indexToUpdate)
	t.CurrentNode.Parent.Key[indexToUpdate+1].updateNextNode(newNode)

	//link current Node and decrees
	t.CurrentNode.Pointer -= (len(t.CurrentNode.Key[:middleKey]) + 1)
	t.CurrentNode.link(newNode)
}

func (n *Node) appendKeys(key []Key, position int) {
	copy(n.Key[position:], key)
	n.Pointer += len(key)
}

func (t *BPTree) split(middleKey int) *Node {
	newNode := t.createNode()

	newNode.appendKeys(t.CurrentNode.Key[middleKey:t.Degree], 0)
	newNode.Parent = t.CurrentNode.Parent

	return newNode
}

func (t *BPTree) newParent() {
	if t.CurrentNode.Parent == nil {
		t.CurrentNode.Parent = t.createNode()
		t.Root = t.CurrentNode.Parent
	}
}

func (t BPTree) shouldUpdate(positionInsert int, key Key) bool {
	return positionInsert > 0 && t.CurrentNode.Key[positionInsert-1] == key
}

func (t *BPTree) updateExistingKey(positionInsert int, key Key) {
	t.CurrentNode.Key[positionInsert-1] = key
}

func (t *BPTree) insertNewKey(positionInsert int, key Key) {
	t.CurrentNode.insertKey(key, positionInsert)
}

func (n *Node) increasePointer() {
	n.Pointer++
}

func (n *Node) appendKey(key Key, position int) {
	n.Key[position] = key
}

func (n *Node) insertKey(key Key, position int) {
	copy(n.Key[position+1:], n.Key[position:])
	n.appendKey(key, position)
	n.increasePointer()
}

type BPTree struct {
	Memory []*Node
	Root   *Node
	Degree int
	ProcessingNode
}

func NewBPTree(capacity, degree int) *BPTree {
	numberOfNode := capacity / degree
	treeHeight := math.Log2(float64(capacity))

	return &BPTree{
		Memory: make([]*Node, 0, numberOfNode),
		Degree: degree,
		ProcessingNode: ProcessingNode{
			Stack: newStack(int(treeHeight)),
		},
	}
}

type ProcessingNode struct {
	CurrentNode *Node
	Stack       //have history of indexes
}
