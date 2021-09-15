package nonBlockQueue


import (
	"sync/atomic"
	"math"
	"time"
	"math/rand"
	"unsafe"
	"fmt"
	"strconv"

)

type NBQueue struct {
	header *Node
	tail *Node
	ListSize int
	dict map[int32]int //dict taking generic pointers?
	nextID int32

	headerID int32
	tailID int32
	dictFlag int32
}

type Node struct{
	sentry bool
	id int32					//keep dict outside?
	marked bool
	next *Node
	lockFlag int32				//trying lock-free
	pred *Node
}

func NewNode(item int, queue *NBQueue) *Node{
	node := Node{}
	id := atomic.AddInt32(&queue.nextID, 1)	//generate next id
	node.id = id
	for atomic.CompareAndSwapInt32(&queue.dictFlag, 0, 1){
	}
	queue.dict[id] = item		//store item as key, id as value
	atomic.StoreInt32(&queue.dictFlag, 0)
	return &node
}

func NewNBQueue() *NBQueue{
	newList :=  NBQueue{}

	newList.dict = make((map[int32]int))

	header := NewNode(0, &newList)		//errors?
	tail := NewNode(0, &newList)

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

func (queue *NBQueue) Enqueue(item int) bool{			//pass in pointers
	node := NewNode(item, queue)
	for{
		if queue.TryEnqueue(node){
			queue.ListSize += 1
			return true
		} else{
			randi := rand.Intn(10)			//once fail delay by a certain tim
			time.Sleep(time.Duration(randi % 2) * time.Millisecond)
		}
	}
}

func (queue *NBQueue) TryEnqueue(node *Node) bool{
	end := queue.tail.pred
	node.next = queue.tail
	addr := (*unsafe.Pointer)(unsafe.Pointer(&queue.tail.pred))
	old := unsafe.Pointer(end)
	new := unsafe.Pointer(node)

	if atomic.CompareAndSwapPointer(addr, old, new){		//change last item next
		end.next = node
		
		printStatement := strconv.Itoa(int(node.id))
		fmt.Println("Queued, ", printStatement)
		return true
	} else{
		return false
	}

}

func (queue *NBQueue) Dequeue() int{
	for{
		returnNodeID := queue.TryDequeue()
		if returnNodeID != int32(math.Inf(1)){
			queue.ListSize -= 1
			return queue.dict[returnNodeID]
		} else{
			return 0	//empty queue
		}

	}
}

func (queue *NBQueue) TryDequeue() int32{
	oldFront := queue.header.next
	dq := oldFront.id
	if oldFront.sentry{
		return int32(math.Inf(1))		//empty queue
	}

	newFront := oldFront.next

	addr := (*unsafe.Pointer)(unsafe.Pointer(&queue.header.next))
	old := unsafe.Pointer(oldFront)
	new := unsafe.Pointer(newFront)

	if atomic.CompareAndSwapPointer(addr, old, new){
		printStatement := strconv.Itoa(int(dq))
		fmt.Println("Dequeued, ", printStatement)
		return oldFront.id
	}else{
		return int32(math.Inf(1))
	}
}