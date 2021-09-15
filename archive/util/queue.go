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


type Queue struct {
	enqLock	*sync.Mutex
	deqLock	*sync.Mutex
	queueSize int
	queueLimit int
	condEnq *CondVar
	condDeq *CondVar
	header *Node
	tail *Node
}

type Node struct{
	sentry bool
	id int
	next *Node
}

func NewQueue(enqLock *sync.Mutex, deqLock *sync.Mutex, queueLimit int, condEnq *CondVar, condDeq *CondVar) *Queue{
	newQueue := Queue{enqLock: enqLock, deqLock: deqLock, queueLimit: queueLimit, condEnq: condEnq, condDeq: condDeq}

	header := Node{sentry: true}
	newQueue.header = &header

	newQueue.tail = &header

	return &newQueue
}

func (queue *Queue) Enqueue(id int) {
	queue.enqLock.Lock()

	for (queue.queueSize >= queue.queueLimit){
		queue.condEnq.Wait()
	}

	newNode := Node{id: id, sentry: false}
	queue.tail.next = &newNode			//what if this queue tail gets suddenly dequeued
	queue.tail = &newNode


	queue.enqLock.Unlock()

}

func (queue *Queue) Dequeue() {

}

func (queue *Queue) QueueSize() {

}

func (queue *Queue) QueueEmpty() {

}