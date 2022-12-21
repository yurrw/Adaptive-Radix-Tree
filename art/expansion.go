package main

// lazy expansion
func (n *Node) shrink() {
	in := n.innerNode

	switch n.Type() {
	case Node4:
		c := in.children[0]
		if !c.IsLeaf() {
			child := c.innerNode
			currentPrefixLen := in.prefixLen

			if currentPrefixLen < MaxPrefixLen {
				in.prefix[currentPrefixLen] = in.keys[0]
				currentPrefixLen++
			}

			if currentPrefixLen < MaxPrefixLen {
				childPrefixLen := min(child.prefixLen, MaxPrefixLen-currentPrefixLen)
				copyBytes(in.prefix[currentPrefixLen:], child.prefix, childPrefixLen)
				currentPrefixLen += childPrefixLen
			}

			copyBytes(child.prefix, in.prefix, min(currentPrefixLen, MaxPrefixLen))
			child.prefixLen += in.prefixLen + 1
		}

		replaceNode(n, c)

	case Node16:
		n4 := newNode4()
		n4inner_node := n4.innerNode
		n4inner_node.copy(n.innerNode)
		n4inner_node.size = 0

		for i := 0; i < len(n4inner_node.keys); i++ {
			n4inner_node.keys[i] = in.keys[i]
			n4inner_node.children[i] = in.children[i]
			n4inner_node.size++
		}

		replace(n, n4)

	}
}


func (n *innerNode) grow() {
	switch n.nodeType {
	case Node4:
		n16 := newNode16().innerNode
		n16.copy(n)
		for i := 0; i < n.size; i++ {
			n16.keys[i] = n.keys[i]
			n16.children[i] = n.children[i]
		}
		replace(n, n16)
	}
}

