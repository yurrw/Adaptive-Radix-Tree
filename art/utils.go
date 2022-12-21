package main

import "bytes"

// utils
func Terminate(key []byte) []byte {
	// adiciona 0 no final
	index := bytes.Index(key, []byte{0})
	if index < 0 {
		key = append(key, byte(0))
	}
	return key
}

func (n *innerNode) minSize() int {
	switch n.nodeType {
	case Node4:
		return Node4Min
	case Node16:
		return Node16Min
	case Node48:
		return Node48Min
	case Node256:
		return Node256Min
	default:
	}
	return 0
}


func (n *Node) minimum() *Node {
	in := n.innerNode

	switch n.Type() {
	case Node4, Node16:
		return in.children[0].minimum()

	case Node48:
		i := 0
		for in.keys[i] == 0 {
			i++
		}

		child := in.children[in.keys[i]-1]

		return child.minimum()

	case Node256:
		i := 0
		for in.children[i] == nil {
			i++
		}
		return in.children[i].minimum()

	case Leaf:
		return n
	}

	return n
}

func (n *leafNode) IsMatch(key []byte) bool {
	return bytes.Equal(n.key, key)
}


func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func copyBytes(dest []byte, src []byte, numBytes int) {
	for i := 0; i < numBytes && i < len(src) && i < len(dest); i++ {
		dest[i] = src[i]
	}
}


func (n *innerNode) maxSize() int {
	switch n.nodeType {
	case Node4:
		return Node4Max
	case Node16:
		return Node16Max
	case Node48:
		return Node48Max
	case Node256:
		return Node256Max
	default:
	}
	return 0
}


func (n *innerNode) isFull() bool { return n.size == n.maxSize() }



func (n *Node) Type() int { // mudar nome para is leaf
	if n.innerNode != nil {
		return n.innerNode.nodeType
	}
	if n.leaf != nil {
		return Leaf
	}
	return -1
}


func replace(old, new *innerNode) {
	*old = *new
}

func replaceNode(old, new *Node) {
	*old = *new
}

func (n *innerNode) copy(src *innerNode) {
	n.data = src.data
}