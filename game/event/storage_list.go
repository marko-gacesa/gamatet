// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package event

type List struct {
	first *eventsNode
	last  *eventsNode
}

type eventsNode struct {
	e    Event
	prev *eventsNode
	next *eventsNode
}

var _ Pusher = (*List)(nil)

func (q *List) Push(e Event) {
	n := &eventsNode{
		e: e,
	}

	if q.last != nil {
		n.prev = q.last
		q.last.next = n
	}
	q.last = n

	if q.first == nil {
		q.first = n
	}
}

func (q *List) Dequeue() (e Event) {
	if q.first == nil {
		return
	}

	e = q.first.e

	if second := q.first.next; second != nil {
		second.prev = nil
		q.first.next = nil
		q.first = second
	}

	return
}

func (q List) Range(f func(e Event)) {
	for curr := q.first; curr != nil; curr = curr.next {
		f(curr.e)
	}
}

func (q List) RangeReverse(f func(e Event)) {
	for curr := q.last; curr != nil; curr = curr.prev {
		f(curr.e)
	}
}

func (q *List) Clear() {
	for curr, next := q.first, q.first; curr != nil; curr = next {
		next = curr.next
		curr.prev = nil
		curr.next = nil
	}
	q.first = nil
	q.last = nil
}

func (q List) IsEmpty() bool {
	return q.first == nil
}

func (q List) Size() int {
	n := 0
	for curr := q.first; curr != nil; curr = curr.next {
		n++
	}
	return n
}

func (q List) Get(idx int) Event {
	for curr, i := q.first, 0; curr != nil; curr = curr.next {
		if i == idx {
			return curr.e
		}
	}
	return nil
}

func (q List) FirstEquals(event Event) bool {
	return q.first == nil || q.first.e.Equals(event)
}
