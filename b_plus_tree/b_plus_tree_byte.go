package b_plus_tree

import (
	"bytes"
	"math"
)

// insert or replace new value to tree
func (t *BPTreeByte[V]) InsertByte(key []byte, value V) {
	realKey := NewKeyByte(key, value) //make real key

	if t.root == nil { //Check if root exist
		t.root = t.createNodeByte() //create new root
	}

	t.searchLeafByte(realKey.key) //find leaf
	t.appendToLeaf(realKey)       //add to leaf

	t.splitParent() //check if will parent split

	t.clear()
}

func NewBPTreeByte[V any](capacity, degree int) *BPTreeByte[V] {
	numberOfNode := capacity / degree
	treeHeight := math.Log2(float64(capacity))

	return &BPTreeByte[V]{
		memory:    make([]*nodeByte[V], 0, numberOfNode),
		degree:    degree,
		middleKey: degree / 2,
		processingNodeByte: processingNodeByte[V]{
			stack: newStack(int(treeHeight)),
		},
	}
}

func (t *BPTreeByte[V]) splitParent() {
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

func (t *BPTreeByte[V]) insertOrUpdate(key KeyByte[V]) {
	positionInsert := t.next()

	if t.shouldUpdate(positionInsert, key) {
		t.updateExistingKey(positionInsert, key)
	} else {
		t.insertNewKey(positionInsert, key)
	}
}

func (t *BPTreeByte[V]) insertNewKey(positionInsert int, key KeyByte[V]) {
	t.currentNode.insertKey(key, positionInsert)
}

// insert key on special position
func (n *nodeByte[V]) insertKey(key KeyByte[V], position int) {
	copy(n.key[position+1:], n.key[position:])
	n.appendKey(key, position)
	n.increasePointer()
}

func (n *nodeByte[V]) increasePointer() {
	n.pointer++
}

func (n *nodeByte[V]) appendKey(key KeyByte[V], position int) {
	n.key[position] = key
}

func (t *BPTreeByte[V]) updateExistingKey(positionInsert int, key KeyByte[V]) {
	t.currentNode.key[positionInsert-1] = key
}

func (t BPTreeByte[V]) shouldUpdate(positionInsert int, key KeyByte[V]) bool {
	return positionInsert > 0 && bytes.Compare(t.currentNode.key[positionInsert-1].key, key.key) == 0
} //t.currentNode.key[positionInsert-1].key == key.key

func (t *BPTreeByte[V]) appendToLeaf(key KeyByte[V]) {
	t.insertOrUpdate(key)

	if t.currentNode.pointer == t.degree {
		t.newParent()

		newNode := t.split(leafCode)

		t.updateState(newNode)
	}

	t.nextParent()
}

func (t *BPTreeByte[V]) nextParent() {
	t.currentNode = t.currentNode.parent
}

// update parent, new spliced node
func (t *BPTreeByte[V]) updateState(newNode *nodeByte[V]) {
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

func (n *nodeByte[V]) link(node *nodeByte[V]) {
	if n.linkNodeRight != nil {
		node.linkNodeRight = n.linkNodeRight //right connection
		n.linkNodeRight.linkNodeLeft = node  //left connection
	}

	n.linkNodeRight = node //right connection
	node.linkNodeLeft = n  //left connection
}

func (k *KeyByte[V]) updateNextNode(n *nodeByte[V]) {
	k.nextNode = n
}

func (t *BPTreeByte[V]) split(code int) *nodeByte[V] {
	newNode := t.createNodeByte()

	start, end := t.middleKey+code, t.degree+code
	sliceOfKey := t.currentNode.key[start:end]

	newNode.appendKeys(sliceOfKey, 0, code)
	newNode.parent = t.currentNode.parent

	return newNode
}

func (n *nodeByte[V]) appendKeys(key []KeyByte[V], position, code int) {
	copy(n.key[position:], key)
	n.pointer += len(key) - code
}

func (t *BPTreeByte[V]) newParent() {
	if t.currentNode.parent == nil {
		t.currentNode.parent = t.createNodeByte()
		t.root = t.currentNode.parent
	}
}

func (t *BPTreeByte[V]) createNodeByte() *nodeByte[V] {
	newNode := newNodeByte[V](t.degree)
	t.memory = append(t.memory, newNode)

	newNode.emptyKeyByte(0)

	return newNode
}

func (n *nodeByte[V]) emptyKeyByte(position int) {
	n.key[position] = KeyByte[V]{}
}

func newNodeByte[V any](degree int) *nodeByte[V] {
	return &nodeByte[V]{
		key: make([]KeyByte[V], degree+1),
	}
}

func (t *BPTreeByte[V]) searchLeafByte(key []byte) {
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

func (t *BPTreeByte[V]) nextNode(nextIndex int) *nodeByte[V] {
	return t.currentNode.key[nextIndex].nextNode
}

func (n *nodeByte[V]) search(key []byte) int {
	for i, currentKey := range n.key[:n.pointer] {
		if bytes.Compare(key, currentKey.key) == -1 { //key < currentKey.key
			return i
		}
	}

	return n.pointer
}

func NewKeyByte[V any](realKey []byte, value V) KeyByte[V] {
	return KeyByte[V]{
		key:   realKey,
		value: value,
	}
}

type BPTreeByte[V any] struct {
	memory    []*nodeByte[V]
	root      *nodeByte[V]
	degree    int
	middleKey int
	processingNodeByte[V]
	keyPointerByte[V]
}

type keyPointerByte[V any] struct {
	pointerPosition int
	pointerNode     *nodeByte[V]
}

type processingNodeByte[V any] struct {
	currentNode *nodeByte[V]
	stack       //have history of indexes
}

type nodeByte[V any] struct {
	pointer       int
	parent        *nodeByte[V]
	linkNodeRight *nodeByte[V]
	linkNodeLeft  *nodeByte[V]
	key           []KeyByte[V]
}

type KeyByte[V any] struct {
	key      []byte
	value    V
	nextNode *nodeByte[V]
}
