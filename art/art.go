package main

import (
	"math/bits"
)

func (t *Tree) Insert(key []byte, value interface{}) bool {
	key = Terminate(key)
	updated := t.insert(&t.root, key, value, 0)
	if !updated {
		t.size++
	}
	return updated
}

func replaceNodeRef(oldNode **Node, newNode *Node) {
	*oldNode = newNode
}

func newLeafNode(key []byte, value interface{}) *Node {
	newKey := make([]byte, len(key))
	copy(newKey, key)

	newLeaf := &leafNode{newKey, value}
	return &Node{leaf: newLeaf}
}



func (n *Node) prefixMatchIndex(key []byte, depth int) int {
	idx := 0
	in := n.innerNode
	p := in.prefix

	for ; idx < in.prefixLen && depth+idx < len(key) && key[depth+idx] == p[idx]; idx++ {
		if idx == MaxPrefixLen-1 {
			min := n.minimum()
			p = min.leaf.key[depth:]
		}
	}
	return idx
}


func (n *innerNode) findChild(key byte) **Node {
	if n == nil {
		return nil
	}

	index := n.index(key)

	switch n.nodeType {
	case Node4, Node16, Node48:
		if index >= 0 {
			return &n.children[index]
		}

		return nil

	case Node256:
		child := n.children[key]
		if child != nil {
			return &n.children[key]
		}
	}

	return nil
}
func (n *innerNode) index(key byte) int {
	switch n.nodeType {
	case Node4:
		for i := 0; i < n.size; i++ {
			if n.keys[i] == key {
				return int(i)
			}
		}
		return -1
	case Node16:
		bitfield := uint(0)
		for i := 0; i < n.size; i++ {
			if n.keys[i] == key {
				bitfield |= (1 << i)
			}
		}
		mask := (1 << n.size) - 1
		bitfield &= uint(mask)
		if bitfield != 0 {
			return bits.TrailingZeros(bitfield)
		}
		return -1
	case Node48:
		index := int(n.keys[key])
		if index > 0 {
			return int(index) - 1
		}

		return -1
	case Node256:
		return int(key)
	}

	return -1
}


func (n *innerNode) addChild(key byte, node *Node) {
	if n.isFull() {
		n.grow()
		n.addChild(key, node)
		return
	}

	switch n.nodeType {
	case Node4:
		idx := 0
		for ; idx < n.size; idx++ {
			if key < n.keys[idx] {
				break
			}
		}

		for i := n.size; i > idx; i-- {
			if n.keys[i-1] > key {
				n.keys[i] = n.keys[i-1]
				n.children[i] = n.children[i-1]
			}
		}

		n.keys[idx] = key
		n.children[idx] = node
		n.size += 1

	case Node16:
		idx := n.size
		bitfield := uint(0)
		for i := 0; i < n.size; i++ {
			if n.keys[i] >= key {
				bitfield |= (1 << i)
			}
		}
		mask := (1 << n.size) - 1
		bitfield &= uint(mask)
		if bitfield != 0 {
			idx = bits.TrailingZeros(bitfield)
		}

		for i := n.size; i > idx; i-- {
			if n.keys[i-1] > key {
				n.keys[i] = n.keys[i-1]
				n.children[i] = n.children[i-1]
			}
		}

		n.keys[idx] = key
		n.children[idx] = node
		n.size += 1

	case Node48:
		idx := 0
		for i := 0; i < len(n.children); i++ {
			if n.children[idx] != nil {
				idx++
			}
		}
		n.children[idx] = node
		n.keys[key] = byte(idx + 1)
		n.size += 1

	case Node256:
		n.children[key] = node
		n.size += 1
	}
}


func (n *leafNode) prefixMatchIndex(leaf *leafNode, depth int) int {
	limit := min(len(n.key), len(leaf.key)) - depth

	i := 0
	for ; i < limit; i++ {
		if n.key[depth+i] != leaf.key[depth+i] {
			return i
		}
	}
	return i
}

