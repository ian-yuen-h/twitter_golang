package NBqueue

import (
	"sync/atomic"
	"unsafe"
)

//////Structs
type Ticket struct {
	Command string 
	ID int 
	Body string
	Timestamp float64
}

type InterfaceQueue interface {
	Enqueue(data *Ticket)
	Dequeue() *Ticket
	Empty() bool 
}

type queueTrack struct {
	head *Node
	tail *Node
}

type Node struct {
	data *Ticket
	next *Node
}


///////Functions

func newNode(data *Ticket, next *Node) *Node {
	return &Node{data: data, next: next}
}

func NewQueue() InterfaceQueue {
	sentinel := newNode(nil, nil)
	return &queueTrack{head: sentinel, tail: sentinel}
}

func (q *queueTrack) Empty() bool { 
	if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head.next))) == nil {
		return true
	} else {
		return false
	}
}

//implementation from pseudocode in textbook

func (q *queueTrack) Enqueue(data *Ticket) {
	new := newNode(data, nil)
	for {
		last := q.tail
		next := last.next
		if last == q.tail {
			if next == nil {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&last.next)), nil, unsafe.Pointer(new)) {
					atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(last), unsafe.Pointer(new))
					return
				}
			} else {
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(last), unsafe.Pointer(next))
			}
		}
	}
}

func (q *queueTrack) Dequeue() *Ticket {
	for {
		first := q.head
		last := q.tail
		next := first.next 
		if first == q.head {
			if first == last {
				if next == nil {
					return nil
				}
			} else {
				m := next.data
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(first), unsafe.Pointer(next)) {
					return m
				}
			}
		}
	}
}
