package main

import (
	"fmt"
	"proj1/lock"
	"sync"
)

func worker(group *sync.WaitGroup, queue *nonBlockQueue.NBQueue, id int){

	queue.Enqueue(id)
	// fmt.Println("queued", id)

	// printStatement := strconv.Itoa(int(queue.ListSize))
	// fmt.Println("queue_size:", printStatement)
	group.Done()

}

func worker2(group *sync.WaitGroup, queue *nonBlockQueue.NBQueue){

	queue.Dequeue()
	// fmt.Println("queued", id)
	// printStatement := strconv.Itoa(int(result))
	// fmt.Println("output:", result)
	group.Done()

}


func main() {
	threads := 10
	fmt.Println("Hello World")

	//initialize list?
	//spawn workers
	//spawn readers

	queue := nonBlockQueue.NewNBQueue()

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

}