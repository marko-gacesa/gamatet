// Copyright (c) 2020 by Marko Gaćeša

package event

type Slice []Event

var _ Pusher = (*Slice)(nil)

func (q *Slice) Push(e Event) {
	*q = append(*q, e)
}

func (q Slice) Range(f func(Event)) {
	for i := range q {
		f(q[i])
	}
}

func (q Slice) RangeReverse(f func(Event)) {
	for i := len(q) - 1; i >= 0; i-- {
		f(q[i])
	}
}

func (q *Slice) Clear() {
	*q = (*q)[:0]
}

func (q Slice) Size() int {
	return len(q)
}

func (q Slice) Get(idx int) Event {
	return q[idx]
}
