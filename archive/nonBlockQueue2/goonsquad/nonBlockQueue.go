package nonBlockQueue2


import (
	"sync/atomic"
	"unsafe"
	// "sync"
	"math"
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

	//test code
	Id int32
}

func NewNode(item *TaskRequest) *Node{
	newwNode := Node{}
	newwNode.value = item
	return &newwNode
}

type NBQueue struct {
	header *Node
	tail *Node
	Done int32

	//testing code
	Counter int32

	QueueStart []int32
	QueueEnd []int32
}

func NewNBQueue() *NBQueue{
	newQ :=  NBQueue{}


	header := NewNode(nil)		

	header.sentry = true

	newQ.header = header
	newQ.tail = header

	newQ.Done = int32(math.Inf(1))

	return &newQ
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

					//testing code
					id := atomic.AddInt32(&queue.Counter, 1)
					newNode.Id = id
					queue.QueueStart = append(queue.QueueStart, id)
					// fmt.Println(queue.QueueStart)


					// fmt.Println("previous", last)
					// fmt.Println("new ending:", last.next)
					// fmt.Println("new node addr:", &newNode)
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
					// fmt.Println("dequeued", next)
					// fmt.Println("id:", next.Id)
					//testing code
					aa := next.Id
					// fmt.Println(aa)
					queue.QueueEnd = append(queue.QueueEnd, aa)
					// fmt.Println(queue.Qu)
					return value
				}
			}
		}
	}
}