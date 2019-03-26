package iradix

import (
	"bytes"
	"sort"
)

// Seek is used together with Next to seek/iterate over the tree.
// When prefix is empty or can't be found in the tree, it works
// similar to Iterator/Next. When prefix is non-empty, (*Seeker).Next,
// unlike (*Iterator).Next, returns not the element seeked to, but
// the one after it, and keeps iterating until the last element
// of the tree, in order, including those elements that are after
// but not under prefix.
func (n *Node) Seek(prefix []byte) *Seeker {
	search := prefix
	p := &pos{n: n}
	for {
		// Check for key exhaustion
		if len(search) == 0 {
			return &Seeker{p}
		}

		num := len(n.edges)
		idx := sort.Search(num, func(i int) bool {
			return n.edges[i].label >= search[0]
		})
		p.current = idx
		if idx < len(n.edges) {
			n = n.edges[idx].node
			if bytes.HasPrefix(search, n.prefix) && len(n.edges) > 0 {
				search = search[len(n.prefix):]
				p.current++
				p = &pos{n: n, prev: p}
				continue
			}
		}
		p.current++
		return &Seeker{p}
	}
}

// Seeker is used to iterate over the tree nodes.
type Seeker struct {
	*pos
}

type pos struct {
	n       *Node
	current int
	prev    *pos
	isLeaf  bool
}

// Next returns the next node from the tree. See Seek for details.
func (s *Seeker) Next() (k []byte, v interface{}, ok bool) {
	if s.current >= len(s.n.edges) {
		if s.prev == nil {
			return nil, nil, false
		}
		s.pos = s.prev
		return s.Next()
	}

	edge := s.n.edges[s.current]
	s.current++
	if edge.node.leaf != nil && !s.isLeaf {
		s.isLeaf = true
		s.current--
		return edge.node.leaf.key, edge.node.leaf.val, true
	}
	s.isLeaf = false
	s.pos = &pos{n: edge.node, prev: s.pos}
	return s.Next()
}
