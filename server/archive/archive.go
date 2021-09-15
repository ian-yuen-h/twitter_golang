package server

import (
	"encoding/json"
	"proj1/feed"
	"proj1/nonBlockQueue2"
	"sync"
	"sync/atomic"
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


func Run(config Config) bool{
	// fmt.Println("running server")

	// defer os.Exit(3)

	Queue := nonBlockQueue2.NewNBQueue()

	Feed := feed.NewFeed()

	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)

	var group sync.WaitGroup

	go producer(config.Decoder, Queue, cond, &group)			//start listening
	group.Add(1)
	// fmt.Println("spawned producer")

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
	return true

}

//request, producer, consumer
//no queue
// as long as task it not done, keep going
//decode, 

func producer(dec *json.Decoder, queue *nonBlockQueue2.NBQueue, cond *sync.Cond, group *sync.WaitGroup)bool{
	// fmt.Println(1)
	checking := &queue.Done
	val := int32(math.Inf(1))
	// fmt.Println(*checking)
	// fmt.Println(val)
	// fmt.Println((*checking == val))
	for{
		// fmt.Println(*checking)
		if (*checking != val){
			break
		}else{
			var item nonBlockQueue2.TaskRequest
			err := dec.Decode(&item)
			// fmt.Println("action")
			if err == nil {
				queue.Enqueue(&item)
				cond.Signal()
				// fmt.Println("equeued:", item)
				// if (atomic.LoadInt32(&queue.ListSize)) > int32(0){ //tis is alway true?
				// 	cond.Signal()
				// }
			}else if (err != nil) {
				break
			}
			if item.Command == "DONE"{
				break
			}

		}
	  }
	
	// fmt.Println("here")
	group.Done()
	// fmt.Println("producer quits")
	return true
}	

// func queueUp(item *nonBlockQueue.TaskRequest, queue *nonBlockQueue.NBQueue){
// 	queue.Enqueue(item)
// }

func consumer(queue *nonBlockQueue2.NBQueue, ss *feed.Feed, cond *sync.Cond, config *Config, group *sync.WaitGroup) bool{

	//if flag is not xx, keep looping
	//if get done task, set the flag
	checking := &queue.Done
	val := int32(math.Inf(1))
	for (*checking == val){

		// for {
		// 	checking2 := atomic.LoadInt32(&queue.ListSize)
		// 	if checking2 == 0{
		// 		cond.Wait()
		// 	}
		// }
		tasksOnHand := make([]*nonBlockQueue2.TaskRequest, 0)

		counter := 0

		for{
			aa := atomic.LoadInt32(&queue.ListSize)	//can delete atomic add to queue size then
			if (counter<=  config.BlockSize && (aa) > int32(0)){ //if return is nil for queue size?
				taskObj := queue.Dequeue()
				// fmt.Println(12)
				if taskObj != nil{
					tasksOnHand = append(tasksOnHand, taskObj)
					counter +=1
					// fmt.Println("dequeued successfully")
				}
				continue
			} else if counter == config.BlockSize{
				if (aa) > int32(0){
					cond.Signal()
				}
				break
			} else if counter >= 0 && aa == int32(0){
				break
			} else if aa == int32(0){
				cond.Wait()
			}
		}

		// if *checking2 > int32(0){
		// 	cond.Signal()
		// }

		if counter >0{
			// fmt.Println(12)
			for _, taskObj := range tasksOnHand {
				command := taskObj.Command
				//look at task queue
				//depending on what the task is, edit the feed
			
				switch command{
				case "ADD":
					atomic.AddInt32(&queue.ProcessingCount, 1)
					feed.Feed.Add(*ss, taskObj.Body, float64(taskObj.Timestamp))
			
					response := responseClient{Success: true, Id: taskObj.Id}
			
					if resultsOut(response, config.Encoder){
						atomic.AddInt32(&queue.ProcessingCount, -1)
					}

			
				case "REMOVE":
					atomic.AddInt32(&queue.ProcessingCount, 1)
					returnVal := feed.Feed.Remove(*ss, float64(taskObj.Timestamp))

					var response responseClient
			
					if returnVal == true{
						response = responseClient{Success: true, Id: taskObj.Id}
					}else{
						response = responseClient{Success: false, Id: taskObj.Id}
					}
			
					if resultsOut(response, config.Encoder){
						atomic.AddInt32(&queue.ProcessingCount, -1)
					}
			
				case "CONTAINS":
					atomic.AddInt32(&queue.ProcessingCount, 1)
					returnVal := feed.Feed.Contains(*ss, float64(taskObj.Timestamp))
					
					var response responseClient
					if returnVal == true{
						response = responseClient{Success: true, Id: taskObj.Id}
					}else{
						response = responseClient{Success: false, Id: taskObj.Id}
					}
			
					if resultsOut(response, config.Encoder){
						atomic.AddInt32(&queue.ProcessingCount, -1)
					}
			
				case "FEED":
					atomic.AddInt32(&queue.ProcessingCount, 1)
					results := feed.Prints(*ss)
					// response := responseFeed{Id: taskObj.Id, Feed: results}
					results.Id = int64(taskObj.Id)
					// fmt.Println(results.Id)

					// new := modelStruct{id: int(results.Id), feed: results.Feed}
					// fmt.Println("new", new)
					
					// fmt.Println("response crafted")
					// fmt.Println("printing out", v)
					// fmt.Println("results:", results)

					enc := json.NewEncoder(os.Stdout)
					err := enc.Encode(&results)
					// fmt.Println("printed")
					if err != nil{
						// fmt.Println("weird")
						// fmt.Println(err)
						_= err
					}

					// fmt.Println("my results:", results)

					atomic.AddInt32(&queue.ProcessingCount, -1)
					// resultsOutFeed(response, config.Encoder)
					// fmt.Println("done")
			
				case "DONE":
					//check queue size, comapre and swap 0, 0
					//until it is empty, then shut it down, signal with another flag?
					// fmt.Println("command processed")
					// checking := &queue.ListSize
					// checking2 := &queue.ProcessingCount
					// // fmt.Println(*checking2)
					// for ((*checking != int32(0)) && (*checking2 != 0) ){
					// 	// fmt.Println("hung")
					// }
					// // fmt.Println("done done")
					// atomic.StoreInt32(&queue.Done, int32(1))
					// echoDone(config.Encoder)
					// fmt.Println("incremented")
					go setDone(queue)
				}
			
			}
			// checking2 := &queue.ProcessingCount
			// fmt.Println(*checking2)
		}
		counter = 0
	}
	// fmt.Println("consumer quits")

	group.Done()

	return true
	//call encoder, send back response
	//check for task again
}

