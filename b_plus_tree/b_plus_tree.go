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

type ProcessingNode[K string | int | float64, V any] struct {
	CurrentNode *Node[K, V]
	Stack       //have history of indexes
}

type BPTree[K string | int | float64, V any] struct {
	Memory    []*Node[K, V]
	Root      *Node[K, V]
	Degree    int
	MiddleKey int
	ProcessingNode[K, V]
}

// new B+ Tree
func NewBPTree[K string | int | float64, V any](capacity, degree int) *BPTree[K, V] {
	numberOfNode := capacity / degree
	treeHeight := math.Log2(float64(capacity))

	return &BPTree[K, V]{
		Memory:    make([]*Node[K, V], 0, numberOfNode),
		Degree:    degree,
		MiddleKey: degree / 2,
		ProcessingNode: ProcessingNode[K, V]{
			Stack: newStack(int(treeHeight)),
		},
	}
}

// insert or replace new value to tree
func (t *BPTree[K, V]) Insert(key K, value V) {
	realKey := newKey(key, value) //make real key

	if t.Root == nil { //Check if root exist
		t.Root = t.createNode() //create new root
	}

	t.searchLeaf(realKey)   //find leaf
	t.appendToLeaf(realKey) //add to leaf

	t.splitParent() //check if will parent split

	t.clear()
}

// make more memory for the tree
func (t *BPTree[K, V]) ClearTree() {
	t.Memory = t.Memory[:0]
}

// create new node with one empty key
func (t *BPTree[K, V]) createNode() *Node[K, V] {
	newNode := newNode[K, V](t.Degree)
	t.Memory = append(t.Memory, newNode)

	newNode.EmptyKey(0)

	return newNode
}

// split node on two identic
func (t *BPTree[K, V]) split(code int) *Node[K, V] {
	newNode := t.createNode()

	start, end := t.MiddleKey+code, t.Degree+code
	sliceOfKey := t.CurrentNode.Key[start:end]

	newNode.appendKeys(sliceOfKey, 0, code)
	newNode.Parent = t.CurrentNode.Parent

	return newNode
}

// create new parent of node
func (t *BPTree[K, V]) newParent() {
	if t.CurrentNode.Parent == nil {
		t.CurrentNode.Parent = t.createNode()
		t.Root = t.CurrentNode.Parent
	}
}

// should i update existed key value
func (t BPTree[K, V]) shouldUpdate(positionInsert int, key Key[K, V]) bool {
	return positionInsert > 0 && t.CurrentNode.Key[positionInsert-1].Key == key.Key
}

// update value
func (t *BPTree[K, V]) updateExistingKey(positionInsert int, key Key[K, V]) {
	t.CurrentNode.Key[positionInsert-1] = key
}

// insert new key in tree
func (t *BPTree[K, V]) insertNewKey(positionInsert int, key Key[K, V]) {
	t.CurrentNode.insertKey(key, positionInsert)
}

// returns nex node in level tree
func (t *BPTree[K, V]) nextNode(nextIndex int) *Node[K, V] {
	return t.CurrentNode.Key[nextIndex].NextNode
}

// find working leaf
func (t *BPTree[K, V]) searchLeaf(key Key[K, V]) {
	t.CurrentNode = t.Root

	for {
		nextIndex := t.CurrentNode.search(key)
		t.add(nextIndex)

		if nextNode := t.nextNode(nextIndex); nextNode == nil {
			break
		} else {
			t.CurrentNode = nextNode
		}
	}
}

// split all full parents node
func (t *BPTree[K, V]) splitParent() {
	for t.CurrentNode != nil && t.Degree == t.CurrentNode.Pointer {
		t.newParent()

		newNode := t.split(ParentCode)

		t.updateState(newNode)

		//update parents
		for i := 0; i < newNode.Pointer+1; i++ {
			index := t.MiddleKey + 1 + i
			t.CurrentNode.Key[index].NextNode.Parent = newNode
		}

		t.Next()
	}
}

// next node
func (t *BPTree[K, V]) Next() {
	t.CurrentNode = t.CurrentNode.Parent
}

// check if key exist then update value if not then add new
func (t *BPTree[K, V]) InsertOrUpdate(key Key[K, V]) {
	positionInsert := t.next()

	if t.shouldUpdate(positionInsert, key) {
		t.updateExistingKey(positionInsert, key)
	} else {
		t.insertNewKey(positionInsert, key)
	}
}

// add new value to the leaf
func (t *BPTree[K, V]) appendToLeaf(key Key[K, V]) {
	t.InsertOrUpdate(key)

	if t.CurrentNode.Pointer == t.Degree {
		t.newParent()

		newNode := t.split(LeafCode)

		t.updateState(newNode)
	}

	t.Next()
}

// update parent, new spliced node
func (t *BPTree[K, V]) updateState(newNode *Node[K, V]) {
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

type Node[K string | int | float64, V any] struct {
	Pointer  int
	Parent   *Node[K, V]
	LinkNode *Node[K, V]
	Key      []Key[K, V]
}

// constructor for new node
func newNode[K string | int | float64, V any](degree int) *Node[K, V] {
	return &Node[K, V]{
		Key: make([]Key[K, V], degree+1),
	}
}

// make empty key
func (n *Node[K, V]) EmptyKey(position int) {
	n.Key[position] = Key[K, V]{}
}

// search returns the index where the specified key should be inserted in the sorted keys array.
func (n *Node[K, V]) search(key Key[K, V]) int {
	for i, currentKey := range n.Key[:n.Pointer] {
		if key.Key < currentKey.Key {
			return i
		}
	}

	return n.Pointer
}

// link current node with next one
func (n *Node[K, V]) link(node *Node[K, V]) {
	if n.LinkNode != nil {
		node.LinkNode = n.LinkNode
	}

	n.LinkNode = node
}

// appends all key to current node
func (n *Node[K, V]) appendKeys(key []Key[K, V], position, code int) {
	copy(n.Key[position:], key)
	n.Pointer += len(key) - code
}

func (n *Node[K, V]) increasePointer() {
	n.Pointer++
}

// just append one key
func (n *Node[K, V]) appendKey(key Key[K, V], position int) {
	n.Key[position] = key
}

// insert key on special position
func (n *Node[K, V]) insertKey(key Key[K, V], position int) {
	copy(n.Key[position+1:], n.Key[position:])
	n.appendKey(key, position)
	n.increasePointer()
}

type Key[K string | int | float64, V any] struct {
	Key      K
	Value    V
	NextNode *Node[K, V]
}

func newKey[K string | int | float64, V any](key K, value V) Key[K, V] {
	return Key[K, V]{
		Key:   key,
		Value: value,
	}
}

func (k *Key[K, V]) updateNextNode(n *Node[K, V]) {
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
func (t *BPTree[K, V]) All() {
	current := t.Root

	make := make(map[K]struct{})

	//go to left first key
	for current.Key[0].NextNode != nil {
		current = current.Key[0].NextNode
	}

	var counter int

	var less K

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
