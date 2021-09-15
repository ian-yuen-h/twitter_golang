package nonBlockQueue


import (
	"sync/atomic"
	"math"
	"time"
	"math/rand"
	"unsafe"
	"sync"
	"fmt"

)

type TaskRequest struct{
	Command string
	Id int
	Body string
	Timestamp int
}

type NBQueue struct {
	header *Node
	tail *Node
	ListSize int32
	nextID int32

	headerID int32
	tailID int32

	//code for testing
	QueueStart [50]int32
	QueueEnd [50]int32
	queueLock *sync.Mutex
	Counter int32
	Done int32
}

type Node struct{
	sentry bool
	id int32					//keep dict outside?
	next *Node
	pred *Node
	data *TaskRequest
}

func NewNode(item *TaskRequest, queue *NBQueue) *Node{
	node := Node{}
	id := atomic.AddInt32(&queue.nextID, 1)	//generate next id
	node.id = id
	node.data = item
	node.sentry = false

	// queue.dictLock.Lock()
	// queue.dict[int32(id)] = int(item)		//store item as key, id as value
	// queue.dictLock.Unlock()
	return &node
}

func NewNBQueue() *NBQueue{
	newList :=  NBQueue{}

	//code for testing
	lock2 := sync.Mutex{}
	newList.queueLock = &lock2

	header := NewNode(nil, &newList)		//errors?
	tail := NewNode(nil, &newList)

	header.id = int32(0)
	tail.id = int32(math.Inf(1))

	newList.headerID = header.id
	newList.tailID = tail.id

	header.sentry = true		//marking as sentry
	tail.sentry = true

	header.next = tail
	tail.pred = header

	newList.header = header
	newList.tail = tail

	newList.ListSize = 0
	newList.nextID = int32(0)

	
	return &newList
}

func (queue *NBQueue) Enqueue(item *TaskRequest) bool{			//pass in pointers
	node := NewNode(item, queue)
	for{
		if queue.TryEnqueue(node){
			atomic.AddInt32(&queue.ListSize, 1)
			// fmt.Println("queuesize:", queue.ListSize)
			// fmt.Println("address enqueued", item)
			return true
		} else{
			randi := rand.Intn(10)			//once fail delay by a certain tim
			time.Sleep(time.Duration(randi % 2) * time.Millisecond)
		}
	}
}

func (queue *NBQueue) TryEnqueue(node *Node) bool{
	end := queue.tail.pred
	addr := (*unsafe.Pointer)(unsafe.Pointer(&queue.tail.pred))
	old := unsafe.Pointer(end)
	new := unsafe.Pointer(node)

	addr2 := (*unsafe.Pointer)(unsafe.Pointer(&end.next))
	old2 := unsafe.Pointer(queue.tail)
	new2 := unsafe.Pointer(node)

	addr3 := (*unsafe.Pointer)(unsafe.Pointer(&node.next))
	new3 := unsafe.Pointer(queue.tail)

	// if end.next.sentry == true{

	// }

	if atomic.CompareAndSwapPointer(addr, old, new){	//changed tail.pred to new node addr
		// if end.next.sentry == false{
		// 	atomic.CompareAndSwapPointer(addr2, old2, new2) //changed oldEnd.next to new node
		// }
		atomic.CompareAndSwapPointer(addr2, old2, new2)
		atomic.StorePointer(addr3, new3)		//new node.next to tail

		fmt.Println("new end:", node)
		fmt.Println("old end:", end)
		fmt.Println("old end's next", end.next)
		//testing code
		// queue.queueLock.Lock()

		// n := atomic.AddInt32(&queue.Counter, 1) -1
		// queue.QueueStart[n] = node.id
		
		// queue.queueLock.Unlock()
		// printStatement := strconv.Itoa(int(node.id))
		// fmt.Println("Queued, ", printStatement)
		return true
	} else{
		return false
	}

}

func (queue *NBQueue) Dequeue() *TaskRequest{
	for{
		returnVal, returnBool := queue.TryDequeue()
		if returnVal!= nil{
			fmt.Println("successfully dequeed") 
			atomic.AddInt32(&queue.ListSize, -1)
			return returnVal
		} else if returnBool == false{
			return nil
		} else {
			randi := rand.Intn(10)			//once fail delay by a certain tim
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
	fmt.Println("magic code")
	// oldFront = *oldFront
	dq := oldFront.data
	if oldFront.id == int32(math.Inf(1)){
		// fmt.Println(15)
		return nil, false		//empty queue
	}

	// newFront := oldFront.next

	addrD := (*unsafe.Pointer)(unsafe.Pointer(&queue.header.next))
	oldD := unsafe.Pointer(oldFront)
	newD := unsafe.Pointer(oldFront.next)
	fmt.Println("pre-atomic")
	fmt.Println("thing dequeued:", oldFront)
	fmt.Println("supposed next:", oldFront.next)
	fmt.Println("next in line:", queue.header.next)

	fmt.Println(54)

	if atomic.CompareAndSwapPointer(addrD, oldD, newD){ //changed queuehead.next, point to deququed node.next
		fmt.Println("successful dequeu pointer swap")
		fmt.Println("post-atomic")
		fmt.Println("thing dequeued:", oldFront)
		fmt.Println("supposed next:", oldFront.next)
		fmt.Println("next in line:", queue.header.next)

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
		return dq, true
	}else{
		return nil, true
	}
}