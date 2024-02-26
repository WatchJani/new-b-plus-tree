package b_plus_tree

import (
	"fmt"
	"math"
)

const (
	NotExist   int = -1
	LeafCode   int = 0
	ParentCode int = 1
)

type ProcessingNode struct {
	CurrentNode *Node
	Stack       //have history of indexes
}

type BPTree struct {
	Memory    []*Node
	Root      *Node
	Degree    int
	MiddleKey int
	ProcessingNode
}

// new B+ Tree
func NewBPTree(capacity, degree int) *BPTree {
	numberOfNode := capacity / degree
	treeHeight := math.Log2(float64(capacity))

	return &BPTree{
		Memory:    make([]*Node, 0, numberOfNode),
		Degree:    degree,
		MiddleKey: degree / 2,
		ProcessingNode: ProcessingNode{
			Stack: newStack(int(treeHeight)),
		},
	}
}

func (t *BPTree) Insert(key, value int) {
	realKey := newKey(key, value) //make real key

	if t.Root == nil { //Check if root exist
		t.Root = t.createNode() //create new root
	}

	t.searchLeaf(realKey)   //find leaf
	t.appendToLeaf(realKey) //add to leaf

	t.splitParent() //check if will parent split

	t.clear()
}

func (t *BPTree) createNode() *Node {
	newNode := newNode(t.Degree)
	t.Memory = append(t.Memory, newNode)

	newNode.Key[0] = Key{}

	return newNode
}

func (t *BPTree) split(code int) *Node {
	newNode := t.createNode()

	newNode.appendKeys(t.CurrentNode.Key[t.MiddleKey+code:t.Degree+code], 0, code)
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
	return positionInsert > 0 && t.CurrentNode.Key[positionInsert-1].Key == key.Key
}

func (t *BPTree) updateExistingKey(positionInsert int, key Key) {
	t.CurrentNode.Key[positionInsert-1] = key
}

func (t *BPTree) insertNewKey(positionInsert int, key Key) {
	t.CurrentNode.insertKey(key, positionInsert)
}

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

func (t *BPTree) splitParent() {
	for t.CurrentNode != nil && t.Degree == t.CurrentNode.Pointer {
		t.newParent()

		newNode := t.split(ParentCode)

		t.updateState(newNode)

		//update parents
		for i := 0; i < newNode.Pointer+1; i++ {
			t.CurrentNode.Key[t.MiddleKey+1+i].NextNode.Parent = newNode
		}

		t.Next()
	}
}

func (t *BPTree) Next() {
	t.CurrentNode = t.CurrentNode.Parent
}

func (t *BPTree) InsertOrUpdate(key Key) {
	positionInsert := t.next()

	if t.shouldUpdate(positionInsert, key) {
		t.updateExistingKey(positionInsert, key)
	} else {
		t.insertNewKey(positionInsert, key)
	}
}

func (t *BPTree) appendToLeaf(key Key) {
	t.InsertOrUpdate(key)

	if t.CurrentNode.Pointer == t.Degree {
		t.newParent()

		newNode := t.split(LeafCode)

		t.updateState(newNode)
	}

	t.Next()
}

func (t *BPTree) updateState(newNode *Node) {
	indexToUpdate := t.next()

	parentKeyIndex := t.MiddleKey
	parentNode := t.CurrentNode.Parent
	parentNode.insertKey(t.CurrentNode.Key[parentKeyIndex], indexToUpdate)

	parentNode.Key[indexToUpdate+1].updateNextNode(newNode)
	parentNode.Key[indexToUpdate].updateNextNode(t.CurrentNode)

	removedPointers := len(t.CurrentNode.Key[:parentKeyIndex]) + t.Degree%2
	t.CurrentNode.Pointer -= removedPointers

	t.CurrentNode.link(newNode)
}

//perfect

type Node struct {
	Pointer  int
	Parent   *Node
	LinkNode *Node
	Key      []Key
}

func newNode(degree int) *Node {
	return &Node{
		Key: make([]Key, degree+1),
	}
}

// search returns the index where the specified key should be inserted in the sorted keys array.
func (n *Node) search(key Key) int {
	for i, currentKey := range n.Key[:n.Pointer] {
		if key.Key < currentKey.Key {
			return i
		}
	}

	return n.Pointer
}

func (n *Node) link(node *Node) {
	if n.LinkNode != nil {
		node.LinkNode = n.LinkNode
	}

	n.LinkNode = node
}

func (n *Node) appendKeys(key []Key, position, code int) {
	copy(n.Key[position:], key)
	n.Pointer += len(key) - code
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

type Key struct {
	Key      int
	Value    int
	NextNode *Node
}

func newKey(key, value int) Key {
	return Key{
		Key:   key,
		Value: value,
	}
}

func (k *Key) updateNextNode(n *Node) {
	k.NextNode = n
}

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

//======================================================================================00

// for testing nothing special
func (t *BPTree) All() {
	current := t.Root

	make := make(map[int]struct{})

	//go to left first key
	for current.Key[0].NextNode != nil {
		current = current.Key[0].NextNode
	}

	var counter int

	var less int = -1

	for current != nil {
		for i := 0; i < current.Pointer; i++ {
			make[current.Key[i].Key] = struct{}{}
			if less <= current.Key[i].Key {
				counter++
				less = current.Key[i].Key
			} else {
				break
			}

			fmt.Println(current.Key[i])
		}

		fmt.Println()

		current = current.LinkNode
	}

	fmt.Println(counter)
	fmt.Println(len(make))
}
