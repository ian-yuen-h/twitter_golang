package main

import (
	"proj1/nonBlockQueue2"
	"sync"
	"fmt"
	// "sync/atomic"
)

func worker(group *sync.WaitGroup, queue *nonBlockQueue2.NBQueue, id int){

	item := nonBlockQueue2.TaskRequest{Id: id}
	queue.Enqueue(&item)
	// fmt.Println("queued", id)

	// printStatement := strconv.Itoa(int(queue.ListSize))
	// fmt.Println("queue_size:", printStatement)
	group.Done()

}

func worker2(group *sync.WaitGroup, queue *nonBlockQueue2.NBQueue){

	queue.Dequeue()
	// fmt.Println("queued", id)
	// printStatement := strconv.Itoa(int(result))
	// fmt.Println("output:", result)
	group.Done()

}


func main() {
	threads := 30
	queue := nonBlockQueue2.NewNBQueue()

	var group sync.WaitGroup

	for i := 0; i < threads; i++ {
		group.Add(1)
		go worker(&group, queue, i)
	}
	group.Wait()

	for i := 0; i < threads; i++ {
		group.Add(1)
		go worker2(&group, queue)
	}
	group.Wait()

	fmt.Println("received", queue.QueueStart)
	fmt.Println("received popped:", queue.QueueEnd)

}
  