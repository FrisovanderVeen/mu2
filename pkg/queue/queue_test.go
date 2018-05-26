package queue_test

import (
	"reflect"
	"testing"

	"github.com/fvdveen/mu2/pkg/queue"
)

func TestReorder(t *testing.T) {
	for _, test := range reorderTests {
		q := queue.New()
		q.PushBack(test.is...)
		for _, f := range test.as {
			err := f(q)
			if err != nil {
				t.Error(err)
			}
		}
		var res []interface{}
		for q.Length() > 0 {
			res = append(res, q.PopFront())
		}

		if len(res) == 0 && len(test.res) == 0 {
		} else if e := reflect.DeepEqual(res, test.res); !e {
			t.Errorf("expected: %v got: %v", test.res, res)
		}
	}
}

func TestRemove(t *testing.T) {
	for _, test := range removeTests {
		q := queue.New()
		q.PushBack(test.is...)
		for _, f := range test.as {
			err := f(q)
			if err != nil {
				t.Error(err)
			}
		}
		var res []interface{}
		for q.Length() > 0 {
			res = append(res, q.PopFront())
		}

		if len(res) == 0 && len(test.res) == 0 {
		} else if e := reflect.DeepEqual(res, test.res); !e {
			t.Errorf("expected: %v got: %v", test.res, res)
		}
	}
}
