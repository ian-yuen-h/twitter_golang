package cond_var


import (
	"fmt"
	"sync/atomic"
	"proj1/non_block_queue"
)

type CondVar struct {
	queue *non_block_queue.NBQueue
	currThread int32
	nextThread map[int32]bool
}

func NewCondVar(queue *non_block_queue.NBQueue) *CondVar{
	newCond := CondVar{queue: queue}

	return &newCond
}

func (cond *CondVar) Wait() {
	id := atomic.AddInt32(&cond.currThread, 1)
	cond.queue.Enqueue(id)
	cond.nextThread[id] = false
	for cond.nextThread[id] {
	}
}


func (cond *CondVar) Broadcast() {
	for key, _ := range cond.nextThread {
		cond.nextThread[key] = true
	}
}

func (cond *CondVar) Signal() {
	nextID := cond.queue.Dequeue()
	cond.nextThread[nextID] = true
}