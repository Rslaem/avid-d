package tmaabe

import (
    "errors"
    //"fmt"
)

type ArrayStack struct {
    Elements []string
}

// New creates a new array stack.
func NewStack() *ArrayStack {
    return &ArrayStack{}
}

// Size returns the number of Elements in the stack.
func (s *ArrayStack) Size() int {
    return len(s.Elements)
}

// Empty returns true or false whether the stack has zero Elements or not.
func (s *ArrayStack) Empty() bool {
    return len(s.Elements) == 0
}

// Clear clears the stack.
func (s *ArrayStack) Clear() {
    s.Elements = make([]string, 0, 10)
}

// Push adds an Element to the stack.
func (s *ArrayStack) Push(e string) {
    s.Elements = append(s.Elements, e)
}

// Pop fetches the top Element of the stack and removes it.
func (s *ArrayStack) Pop() (string, error) {
    if s.Empty() {
        return "", errors.New("Pop: the stack cannot be empty")
    }
    result := s.Elements[len(s.Elements)-1]
    s.Elements = s.Elements[:len(s.Elements)-1]
    return result, nil
}

// Top returns the top of Element from the stack, but does not remove it.
func (s *ArrayStack) Top() (string, error) {
    if s.Empty() {
        return "", errors.New("Top: stack cannot be empty")
    }
    return s.Elements[len(s.Elements)-1], nil
}
