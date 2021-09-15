//non-lock
//enqueue
//dequeue
//queue size/ empty?

package util


import (
	"fmt"
	"sync"
	"cond_var"
)


type NQueue struct {
	enqFlag int
	deqFlag int
	queueSize int
	queueLimit int
	cond *CondVar
	header *Node
	tail *Node
}

type NNode struct{
	sentry bool
	id int
	next *Node
}

func NewNQueue(enqLock *sync.Mutex, deqLock *sync.Mutex, queueLimit int, cond *CondVar) *Queue{
	newQueue := NQueue{enqLock: enqLock, deqLock: deqLock, queueLimit: queueLimit, cond: cond}

	header := Node{sentry: true}
	newQueue.header = &header

	return &newQueue
}

func (queue *NQueue) Enqueue() {
	queue.enqLock.Lock()

	for (queue.queueSize >= queue.queueLimit){
		queue.cond.Wait()
	}

	queue.enqLock.Unlock()

}

func (queue *NQueue) Dequeue() {

}

func (queue *NQueue) QueueSize() {

}

func (queue *NQueue) QueueEmpty() {

}