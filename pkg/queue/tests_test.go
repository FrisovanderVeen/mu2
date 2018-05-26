package queue_test

import (
	"fmt"

	"github.com/fvdveen/mu2/pkg/queue"
)

var reorderTests = []struct {
	is  []interface{}
	as  []func(q *queue.Queue) error
	res []interface{}
}{
	{
		is: []interface{}{
			1, 2, 3, 4,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				return q.Reorder(0, 2)
			},
		},
		res: []interface{}{
			2, 3, 1, 4,
		},
	},
	{
		is: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				if err := q.Reorder(0, 2); err != nil {
					return err
				}
				return q.Reorder(5, 8)
			},
		},
		res: []interface{}{
			2, 3, 1, 4, 5, 7, 8, 9, 6, 10,
		},
	},
	{
		is: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				if err := q.Reorder(-1, 2); err != queue.ErrOutOfBounds {
					return fmt.Errorf("expected: %v got: %v", queue.ErrOutOfBounds, err)
				}
				if err := q.Reorder(5, 11); err != queue.ErrOutOfBounds {
					return fmt.Errorf("expected: %v got: %v", queue.ErrOutOfBounds, err)
				}
				return nil
			},
		},
		res: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
	},
	{
		is: []interface{}{
			1, 2, 3, 4,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				return q.Reorder(3, 0)
			},
		},
		res: []interface{}{
			4, 1, 2, 3,
		},
	},
	{
		is: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				if err := q.Reorder(3, 0); err != nil {
					return err
				}
				return nil
			},
		},
		res: []interface{}{
			4, 1, 2, 3, 5, 6, 7, 8, 9, 10,
		},
	},
	{
		is: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				return q.Reorder(5, 5)
			},
		},
		res: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
	},
	{
		is: []interface{}{},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				if err := q.Reorder(5, 5); err != queue.ErrOutOfBounds {
					return err
				}
				return nil
			},
		},
		res: []interface{}{},
	},
}

var removeTests = []struct {
	is  []interface{}
	as  []func(q *queue.Queue) error
	res []interface{}
}{
	{
		is: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		as: []func(q *queue.Queue) error{
			func(q *queue.Queue) error {
				if err := q.Remove(10); err != queue.ErrOutOfBounds {
					return fmt.Errorf("expected: %v got: %v", queue.ErrOutOfBounds, err)
				}
				if err := q.Remove(9); err != nil {
					return err
				}

				return nil
			},
		},
		res: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9,
		},
	},
}
