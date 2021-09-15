
//a queue, is just a lazy linked list, with adding = Enqueue, removing = Dequeue, special type of adding, removing

package lazyLinked


import (
	"fmt"
	"sync/atomic"
	"proj1/cond_var"
	"math"
)

type LazyList struct {
	header *Node
	tail *Node
	listSize int
	dict map[interface{}]int32 //dict taking generic pointers?
	nextID int32

	headerID int32
	tailID int32
}

type Node struct{
	sentry bool
	id int32					//keep dict outside?
	marked bool
	next *Node
	lockFlag int32				//trying lock-free
}

func NewNode(item interface{}, list *LazyList) *Node{
	node := Node{}
	id := atomic.AddInt32(&list.nextID, 1)	//generate next id
	node.id = id
	list.dict[item] = id		//store item as key, id as value
	return &node
}

func NewLazyList() *LazyList{
	newList := LazyList{}

	header := NewNode(nil, &newList)		//errors?
	tail := NewNode(nil, &newList)

	tail.id = math.Inf()

	newList.headerID = header.id
	newList.tailID = tail.id

	header.sentry = true		//marking as sentry
	tail.sentry = true

	header.next = &tail
	newList.header = &header
	newList.tail = &tail


	return &newList
}

func (list *LazyList) Validate(pred *Node, curr *Node) {

	return ((!pred.marked) && (!curr.marked) && (pred.next ==curr))		//pointer?

}

func (list *LazyList) Contains(item interface{}) {

	targetID := list.dict[item]

	inspected := list.header

	for inspected.id < targetID {
		inspected = inspected.next
	}

	return ((inspected.id == targetID) && (!inspected.marked))

}


func (list *LazyList) Add(item interface{}) bool{


	node := NewNode(item, &list)

	pred := list.header
	curr := list.header.next

	defer atomic.StoreInt32(&pred.lockFlag, 1)

	for curr.id < node.id{				// no need for double loop?
		pred = curr
		curr = curr.next
	}


	if atomic.CompareAndSwapInt32(&pred.lockFlag, 0, 1) && atomic.CompareAndSwapInt32(&curr.lockFlag, 0, 1){
		if list.Validate(pred, curr){
			if curr.id == node.id{		//already added?
				atomic.StoreInt32(&pred.lockFlag, 1)
				atomic.StoreInt32(&curr.lockFlag, 1)
				return false
			} else{
				pred.next = &node
				node.next = &curr
				atomic.StoreInt32(&pred.lockFlag, 1)
				atomic.StoreInt32(&curr.lockFlag, 1)
				return true
			}
		}
	}

}

func (list *LazyList) Remove(item interface{}) bool{
	node := NewNode(item, &list)

	pred := list.header
	curr := list.header.next

	defer atomic.StoreInt32(&pred.lockFlag, 1)

	for curr.id < node.id{				// no need for double loop?
		pred = curr
		curr = curr.next
	}

	if atomic.CompareAndSwapInt32(&pred.lockFlag, 0, 1) && atomic.CompareAndSwapInt32(&curr.lockFlag, 0, 1){
		if list.Validate(pred, curr){
			if curr.id == node.id{
				curr.marked = true
				pred.next = curr.next
				atomic.StoreInt32(&pred.lockFlag, 1)
				atomic.StoreInt32(&curr.lockFlag, 1)
				return true
			} else{
				atomic.StoreInt32(&pred.lockFlag, 1)
				atomic.StoreInt32(&curr.lockFlag, 1)
				return false
			}
		}
	}

}