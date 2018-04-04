// Package singlyLinkedList implements a singly linked list
package list

// SinglyLinkedList is a singly linked list implementation
type SinglyLinkedList struct {
	front *Node

	length int
}

// Init initializes an empty list
func (s *SinglyLinkedList) Init() *SinglyLinkedList {
	s.length = 0
	return s
}

// New returns an initialized list
func New() *SinglyLinkedList {
	return new(SinglyLinkedList).Init()
}

// Front returns the first node in list s
func (s *SinglyLinkedList) Front() NodeService {
	return s.front
}

// Back returns the last node in list s
func (s *SinglyLinkedList) Back() NodeService {
	currentNode := s.front
	for currentNode != nil && currentNode.next.(*Node) != nil {
		currentNode = currentNode.next.(*Node)
	}
	return currentNode
}

// Append appends node n to list s
func (s *SinglyLinkedList) Append(n NodeService) {
	if s.front == nil {
		s.front = n.(*Node)
	} else {
		currentNode := s.front

		for currentNode.next.(*Node) != nil {
			currentNode = currentNode.next.(*Node)
		}

		currentNode.next = n
	}

	s.length++
}

// Prepend prepends node n to list s
func (s *SinglyLinkedList) Prepend(n NodeService) {
	if s.front == nil {
		s.front = n.(*Node)
	} else {
		n.next = s.front.(*Node)
		s.front = n.(*Node)
	}

	s.length++
}

// InsertBefore inserts node insert before node before in list s
func (s *SinglyLinkedList) InsertBefore(insert *Node, before *Node) {
	if s.front == before {
		insert.next = before
		s.front = insert
		s.length++
	} else {
		currentNode := s.front
		for currentNode != nil && currentNode.next.(*Node) != nil && currentNode.next.(*Node) != before {
			currentNode = currentNode.next.(*Node)
		}

		if currentNode.next.(*Node) == before {
			insert.next = before
			currentNode.next = insert
			s.length++
		}
	}
}

// InsertAfter inserts node insert after node after in list s
func (s *SinglyLinkedList) InsertAfter(insert *Node, after *Node) {
	currentNode := s.front
	for currentNode != nil && currentNode != after && currentNode.next.(*Node) != nil {
		currentNode = currentNode.next.(*Node)
	}

	if currentNode == after {
		insert.next = after.next.(*Node)
		after.next = insert
		s.length++
	}
}

// Remove removes node n from list s
func (s *SinglyLinkedList) Remove(n *Node) {
	if s.front == n {
		s.front = n.next.(*Node)
		s.length--
	} else {
		currentNode := s.front

		// search for node n
		for currentNode != nil && currentNode.next.(*Node) != nil && currentNode.next.(*Node) != n {
			currentNode = currentNode.next.(*Node)
		}

		// see if current's next node is n
		// if it's not n, then node n wasn't found in list s
		if currentNode.next.(*Node) == n {
			currentNode.next = currentNode.next.(*Node).next.(*Node)
			s.length--
		}
	}
}

// RemoveBefore removes node before node before
func (s *SinglyLinkedList) RemoveBefore(before *Node) {
	if s.front != nil && s.front != before {
		if s.front.next.(*Node) == before {
			s.front = before
		} else {
			currentNode := s.front
			for currentNode.next.(*Node).next.(*Node) != nil && currentNode.next.(*Node).next.(*Node) != before {
				currentNode = currentNode.next.(*Node)
			}
			if currentNode.next.(*Node).next.(*Node) == before {
				currentNode.next = before
			}
		}
	}
}

// RemoveAfter removes node after node after
func (s *SinglyLinkedList) RemoveAfter(after *Node) {
	if s.front != nil && s.front.next.(*Node) != nil {
		currentNode := s.front
		for currentNode != after && currentNode.next.(*Node) != nil {
			currentNode = currentNode.next.(*Node)
		}

		if currentNode == after {
			currentNode.next = currentNode.next.(*Node).next.(*Node)
		}
	}
}

// GetAtPos returns the node at index in list s
func (s *SinglyLinkedList) GetAtPos(index int) *Node {
	currentNode := s.front
	count := 0
	for count < index && currentNode != nil && currentNode.next.(*Node) != nil {
		currentNode = currentNode.next.(*Node)
		count++
	}

	if count == index {
		return currentNode
	} else {
		return nil
	}
}

// Find returns the node with matching value or nil if not found
func (s *SinglyLinkedList) Find(value interface{}) *Node {
	currentNode := s.front
	for currentNode != nil && currentNode.Value != value && currentNode.next.(*Node) != nil {
		currentNode = currentNode.next.(*Node)
	}

	return currentNode
}

// Length returns the amount of nodes in list s
func (s *SinglyLinkedList) Length() int {
	return s.length
}
