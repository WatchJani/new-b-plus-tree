package b_plus_tree

const NotExist int = -1

type Key struct {
	Value    int
	NextNode *Node
}

// *fix that (maybe need the pointer)
func NewKey(value int) Key {
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

// Search returns the index where the specified key should be inserted in the sorted keys array.
func (n *Node) Search(key Key) int {
	for i, currentKey := range n.Key[:n.Pointer] {
		if key.Value < currentKey.Value {
			return i
		}
	}

	return n.Pointer
}

// find the leaf, parent index if exist, and find position in node to insert new element
func (n *Node) SearchLeaf(key Key) (*Node, int, int) {
	currentNode, nextIndex, prevuesIndex := n, n.Search(key), NotExist

	for currentNode.Key[nextIndex].NextNode != nil {
		prevuesIndex = nextIndex
		currentNode = currentNode.Key[nextIndex].NextNode
		nextIndex = currentNode.Search(key)
	}

	return currentNode, nextIndex, prevuesIndex
}

func (n *Node) UpdateValue(key Key, position int) {
	n.Key[position] = key
}

func (n *Node) PointerIncrease() {
	n.Pointer++
}

func (n *Node) AppendKey(key Key, position int) {
	n.UpdateValue(key, position)
	n.PointerIncrease()
}

func (n *Node) InsertKey(key Key, position int) {
	copy(n.Key[position+1:], n.Key[position:])
	n.AppendKey(key, position)
}

type BPTree struct {
	Memory []Node
	Root   *Node
	Degree int
	ProcessingNode
}

type ProcessingNode struct {
	CurrentNode        *Node
	ParentKeyPosition  int
	CurrentKeyPosition int
}
