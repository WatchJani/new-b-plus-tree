package b_plus_tree

import (
	"errors"
	"fmt"
	"math"
)

const (
	notExist   int = -1
	leafCode   int = 0
	parentCode int = 1
)

type processingNode[K string | int | float64, V any] struct {
	currentNode *node[K, V]
	stack       //have history of indexes
}

type bPTree[K string | int | float64, V any] struct {
	memory    []*node[K, V]
	root      *node[K, V]
	degree    int
	middleKey int
	processingNode[K, V]
	keyPointer[K, V]
}

type keyPointer[K string | int | float64, V any] struct {
	pointerPosition int
	pointerNode     *node[K, V]
}

// new B+ Tree
func NewBPTree[K string | int | float64, V any](capacity, degree int) *bPTree[K, V] {
	numberOfNode := capacity / degree
	treeHeight := math.Log2(float64(capacity))

	return &bPTree[K, V]{
		memory:    make([]*node[K, V], 0, numberOfNode),
		degree:    degree,
		middleKey: degree / 2,
		processingNode: processingNode[K, V]{
			stack: newStack(int(treeHeight)),
		},
	}
}

// insert or replace new value to tree
func (t *bPTree[K, V]) Insert(key K, value V) {
	realKey := NewKey(key, value) //make real key

	if t.root == nil { //Check if root exist
		t.root = t.createNode() //create new root
	}

	t.searchLeaf(realKey)   //find leaf
	t.appendToLeaf(realKey) //add to leaf

	t.splitParent() //check if will parent split

	t.clear()
}

// make more memory for the tree
func (t *bPTree[K, V]) ClearTree() {
	t.memory = t.memory[:0]
}

// create new node with one empty key
func (t *bPTree[K, V]) createNode() *node[K, V] {
	newNode := newNode[K, V](t.degree)
	t.memory = append(t.memory, newNode)

	newNode.emptyKey(0)

	return newNode
}

// split node on two identic
func (t *bPTree[K, V]) split(code int) *node[K, V] {
	newNode := t.createNode()

	start, end := t.middleKey+code, t.degree+code
	sliceOfKey := t.currentNode.key[start:end]

	newNode.appendKeys(sliceOfKey, 0, code)
	newNode.parent = t.currentNode.parent

	return newNode
}

// create new parent of node
func (t *bPTree[K, V]) newParent() {
	if t.currentNode.parent == nil {
		t.currentNode.parent = t.createNode()
		t.root = t.currentNode.parent
	}
}

// should i update existed key value
func (t bPTree[K, V]) shouldUpdate(positionInsert int, key key[K, V]) bool {
	return positionInsert > 0 && t.currentNode.key[positionInsert-1].key == key.key
}

// update value
func (t *bPTree[K, V]) updateExistingKey(positionInsert int, key key[K, V]) {
	t.currentNode.key[positionInsert-1] = key
}

// insert new key in tree
func (t *bPTree[K, V]) insertNewKey(positionInsert int, key key[K, V]) {
	t.currentNode.insertKey(key, positionInsert)
}

// returns nex node in level tree
func (t *bPTree[K, V]) nextNode(nextIndex int) *node[K, V] {
	return t.currentNode.key[nextIndex].nextNode
}

// search current node
func (t *bPTree[K, V]) search(key key[K, V]) (*node[K, V], int) {
	current := t.root

	for {
		nextIndex := current.search(key)

		if nextNode := current.key[nextIndex].nextNode; nextNode == nil {
			return current, nextIndex
		} else {
			current = nextNode
		}
	}
}

// return current position of key | need to use with NextKey() func
func (t *bPTree[K, V]) PositionSearch(key key[K, V]) {
	t.pointerNode, t.pointerPosition = t.search(key)

}

// return to use current value
func (t *bPTree[K, V]) GetValueCurrentKey() (*K, error) {
	if t.pointerPosition != 0 {
		return &t.pointerNode.key[t.pointerPosition-1].key, nil
	}

	return nil, errors.New("Three is empty")
}

func (t *bPTree[K, V]) NextKey() error {
	t.pointerPosition++

	if t.pointerPosition > t.degree {
		return t.resetPointer()
	}

	return nil
}

func (t *bPTree[K, V]) resetPointer() error {
	if t.pointerNode.linkNodeRight != nil {
		t.pointerNode = t.pointerNode.linkNodeRight
		t.pointerPosition = 0
		return nil
	}
	return errors.New("this node does not exist")
}

// check if this element exist
func (t *bPTree[K, V]) Search(key key[K, V]) (bool, *key[K, V]) {
	node, index := t.search(key) // index give us back a first larger element than requested

	if index == 0 || node.key[index-1].key != key.key {
		return false, nil
	}

	return true, &node.key[index-1]
}

