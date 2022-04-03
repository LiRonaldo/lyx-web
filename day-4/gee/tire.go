package gee

import "strings"

type Node struct {
	isWild   bool
	part     string
	partner  string
	children []*Node
}

func (n *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.partner = pattern
		return
	}
	part := parts[height]
	child := n.match(part)
	if child == nil {
		child = &Node{
			part:   part,
			isWild: part[0] == '*' || part[0] == ':'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *Node) search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.partner == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	child := n.matchChildren(part)
	if child != nil {
		for _, children := range child.children {
			result := children.search(parts, height+1)
			if result != nil {
				return result
			}
		}
		return child
	}
	return nil
}
func (n *Node) matchChildren(part string) *Node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
func (n *Node) match(part string) *Node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
