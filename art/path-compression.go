package main

//To implement path compression in this radix tree implementation, you could add a new method to the innerNode type that combines the prefixes of all the child nodes into a single prefix for the current node. This method could be called whenever a new child is added to an inner node, and it would update the prefix of the current node to include the common prefix shared by all its children. You could also modify the insert() method to check for common prefixes between the existing keys in a node and the key being inserted, and combine them into a single prefix for the node if possible. This would reduce the depth of the tree and improve its performance.

//Here is an example of how this could be implemented:

