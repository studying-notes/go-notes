package trie

type DNSCache struct {
	root *Node
}

func NewDNSCache() *DNSCache {
	return &DNSCache{root: New()}
}

func (c *DNSCache) Add(ip, domain string) {
	n := c.root

	for _, r := range ip {
		if _, ok := n.children[r]; !ok {
			n.children[r] = New()
		}
		n = n.children[r]
	}

	n.isWord = true
	n.Value = domain
}

func (c *DNSCache) Get(ip string) string {
	n := c.root

	for _, r := range ip {
		if _, ok := n.children[r]; !ok {
			return ""
		}
		n = n.children[r]
	}

	if n.isWord {
		return n.Value
	}

	return ""
}
