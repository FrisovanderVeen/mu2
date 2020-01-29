package bot

import (
	"errors"
	"fmt"
	"sync"
)

// ErrOutOfBounds is used when a queue index is out of bounds
var ErrOutOfBounds = errors.New("queue index out of bounds")

type queue struct {
	front *queueItem

	mu sync.RWMutex
}

type queueItem struct {
	next *queueItem
	item VoiceItem
}

func (q *queue) Add(vi VoiceItem) {
	q.mu.Lock()
	defer q.mu.Unlock()
	qi := &queueItem{
		next: nil,
		item: vi,
	}

	if q.front == nil {
		q.front = qi
		return
	}

	x := q.front
	for x.next != nil {
		x = x.next
	}

	x.next = qi
}

func (q *queue) Next() VoiceItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.front == nil {
		return nil
	}
	qi := q.front
	q.front = q.front.next
	return qi.item
}

// len is not thread safe
func (q *queue) len() int {
	len := 0

	for qi := q.front; qi != nil; qi = qi.next {
		len++
	}

	return len
}

func (q *queue) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.len()
}

// index assumes that i is in the queue
// checks have to be done before calling
// index is not thread safe
func (q *queue) index(i int) *queueItem {
	qi := q.front
	for x := 0; x < i; x++ {
		qi = qi.next
	}

	return qi
}

func (q *queue) Index(i int) *queueItem {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.index(i)
}

func (q *queue) List() []VoiceItem {
	q.mu.RLock()
	defer q.mu.RUnlock()

	list := []VoiceItem{}

	for qi := q.front; qi != nil; qi = qi.next {
		list = append(list, qi.item)
	}

	return list
}

// Reorder will put item a at the position of item b
// item b wil be moved behind a
func (q *queue) Reorder(a, b int) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if a == b {
		return nil
	}
	len := q.len()
	if a > len-1 || a < 0 {
		return fmt.Errorf("%w: %d", ErrOutOfBounds, a)
	} else if b > len-1 || b < 0 {
		return fmt.Errorf("%w: %d", ErrOutOfBounds, b)
	}

	if a > b {
		x := q.index(a - 1)
		temp := x.next
		x.next = x.next.next
		x = temp

		y := q.index(b - 1)
		x.next = y.next
		y.next = x
		return nil
	}

	y := q.index(b)

	x := q.index(a - 1)
	temp := x.next
	x.next = x.next.next
	x = temp

	x.next = y.next
	y.next = x

	return nil
}

func (q *queue) Remove(i int) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if i < 0 {
		return fmt.Errorf("%w: %d", ErrOutOfBounds, i)
	}
	len := q.len()
	if i > len-1 {
		return fmt.Errorf("%w: %d", ErrOutOfBounds, i)
	}

	if i == 0 {
		q.front = q.front.next
		return nil
	}

	if i == len-1 {
		qi := q.index(i - 1)
		qi.next = nil

		return nil
	}

	qi := q.index(i)

	qi.next = qi.next.next

	return nil
}
