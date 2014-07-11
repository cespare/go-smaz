// Package trie is an implementation of a trie (prefix tree) data structure mapping []bytes to ints. It
// provides a small and simple API for usage as a set as well as a 'Node' API for walking the trie.
package trie

// A Node represents a logical vertex in the trie structure.
type Node struct {
	branches    [256]*Node
	val int
	terminal    bool
	hasChildren bool
}


// A Trie is a a prefix tree.
type Trie struct {
	root *Node
}

// New construct a new, empty Trie ready for use.
func New() *Trie {
	return &Trie{
		root: &Node{},
	}
}

// Put inserts the mapping k -> v into the Trie, overwriting any previous value. It returns true if the
// element was not previously in t.
func (t *Trie) Put(k []byte, v int) bool {
	n := t.root
	for _, c := range k {
		next, ok := n.Walk(c)
		if !ok {
			next = &Node{}
			n.branches[c] = next
			n.hasChildren = true
		}
		n = next
	}
	n.val = v
	if n.terminal {
		return false
	}
	n.terminal = true
	return true
}

// Get the value corresponding to k in t, if any.
func (t *Trie) Get(k []byte) (v int, ok bool) {
	n := t.root
	for _, c := range k {
		next, ok := n.Walk(c)
		if !ok {
			return 0, false
		}
		n = next
	}
	if n.terminal {
		return n.val, true
	}
	return 0, false
}

// Root returns the root node of a Trie. A valid Trie (i.e., constructed with New), always has a non-nil root
// node.
func (t *Trie) Root() *Node { return t.root }

// Walk returns the node reached along edge c, if one exists. The ok value indicates whether such a node
// exists.
func (n *Node) Walk(c byte) (next *Node, ok bool) {
	next = n.branches[int(c)]
	return next, (next != nil)
}

// Terminal indicates whether n is terminal in the trie (that is, whether the path from the root to n
// represents an element in the set). For instance, if the root node is terminal, then []byte{} is in the
// trie.
func (n *Node) Terminal() bool { return n.terminal }

// Leaf indicates whether n is a leaf node in the trie (that is, whether it has children). A leaf node must be
// terminal (else it would not exist). Logically, if n is a leaf node then the []byte represented by the path
// from the root to n is not a proper prefix of any element of the trie.
func (n *Node) Leaf() bool { return !n.hasChildren }

// Val gives the value associated with this node. It panics if n is not terminal.
func (n *Node) Val() int {
	if !n.terminal {
		panic("Val called on non-terminal node")
	}
	return n.val
}