// find working leaf
func (t *bPTree[K, V]) searchLeaf(key key[K, V]) {
	t.currentNode = t.root

	for {
		nextIndex := t.currentNode.search(key)
		t.add(nextIndex)

		if nextNode := t.nextNode(nextIndex); nextNode == nil {
			break
		} else {
			t.currentNode = nextNode
		}
	}
}

// split all full parents node
func (t *bPTree[K, V]) splitParent() {
	for t.currentNode != nil && t.degree == t.currentNode.pointer {
		t.newParent()

		newNode := t.split(parentCode)

		t.updateState(newNode)

		//update parents
		for i := 0; i < newNode.pointer+1; i++ {
			index := t.middleKey + 1 + i
			t.currentNode.key[index].nextNode.parent = newNode
		}

		t.nextParent()
	}
}

// next node in pyramid
func (t *bPTree[K, V]) nextParent() {
	t.currentNode = t.currentNode.parent
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

	if t.currentNode.pointer == t.degree {
		t.newParent()

		newNode := t.split(leafCode)

		t.updateState(newNode)
	}

	t.nextParent()
}

// update parent, new spliced node
func (t *bPTree[K, V]) updateState(newNode *node[K, V]) {
	indexToUpdate := t.next()

	parentKeyIndex := t.middleKey
	parentNode := t.currentNode.parent
	parentNode.insertKey(t.currentNode.key[parentKeyIndex], indexToUpdate)

	parentNode.key[indexToUpdate+1].updateNextNode(newNode)
	parentNode.key[indexToUpdate].updateNextNode(t.currentNode)

	removedPointers := len(t.currentNode.key[:parentKeyIndex]) + t.degree%2
	t.currentNode.pointer -= removedPointers

	t.currentNode.link(newNode)
}

type node[K string | int | float64, V any] struct {
	pointer       int
	parent        *node[K, V]
	linkNodeRight *node[K, V]
	linkNodeLeft  *node[K, V]
	key           []key[K, V]
}

// constructor for new node
func newNode[K string | int | float64, V any](degree int) *node[K, V] {
	return &node[K, V]{
		key: make([]key[K, V], degree+1),
	}
}

// make empty key
func (n *node[K, V]) emptyKey(position int) {
	n.key[position] = key[K, V]{}
}

// search returns the index where the specified key should be inserted in the sorted keys array.
func (n *node[K, V]) search(key key[K, V]) int {
	for i, currentKey := range n.key[:n.pointer] {
		if key.key < currentKey.key {
			return i
		}
	}

	return n.pointer
}

// link current node with next one
func (n *node[K, V]) link(node *node[K, V]) {
	if n.linkNodeRight != nil {
		node.linkNodeRight = n.linkNodeRight //right connection
		n.linkNodeRight.linkNodeLeft = node  //left connection
	}

	n.linkNodeRight = node //right connection
	node.linkNodeLeft = n  //left connection
}

// appends all key to current node
func (n *node[K, V]) appendKeys(key []key[K, V], position, code int) {
	copy(n.key[position:], key)
	n.pointer += len(key) - code
}

func (n *node[K, V]) increasePointer() {
	n.pointer++
}

// just append one key
func (n *node[K, V]) appendKey(key key[K, V], position int) {
	n.key[position] = key
}

// insert key on special position
func (n *node[K, V]) insertKey(key key[K, V], position int) {
	copy(n.key[position+1:], n.key[position:])
	n.appendKey(key, position)
	n.increasePointer()
}

type key[K string | int | float64, V any] struct {
	key      K
	value    V
	nextNode *node[K, V]
}

func NewKey[K string | int | float64, V any](realKey K, value V) key[K, V] {
	return key[K, V]{
		key:   realKey,
		value: value,
	}
}

func (k *key[K, V]) updateNextNode(n *node[K, V]) {
	k.nextNode = n
}

type stack struct {
	stack   []int
	current int
}

func newStack(capacity int) stack {
	return stack{
		stack: make([]int, capacity),
	}
}

func (s *stack) add(index int) {
	s.stack[s.current] = index
	s.current++
}

func (s *stack) next() int {
	if s.current > 0 {
		s.current--
		return s.stack[s.current]
	}

	return 0
}

func (s *stack) clear() {
	s.current = 0
}

//======================================================================================00

// for testing nothing special
// number of key
// number of replication key
func (t *bPTree[K, V]) AllRight() (int, int) {
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

func (t *bPTree[K, V]) AllLeft() (int, int) {
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
