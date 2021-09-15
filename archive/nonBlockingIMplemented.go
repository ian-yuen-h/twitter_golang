
package nonBlockQueue2


import (
	"sync/atomic"
	"unsafe"
	// "sync"
	// "fmt"
)


/////////STRUCTS, CONSTRUCTORS
type TaskRequest struct{
	Command string
	Id int
	Body string
	Timestamp int
}

type Node struct{
	value *TaskRequest
	next *Node
	sentry bool
}

func NewNode(item *TaskRequest) *Node{
	newwNode := Node{}
	newwNode.value = item
}

type NBQueue struct {
	header *Node
	tail *Node

}

func NewNBQueue() *NBQueue{
	newQ :=  NBQueue{}


	header := NewNode(nil)		

	header.sentry = true

	newQ.header = header
	newQ.tail = header
}


/////ENQUEUE
func (queue *NBQueue) Enqueue(item *TaskRequest) bool{	

	newNode := NewNode(item)

	for{
		last := queue.tail
		next := last.next
		if (last == queue.tail){
			if next == nil{
				addrOldTailNext := (*unsafe.Pointer)(unsafe.Pointer(&last.next))
				pNext := unsafe.Pointer(next)
				pNew := unsafe.Pointer(newNode)

				addrTail := (*unsafe.Pointer)(unsafe.Pointer(&queue.tail))
				pLast := unsafe.Pointer(last)
				if atomic.CompareAndSwapPointer(addrOldTailNext, pNext, pNew){
					atomic.CompareAndSwapPointer(addrTail, pLast, pNew)
					return true
				}else{
					atomic.CompareAndSwapPointer(addrTail, pLast, pNext)
				}
			}
		}
	}
}


/////DEQUEUE
func (queue *NBQueue) Dequeue() *TaskRequest{
	for{
		first := queue.header
		last := queue.tail
		next := first.next

		if (first == queue.header){
			if (first == last){
				if (next == nil){
					return nil
				}
				addrTail := (*unsafe.Pointer)(unsafe.Pointer(&queue.tail))
				pLast := unsafe.Pointer(last)
				pNext := unsafe.Pointer(next)
				atomic.CompareAndSwapPointer(addrTail, pLast, pNext)
			}else{
				value := next.value
				addrHead := (*unsafe.Pointer)(unsafe.Pointer(&queue.header))
				pFirst := unsafe.Pointer(first)
				pNext := unsafe.Pointer(next)
				if atomic.CompareAndSwapPointer(addrHead, pFirst, pNext){
					return value
				}
			}
		}
	}
}