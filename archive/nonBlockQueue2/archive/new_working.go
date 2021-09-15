package nonBlockQueue2


import (
	"sync/atomic"
	"time"
	"math/rand"
	"unsafe"
	// "sync"
	// "fmt"
	"math"

)

type TaskRequest struct{
	Command string
	Id int
	Body string
	Timestamp int
}

type NBQueue struct {
	header *Node
	ending *Node
	ListSize int32
	// nextID int32

	// headerID int32
	// tailID int32

	Counter int32
	Done int32
	ProcessingCount int32
	Available int
}

type Node struct{
	sentry bool
	id int32					//keep dict outside?
	next *Node
	data *TaskRequest
	pred *Node
}

func NewNode(item *TaskRequest, queue *NBQueue) *Node{
	node := Node{}
	// id := atomic.AddInt32(&queue.nextID, 1)	//generate next id
	// node.id = id
	node.data = item
	node.sentry = false
	node.next = nil

	// queue.dictLock.Lock()
	// queue.dict[int32(id)] = int(item)		//store item as key, id as value
	// queue.dictLock.Unlock()
	return &node
}

func NewNBQueue() *NBQueue{
	newList :=  NBQueue{}


	header := NewNode(nil, &newList)		//errors?

	// header.id = int32(0)

	// newList.headerID = header.id		//delete?

	header.sentry = true		//marking as sentry

	header.next = nil
	newList.ending = header


	newList.header = header

	newList.ListSize = 0
	// newList.nextID = int32(0)
	newList.Done = int32(math.Inf(1))
	newList.ProcessingCount = int32(0)

	
	return &newList
}

func (queue *NBQueue) Enqueue(item *TaskRequest) bool{			//pass in pointers
	node := NewNode(item, queue)
	for{
		if queue.TryEnqueue(node){
			// atomic.AddInt32(&queue.ListSize, 1)
			// fmt.Println("enqueued")
			return true
		} else{
			randi := rand.Intn(5)			//change to 4?
			time.Sleep(time.Duration(randi % 2) * time.Millisecond)
		}
	}
}

func (queue *NBQueue) TryEnqueue(node *Node) bool{
	addr := (*unsafe.Pointer)(unsafe.Pointer(&queue.ending.next))
	old := unsafe.Pointer(nil)
	new := unsafe.Pointer(node)

	if atomic.CompareAndSwapPointer(addr, old, new){	//changed tail.pred to new node addr


		// fmt.Println("new end:", node)
		// fmt.Println("old end:", end)
		// fmt.Println("old end's next", end.next)

		return true
	} else{
		return false
	}

}

func (queue *NBQueue) Dequeue() *TaskRequest{
	for{
		returnVal, returnBool := queue.TryDequeue()
		if returnVal!= nil{
			// fmt.Println("successfully dequeed") 
			// atomic.AddInt32(&queue.ListSize, -1)
			return returnVal
		} else if returnBool == false{	//empty
			return nil
		} else {
			randi := rand.Intn(5)			//once fail delay by a certain tim
			time.Sleep(time.Duration(randi % 2) * time.Millisecond)
		}
		_ = returnBool
		// if returnBool == false{
		// 	// fmt.Println(33)
		// 	return nil
		// }
	}
}

func (queue *NBQueue) TryDequeue() (*TaskRequest, bool){
	// oldFront := atomic.LoadPointer(queue.header.next)
	// addr2 := (*unsafe.Pointer)(unsafe.Pointer(&queue.header.next))
	// oldFront := atomic.LoadPointer(addr2)
	// oldFrontV:= unsafe.Pointer(uintptr(oldFront)+offset)
	oldFront := queue.header.next

	if oldFront == nil{
		// fmt.Println(15)
		return nil, false		//empty queue
	}

	newFront := oldFront.next


	addrD := (*unsafe.Pointer)(unsafe.Pointer(&queue.header.next))
	oldD := unsafe.Pointer(oldFront)
	newD := unsafe.Pointer(newFront)


	// fmt.Println(54)

	if atomic.CompareAndSwapPointer(addrD, oldD, newD){ //changed queuehead.next, point to deququed node.next
		// fmt.Println("successful dequeu pointer swap")
		// fmt.Println("post-atomic")
		// fmt.Println("thing dequeued:", oldFront)
		// fmt.Println("supposed next:", oldFront.next)
		// fmt.Println("next in line:", queue.header.next)

		// if (oldFront.next).sentry{
		// 	addr2 := (*unsafe.Pointer)(unsafe.Pointer(&queue.tail.pred))
		// 	old2 := unsafe.Pointer(oldD)
		// 	new2 := unsafe.Pointer(&queue.header)
		// 	atomic.CompareAndSwapPointer(addr2, old2, new2)
		// }
		// queue.queueLock.Lock()
		// n := atomic.AddInt32(&queue.Counter, 1) -1
		// queue.QueueEnd[n] = oldFront.id

		// queue.queueLock.Unlock()
		dq := oldFront.data
		return dq, true
	}else{
		return nil, true
	}
}