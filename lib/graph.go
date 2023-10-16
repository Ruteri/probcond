package lib

import (
	"fmt"
	"strings"
)

type Node struct {
	Value string
}

type Edge struct {
	Src *Node
	Dst *Node
}

type DAG struct {
	Nodes []*Node
	Edges []Edge
}

func (dag *DAG) Roots() []*Node {
	seenMap := make(map[*Node]struct{})
	for _, edge := range dag.Edges {
		seenMap[edge.Dst] = struct{}{}
	}

	roots := []*Node{}
	for _, node := range dag.Nodes {
		if _, found := seenMap[node]; !found {
			roots = append(roots, node)
		}
	}
	return roots
}

func (dag *DAG) AddNode(value string) *Node {
	node := &Node{
		Value: value,
	}
	dag.Nodes = append(dag.Nodes, node)
	return node
}

func (dag *DAG) AddEdge(source, destination string) error {
	var srcNode, dstNode *Node

	for _, node := range dag.Nodes {
		if node.Value == source {
			srcNode = node
		}
		if node.Value == destination {
			dstNode = node
		}
	}

	if srcNode == nil {
		srcNode = dag.AddNode(source)
	}

	if dstNode == nil {
		dstNode = dag.AddNode(destination)
	}

	edge := Edge{
		Src: srcNode,
		Dst: dstNode,
	}
	dag.Edges = append(dag.Edges, edge)

	return nil
}

func (dag *DAG) NodeChildren(value string) []*Node {
	children := make([]*Node, 0)
	for _, edge := range dag.Edges {
		if edge.Src.Value == value {
			// Append the destination node to the children slice
			children = append(children, edge.Dst)
		}
	}
	return children
}

func (dag *DAG) NodeParents(value string) []*Node {
	parents := make([]*Node, 0)
	// Find the node with the given value
	for _, node := range dag.Nodes {
		if node.Value == value {
			// Find all incoming edges to the node
			for _, edge := range dag.Edges {
				if edge.Dst == node {
					// Append the source node to the parents slice
					parents = append(parents, edge.Src)
				}
			}
			break
		}
	}
	return parents
}

func (dag *DAG) IsRoot(node *Node) bool {
	for _, r := range dag.Roots() {
		if r.Value == node.Value {
			return true
		}
	}

	return false
}

func (dag *DAG) Traverse(cb func(*Node) bool) {
	visited := make(map[*Node]bool)
	for _, node := range dag.Roots() {
		if !dag.traverseNode(node, visited, cb) {
			return
		}
	}
}

func (dag *DAG) traverseNode(node *Node, visited map[*Node]bool, cb func(*Node) bool) bool {
	// Check if the node has been visited before
	if visited[node] {
		return true
	}
	visited[node] = true
	if !cb(node) {
		return false
	}
	// Traverse children recursively
	for _, child := range dag.NodeChildren(node.Value) {
		if !dag.traverseNode(child, visited, cb) {
			return false
		}
	}

	return true
}

func (dag *DAG) Format() string {
	ret := ""
	for _, edge := range dag.Edges {
		ret += fmt.Sprintf("%s -> %s\n", edge.Src.Value, edge.Dst.Value)
	}

	dag.Traverse(func(node *Node) bool {
		parents := Transform(dag.NodeParents(node.Value), func(node *Node) string { return node.Value })
		children := Transform(dag.NodeChildren(node.Value), func(node *Node) string { return node.Value })

		ret += fmt.Sprintf("[%s] -> %s -> [%s]\n", strings.Join(parents, ", "), node.Value, strings.Join(children, ", "))
		return true
	})

	return ret
}

func Transform[T any, R any](vals []T, cb func(T) R) []R {
	rs := []R{}
	for _, t := range vals {
		rs = append(rs, cb(t))
	}
	return rs
}
