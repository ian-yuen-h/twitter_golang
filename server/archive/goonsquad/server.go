package server

import (
	"encoding/json"
	"proj1/feed"
	"proj1/nonBlockQueue2"
	"sync"
	// "sync/atomic"
	// "fmt"
	"math"
	"os"
)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
	BlockSize      int // Represents the maximum number of tasks a consumer
	// can process at any given point in time.
}

//Run starts up the twitter server based on the configuration
//information provided and only returns when the server is fully
// shutdown.


func Run(config Config) {

	Queue := nonBlockQueue2.NewNBQueue()

	Feed := feed.NewFeed()

	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)

	var group sync.WaitGroup

	go producer(config.Decoder, Queue, cond, &group, config.BlockSize)			//start listening
	group.Add(1)

	if config.Mode == "p"{
		for i := 0; i < config.ConsumersCount; i++ {
			go consumer(Queue, &Feed, cond, &config, &group)
			group.Add(1)
		}
	}else if config.Mode == "s"{
		config.BlockSize = 10
		go consumer(Queue, &Feed, cond, &config, &group)
		group.Add(1)
	}
	group.Wait()

}

//request, producer, consumer
//no queue
// as long as task it not done, keep going
//decode, 

func producer(dec *json.Decoder, queue *nonBlockQueue2.NBQueue, cond *sync.Cond, group *sync.WaitGroup, blocksize int){
	// fmt.Println(1)
	checking := &queue.Done
	val := int32(math.Inf(1))

	for (*checking == val) {


		var item nonBlockQueue2.TaskRequest
		err := dec.Decode(&item)

		if err == nil {
			// fmt.Println("received", item)
			queue.Enqueue(&item)
			cond.Signal()
			// SignalConsumers(queue, cond, blocksize)
		}else if (err != nil) {
			break
		}
		if item.Command == "DONE"{
			break
		}
	  }
	group.Done()
}


func consumer(queue *nonBlockQueue2.NBQueue, ss *feed.Feed, cond *sync.Cond, config *Config, group *sync.WaitGroup) {

	var checking *int32
	checking = &queue.Done
	val := int32(math.Inf(1))
	for (*checking == val){

		tasksOnHand := make([]*nonBlockQueue2.TaskRequest, 0)

		counter := 0

		for (*checking == val) {
			if counter <= config.BlockSize{
				taskObj := queue.Dequeue()
				if (taskObj != nil){ 	//non empty queue
					tasksOnHand = append(tasksOnHand, taskObj)
					counter +=1
					continue
				} else if counter >= 0{
					break
				} else {
					cond.Wait()
				}
			}
		}

		if !(*checking == val){
			break
		}

		if counter >0{
			for _, taskObj := range tasksOnHand {
				command := taskObj.Command

				switch command{
				case "ADD":
					feed.Feed.Add(*ss, taskObj.Body, float64(taskObj.Timestamp))
			
					response := responseClient{Success: true, Id: taskObj.Id}
					
					if err := config.Encoder.Encode(&response); err != nil {
						_= err
					}
			
				case "REMOVE":
					returnVal := feed.Feed.Remove(*ss, float64(taskObj.Timestamp))

					var response responseClient
			
					if returnVal == true{
						response = responseClient{Success: true, Id: taskObj.Id}
					}else{
						response = responseClient{Success: false, Id: taskObj.Id}
					}
					if err := config.Encoder.Encode(&response); err != nil {
						_= err
					}
			
				case "CONTAINS":
					returnVal := feed.Feed.Contains(*ss, float64(taskObj.Timestamp))
					
					var response responseClient
					if returnVal == true{
						response = responseClient{Success: true, Id: taskObj.Id}
					}else{
						response = responseClient{Success: false, Id: taskObj.Id}
					}

					if err := config.Encoder.Encode(&response); err != nil {
						_= err
					}
			
				case "FEED":
					results := feed.Prints(*ss)
					results.Id = int64(taskObj.Id)

					enc := json.NewEncoder(os.Stdout)
					err := enc.Encode(&results)
					if err != nil{
						_= err
					}
				case "DONE":
					queue.Done = int32(1)
					cond.Broadcast()
				}	
			}
		}
	}
	group.Done()
}


type responseClient struct{
	Success bool
	Id int
}


