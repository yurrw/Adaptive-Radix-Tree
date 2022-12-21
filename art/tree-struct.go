package main

type (
	leafNode struct {
		key   []byte
		value interface{}
	}
	data struct {
		prefix    []byte
		prefixLen int
		size      int
	}

	innerNode struct {
		data
		nodeType int
		keys     []byte
		children []*Node
	}

	Node struct {
		innerNode *innerNode
		leaf *leafNode
	}

	Tree struct {
		root *Node
		size uint64
	}
)

func NewTree() *Tree {
	return &Tree{root: nil, size: 0}
}

func (n *Node) IsLeaf() bool { return n.leaf != nil }

func newNode4() *Node {
	innerNode := &innerNode{
		nodeType: Node4,
		keys:     make([]byte, Node4Max),
		children: make([]*Node, Node4Max),
		data: data{
			prefix: make([]byte, MaxPrefixLen),
		},
	}
	return &Node{innerNode: innerNode}
}

func newNode16() *Node {
	innerNode := &innerNode{
		nodeType: Node16,
		keys:     make([]byte, Node16Max),
		children: make([]*Node, Node16Max),
		data: data{
			prefix: make([]byte, MaxPrefixLen),
		},
	}

	return &Node{innerNode: innerNode}
}