func (queue *NBQueue) Enqueue(item *TaskRequest) bool{			//pass in pointers
	node := NewNode(item, queue)
	for{
		if queue.TryEnqueue(node){
			atomic.AddInt32(&queue.ListSize, 1)
			return true
		} else{
			randi := rand.Intn(10)			//once fail delay by a certain tim
			time.Sleep(time.Duration(randi % 2) * time.Millisecond)
		}
	}
}