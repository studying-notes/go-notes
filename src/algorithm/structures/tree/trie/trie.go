package trie

type Node struct {
	children map[rune]*Node // map of child nodes
	isWord   bool           // true if this node is the end of a word
	Value    string         // optional value
}

func New() *Node {
	return &Node{children: make(map[rune]*Node)}
}

func (n *Node) Add(word string) {
	if len(word) == 0 {
		n.isWord = true
		return
	}

	r := rune(word[0])
	if _, ok := n.children[r]; !ok {
		n.children[r] = New()
	}

	n.children[r].Add(word[1:])
}

func (n *Node) Has(word string) bool {
	if len(word) == 0 {
		return n.isWord
	}

	r := rune(word[0])
	if _, ok := n.children[r]; !ok {
		return false
	}

	return n.children[r].Has(word[1:])
}

func (n *Node) Remove(word string) bool {
	if len(word) == 0 {
		if n.isWord {
			n.isWord = false
			return true
		}
		return false
	}

	r := rune(word[0])
	if _, ok := n.children[r]; !ok {
		return false
	}

	if n.children[r].Remove(word[1:]) {
		if len(n.children[r].children) == 0 {
			delete(n.children, r)
		}
		return true
	}

	return false
}

func (n *Node) HasPrefix(prefix string) bool {
	if len(prefix) == 0 {
		return true
	}

	r := rune(prefix[0])
	if _, ok := n.children[r]; !ok {
		return false
	}

	return n.children[r].HasPrefix(prefix[1:])
}

func (n *Node) Words() []string {
	return n.words("")
}

func (n *Node) words(prefix string) []string {
	words := make([]string, 0)

	if n.isWord {
		words = append(words, prefix)
	}

	for r, child := range n.children {
		words = append(words, child.words(prefix+string(r))...)
	}

	return words
}

func (n *Node) Prefixes() []string {
	return n.prefixes("")
}

func (n *Node) prefixes(prefix string) []string {
	prefixes := make([]string, 0)

	if len(n.children) > 0 {
		prefixes = append(prefixes, prefix)
	}

	for r, child := range n.children {
		prefixes = append(prefixes, child.prefixes(prefix+string(r))...)
	}

	return prefixes
}
