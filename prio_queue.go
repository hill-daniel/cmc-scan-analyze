package main

type node[T any] struct {
	value    T
	next     *node[T]
	previous *node[T]
}

type PriorityQueue[T any] struct {
	size        int
	limit       int
	compareFunc func(c1, c2 T) int
	head        *node[T]
	tail        *node[T]
}

func newPriorityQueue[T any](limit int, compareFunc func(c1, c2 T) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{limit: limit, compareFunc: compareFunc}
}

func (pq *PriorityQueue[T]) Add(value T) {
	if pq.head == nil {
		pq.head = &node[T]{value: value, next: nil, previous: nil}
		pq.size += 1
		return
	}

	cur := pq.head

	for {
		if pq.compareFunc(value, cur.value) <= 0 {
			if cur.next == nil {
				n := &node[T]{value: value, next: nil, previous: cur}
				pq.tail = n
				cur.next = n
				break
			}
			cur = cur.next
			continue
		}

		n := &node[T]{value: value, next: cur, previous: nil}
		if cur == pq.head {
			cur.previous = n
			pq.head = n
		} else {
			n.previous = cur.previous
			cur.previous.next = n
			cur.previous = n
		}

		if pq.tail == nil {
			pq.tail = cur
		}
		break
	}

	if pq.size+1 > pq.limit {
		tail := pq.tail
		previous := tail.previous
		previous.next = nil
		pq.tail = previous
		tail.next = nil
		tail.previous = nil
	} else {
		pq.size += 1
	}
}

func (pq *PriorityQueue[T]) GetAll() []T {
	i := 0
	cur := pq.head
	elements := make([]T, pq.size)
	for {
		if cur == nil {
			break
		}
		elements[i] = cur.value
		cur = cur.next
		i++
	}
	return elements
}