func (t *Tree) insert(currentRef **Node, key []byte, value interface{}, depth int) bool {
	current := *currentRef
	if current == nil {
		replaceNodeRef(currentRef, newLeafNode(key, value))
		return false
	}

	if current.IsLeaf() {
		if current.leaf.IsMatch(key) {
			current.leaf.value = value
			return true
		}

		currentLeaf := current.leaf
		newLeaf := newLeafNode(key, value)
		limit := currentLeaf.prefixMatchIndex(newLeaf.leaf, depth)

		n4 := newNode4()
		n4in := n4.innerNode
		n4in.prefixLen = limit

		copyBytes(n4in.prefix, key[depth:], min(n4in.prefixLen, MaxPrefixLen))

		depth += n4in.prefixLen

		n4in.addChild(currentLeaf.key[depth], current)
		n4in.addChild(key[depth], newLeaf)
		replaceNodeRef(currentRef, n4)

		return false
	}

	in := current.innerNode
	if in.prefixLen != 0 {
		mIsmatch := current.prefixMatchIndex(key, depth)

		if mIsmatch != in.prefixLen {
			n4 := newNode4()
			n4in := n4.innerNode
			replaceNodeRef(currentRef, n4)
			n4in.prefixLen = mIsmatch

			copyBytes(n4in.prefix, in.prefix, mIsmatch)

			if in.prefixLen < MaxPrefixLen {
				n4in.addChild(in.prefix[mIsmatch], current)
				in.prefixLen -= (mIsmatch + 1)
				copyBytes(in.prefix, in.prefix[mIsmatch+1:], min(in.prefixLen, MaxPrefixLen))
			} else {
				in.prefixLen -= (mIsmatch + 1)
				minKey := current.minimum().leaf.key
				n4in.addChild(minKey[depth+mIsmatch], current)
				copyBytes(in.prefix, minKey[depth+mIsmatch+1:], min(in.prefixLen, MaxPrefixLen))
			}

			newLeafNode := newLeafNode(key, value)
			n4in.addChild(key[depth+mIsmatch], newLeafNode)

			return false
		}
		depth += in.prefixLen
	}

	next := in.findChild(key[depth])
	if next != nil {
		return t.insert(next, key, value, depth+1)
	}

	in.addChild(key[depth], newLeafNode(key, value))
	return false
}

// delete shrink 
func (t *Tree) Delete(key []byte) bool {
	if t.root == nil {
		return false
	}
	key = Terminate(key)
	deleted := t.delete(&t.root, key, 0)
	if deleted {
		t.size--
		return true
	}
	return false
}

func (t *Tree) delete(currentRef **Node, key []byte, depth int) bool {
	current := *currentRef
	if current.IsLeaf() {
		if current.leaf.IsMatch(key) {
			replaceNodeRef(currentRef, nil)
			return true
		}
	} else {
		in := current.innerNode
		if in.prefixLen != 0 {
			mIsmatch := current.prefixMatchIndex(key, depth)
			if mIsmatch != in.prefixLen {
				return false
			}
			depth += in.prefixLen
		}

		next := in.findChild(key[depth])
		if *next != nil {
			if (*next).IsLeaf() {
				leaf := (*next).leaf
				if leaf.IsMatch(key) {
					current.deleteChild(key[depth])
					return true
				}
			}
		}

		return t.delete(next, key, depth+1)
	}
	return false
}

func (n *Node) deleteChild(key byte) {
	in := n.innerNode

	switch n.Type() {
	case Node4, Node16:
		idx := in.index(key)

		in.keys[idx] = 0
		in.children[idx] = nil

		if idx >= 0 {
			for i := idx; i < in.size-1; i++ {
				in.keys[i] = in.keys[i+1]
				in.children[i] = in.children[i+1]
			}

		}

		in.keys[in.size-1] = 0
		in.children[in.size-1] = nil
		in.size -= 1

	case Node48:
		idx := in.index(key)
		if idx >= 0 {
			child := in.children[idx]
			if child != nil {
				in.children[idx] = nil
				in.keys[key] = 0
				in.size -= 1
			}
		}

	case Node256:
		idx := in.index(key)
		child := in.children[idx]
		if child != nil {
			in.children[idx] = nil
			in.size -= 1
		}
	}
	// delete child
	if in.size < in.minSize() {
		n.shrink()
	}
}

