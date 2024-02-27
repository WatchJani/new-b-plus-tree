package b_plus_tree

import (
	"math"
)

const (
	notExist   int = -1
	leafCode   int = 0
	parentCode int = 1
)

type processingNode[K string | int | float64, V any] struct {
	CurrentNode *node[K, V]
	stack       //have history of indexes
}

type bPTree[K string | int | float64, V any] struct {
	Memory    []*node[K, V]
	Root      *node[K, V]
	Degree    int
	MiddleKey int
	processingNode[K, V]
}

// new B+ Tree
func NewBPTree[K string | int | float64, V any](capacity, degree int) *bPTree[K, V] {
	numberOfNode := capacity / degree
	treeHeight := math.Log2(float64(capacity))

	return &bPTree[K, V]{
		Memory:    make([]*node[K, V], 0, numberOfNode),
		Degree:    degree,
		MiddleKey: degree / 2,
		processingNode: processingNode[K, V]{
			stack: newStack(int(treeHeight)),
		},
	}
}

// insert or replace new value to tree
func (t *bPTree[K, V]) Insert(key K, value V) {
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
func (t *bPTree[K, V]) ClearTree() {
	t.Memory = t.Memory[:0]
}

// create new node with one empty key
func (t *bPTree[K, V]) createNode() *node[K, V] {
	newNode := newNode[K, V](t.Degree)
	t.Memory = append(t.Memory, newNode)

	newNode.emptyKey(0)

	return newNode
}

// split node on two identic
func (t *bPTree[K, V]) split(code int) *node[K, V] {
	newNode := t.createNode()

	start, end := t.MiddleKey+code, t.Degree+code
	sliceOfKey := t.CurrentNode.Key[start:end]

	newNode.appendKeys(sliceOfKey, 0, code)
	newNode.Parent = t.CurrentNode.Parent

	return newNode
}

// create new parent of node
func (t *bPTree[K, V]) newParent() {
	if t.CurrentNode.Parent == nil {
		t.CurrentNode.Parent = t.createNode()
		t.Root = t.CurrentNode.Parent
	}
}

// should i update existed key value
func (t bPTree[K, V]) shouldUpdate(positionInsert int, key key[K, V]) bool {
	return positionInsert > 0 && t.CurrentNode.Key[positionInsert-1].Key == key.Key
}

// update value
func (t *bPTree[K, V]) updateExistingKey(positionInsert int, key key[K, V]) {
	t.CurrentNode.Key[positionInsert-1] = key
}

// insert new key in tree
func (t *bPTree[K, V]) insertNewKey(positionInsert int, key key[K, V]) {
	t.CurrentNode.insertKey(key, positionInsert)
}

// returns nex node in level tree
func (t *bPTree[K, V]) nextNode(nextIndex int) *node[K, V] {
	return t.CurrentNode.Key[nextIndex].NextNode
}

// find working leaf
func (t *bPTree[K, V]) searchLeaf(key key[K, V]) {
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
func (t *bPTree[K, V]) splitParent() {
	for t.CurrentNode != nil && t.Degree == t.CurrentNode.Pointer {
		t.newParent()

		newNode := t.split(parentCode)

		t.updateState(newNode)

		//update parents
		for i := 0; i < newNode.Pointer+1; i++ {
			index := t.MiddleKey + 1 + i
			t.CurrentNode.Key[index].NextNode.Parent = newNode
		}

		t.nextParent()
	}
}

// next node
func (t *bPTree[K, V]) nextParent() {
	t.CurrentNode = t.CurrentNode.Parent
}

// check if key exist then update value if not then add new
func (t *bPTree[K, V]) insertOrUpdate(key key[K, V]) {
	positionInsert := t.next()

	if t.shouldUpdate(positionInsert, key) {
		t.updateExistingKey(positionInsert, key)
	} else {
		t.insertNewKey(positionInsert, key)
	}
}

// add new value to the leaf
func (t *bPTree[K, V]) appendToLeaf(key key[K, V]) {
	t.insertOrUpdate(key)

	if t.CurrentNode.Pointer == t.Degree {
		t.newParent()

		newNode := t.split(leafCode)

		t.updateState(newNode)
	}

	t.nextParent()
}

// update parent, new spliced node
func (t *bPTree[K, V]) updateState(newNode *node[K, V]) {
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

type node[K string | int | float64, V any] struct {
	Pointer  int
	Parent   *node[K, V]
	LinkNode *node[K, V]
	Key      []key[K, V]
}

// constructor for new node
func newNode[K string | int | float64, V any](degree int) *node[K, V] {
	return &node[K, V]{
		Key: make([]key[K, V], degree+1),
	}
}

// make empty key
func (n *node[K, V]) emptyKey(position int) {
	n.Key[position] = key[K, V]{}
}

// search returns the index where the specified key should be inserted in the sorted keys array.
func (n *node[K, V]) search(key key[K, V]) int {
	for i, currentKey := range n.Key[:n.Pointer] {
		if key.Key < currentKey.Key {
			return i
		}
	}

	return n.Pointer
}

// link current node with next one
func (n *node[K, V]) link(node *node[K, V]) {
	if n.LinkNode != nil {
		node.LinkNode = n.LinkNode
	}

	n.LinkNode = node
}

// appends all key to current node
func (n *node[K, V]) appendKeys(key []key[K, V], position, code int) {
	copy(n.Key[position:], key)
	n.Pointer += len(key) - code
}

func (n *node[K, V]) increasePointer() {
	n.Pointer++
}

// just append one key
func (n *node[K, V]) appendKey(key key[K, V], position int) {
	n.Key[position] = key
}

// insert key on special position
func (n *node[K, V]) insertKey(key key[K, V], position int) {
	copy(n.Key[position+1:], n.Key[position:])
	n.appendKey(key, position)
	n.increasePointer()
}

type key[K string | int | float64, V any] struct {
	Key      K
	Value    V
	NextNode *node[K, V]
}

func newKey[K string | int | float64, V any](realKey K, value V) key[K, V] {
	return key[K, V]{
		Key:   realKey,
		Value: value,
	}
}

func (k *key[K, V]) updateNextNode(n *node[K, V]) {
	k.NextNode = n
}

type stack struct {
	Stack   []int
	Current int
}

func newStack(capacity int) stack {
	return stack{
		Stack: make([]int, capacity),
	}
}

func (s *stack) add(index int) {
	s.Stack[s.Current] = index
	s.Current++
}

func (s *stack) next() int {
	if s.Current > 0 {
		s.Current--
		return s.Stack[s.Current]
	}

	return 0
}

func (s *stack) clear() {
	s.Current = 0
}

//======================================================================================00

// // for testing nothing special
// func (t *BPTree[K, V]) all() {
// 	current := t.Root

// 	make := make(map[K]struct{})

// 	//go to left first key
// 	for current.Key[0].NextNode != nil {
// 		current = current.Key[0].NextNode
// 	}

// 	var counter int

// 	var less K

// 	for current != nil {
// 		for i := 0; i < current.Pointer; i++ {
// 			make[current.Key[i].Key] = struct{}{}
// 			if less <= current.Key[i].Key {
// 				counter++
// 				less = current.Key[i].Key
// 			} else {
// 				break
// 			}

// 			fmt.Println(current.Key[i])
// 		}

// 		fmt.Println()

// 		current = current.LinkNode
// 	}

// 	fmt.Println(counter)
// 	fmt.Println(len(make))
// }
