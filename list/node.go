// Package singlyLinkedList implements a singly linked list
package list

type NodeService interface {
	Next()
}

// Node is a node within a singly linked list
type Node struct {
	Value interface{}

	next NodeService
}

// Next returns Node n's next node
func (n *Node) Next() NodeService {
	return n.next
}
