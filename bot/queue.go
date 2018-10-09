package bot

import (
	"errors"
	"sync"
)

var (
	// errOutOfBounds is used when an given index is out of bounds
	errOutOfBounds = errors.New("index out of queue bounds")
)

// queue is a FIFO (first in, first out) queue data structure
type queue struct {
	mu sync.RWMutex

	is []Video
}

// newQueue creates a new queue
func newQueue() *queue {
	q := &queue{
		is: make([]Video, 0),
	}
	return q
}

// Length returns the length of the queue
func (q *queue) Length() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.is)
}

// Back returns the last element of the queue
func (q *queue) Back() Video {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.is) == 0 {
		return nil
	}

	return q.is[len(q.is)-1]
}

// PushBack adds the elements to the back of the queue
func (q *queue) PushBack(i ...Video) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.is = append(q.is, i...)
}

// PopBack returns the last element of the queue and removes it
func (q *queue) PopBack() Video {
	q.mu.RLock()
	if len(q.is) == 0 {
		q.mu.RUnlock()
		return nil
	}
	q.mu.RUnlock()
	q.mu.Lock()
	defer q.mu.Unlock()
	l := len(q.is) - 1
	i := q.is[l]
	q.is = q.is[:l]
	return i
}

// Front returns the first element of the queue
func (q *queue) Front() Video {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.is) == 0 {
		return nil
	}

	return q.is[0]
}

// PushFront adds the elements to the front of the queue
func (q *queue) PushFront(i ...Video) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.is = append(i, q.is...)
}

// PopFront returns the first element of the queue and removes it
func (q *queue) PopFront() Video {
	q.mu.RLock()
	if len(q.is) == 0 {
		q.mu.RUnlock()
		return nil
	}
	q.mu.RUnlock()
	q.mu.Lock()
	defer q.mu.Unlock()
	i := q.is[0]
	q.is = q.is[1:]
	return i
}

// Reorder puts element a at elements b position in the queue
func (q *queue) Reorder(a, b int) error {
	q.mu.RLock()
	if a < 0 || b < 0 || a > len(q.is)-1 || b > len(q.is)-1 {
		q.mu.RUnlock()
		return errOutOfBounds
	}
	if a == b {
		q.mu.RUnlock()
		return nil
	}
	q.mu.RUnlock()

	q.mu.Lock()
	defer q.mu.Unlock()

	i := q.is[a]
	q.is = append(q.is[:a], q.is[a+1:]...)
	q.is = append(q.is[:b], append([]Video{i}, q.is[b:]...)...)

	return nil
}

// Copy returns a copy of the queue at that time
func (q *queue) Copy() []Video {
	q.mu.RLock()
	defer q.mu.RUnlock()
	x := make([]Video, q.Length())
	copy(x, q.is)
	return x
}

// Remove removes the item at index i in the queue
func (q *queue) Remove(i int) error {
	q.mu.RLock()
	if i < 0 || i > len(q.is)-1 {
		q.mu.RUnlock()
		return errOutOfBounds
	}
	q.mu.RUnlock()

	q.mu.Lock()
	defer q.mu.Unlock()
	q.is = append(q.is[:i], q.is[i+1:]...)
	return nil
}