//decoder, decodes the command bit, then send to decode function
//which adds to the task queue

func setDone(queue *nonBlockQueue2.NBQueue)bool{
	// checking := &queue.ListSize				//queue is empty
	checking2 := &queue.ProcessingCount		//process count is zero
	// fmt.Println(*checking2)
	// for ((*checking != int32(0)) && (*checking2 != 0) ){
	// 	// fmt.Println("hung")
	// }
	for (*checking2 != 0){
		// fmt.Println("hung")
	}
	// fmt.Println("done done")
	atomic.StoreInt32(&queue.Done, int32(1))
	return true
}


type responseClient struct{
	Success bool
	Id int
}

// type responseFeed struct{
// 	Id int
// 	Feed []string
// }

// type modelStruct struct{
// 	id int
// 	feed interface{}
// }

func resultsOut(v responseClient, enc *json.Encoder) bool{
	// fmt.Println(v)
	if err := enc.Encode(&v); err != nil {
		// log.Println(err)
		_= err
	}

	return true
}

// func resultsOutFeed(v responseFeed, Encoder *json.Encoder){
// 	fmt.Println("printing out", v)
// 	err := Encoder.Encode(&v)
// 	// fmt.Println("printed")
// 	if err != nil{
// 		// fmt.Println("weird")
// 		_= err
// 	}
// }

// func echoDone(Encoder *json.Encoder){
// 	for i := 1; i < 10; i++ {
// 		var v string
// 		v = "all done"
// 		err := Encoder.Encode(&v)
// 		if err != nil{
// 			_= err
// 		}
// 	}
// 	// var v string
// 	// 	v = "all done"
// 	// 	err := Encoder.Encode(&v)
// 	// 	if err != nil{
// 	// 		_= err
// 	// 	}
// }


						// fmt.Println("dequeued successfully")
				// }else{
				// 	cond.Wait()
				// 	break
				// }
				// else if counter == config.BlockSize{
				// 	// cond.Signal()
				// 	// if (aa) > int32(0){
				// 	// 	cond.Signal()
				// 	// }
				// 	break
				// }