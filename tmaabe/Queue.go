package tmaabe

import (
    "errors"
    //"fmt"
)

type Queue struct {
    Elements []*TreeNode
}

// New creates a new array queue.
func NewQueue() *Queue {
    return &Queue{}
}

// Size returns the number of Elements in the queue.
func (s *Queue) Size() int {
    return len(s.Elements)
}

// Empty returns true or false whether the queue has zero Elements or not.
func (s *Queue) Empty() bool {
    return len(s.Elements) == 0
}

// Clear clears the queue.
func (s *Queue) Clear() {
    s.Elements = make([]*TreeNode, 0, 10)
}

// Push adds an Element to the queue.
func (s *Queue) Add(e *TreeNode) {
    s.Elements = append(s.Elements, e)
}

// Pop fetches the top Element of the queue and removes it.
func (s *Queue) Poll() (*TreeNode, error) {
    if s.Empty() {
        return nil, errors.New("Pop: the queue cannot be empty")
    }
    result := s.Elements[0]
    s.Elements = s.Elements[1:len(s.Elements)]
    return result, nil
}

// Top returns the top of Element from the queue, but does not remove it.
func (s *Queue) Top() (*TreeNode, error) {
    if s.Empty() {
        return nil, errors.New("Top: queue cannot be empty")
    }
    return s.Elements[len(s.Elements)-1], nil
}
