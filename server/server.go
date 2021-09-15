package server

import (
	"encoding/json"
	"proj1/feed"
	"proj1/NBqueue"
	"log"
	"sync"
	"sync/atomic"
)


//////Structs
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

type Response struct {
	Success bool
	ID int
}

type Context struct {
	doneFlag int32

	group sync.WaitGroup
	mutex *sync.Mutex
	cond *sync.Cond

	queue NBqueue.InterfaceQueue
	feed feed.Feed
}

///Functions
func Run(config Config) {
	feed := feed.NewFeed()
	if config.Mode == "p" {
		var mutex sync.Mutex
		condNew := sync.NewCond(&mutex)
		queue := NBqueue.NewQueue()
		ctx := Context{mutex: &mutex, cond: condNew, doneFlag: 1, queue: queue, feed: feed}

		//spawn threads
		ctx.group.Add(config.ConsumersCount)
		for i := 0; i < config.ConsumersCount; i++ {
			go consumer(&config, &ctx)
		} 
		producer(&config, &ctx)
		ctx.group.Wait()
	} else {
		ctx := Context{feed: feed}
		non_parallel(&config, &ctx)
	}
}

func non_parallel(config *Config, ctx *Context) {
	for {
		m := &(NBqueue.Ticket{})
        if err := config.Decoder.Decode(m); err != nil {
            log.Println(err)
            return
        }
		if m.Command == "DONE" {
			return
		}
		doWork(m, config, ctx)
	}
}


func consumer(config *Config, ctx *Context) {
	for {
		tasks := make([]*NBqueue.Ticket, 0) 
		
		for i := 0; i < config.BlockSize; i++ { 
			task := ctx.queue.Dequeue()
			if task == nil {
				break
			} 
			tasks = append(tasks, task)
		}
		for _, task := range tasks {
			doWork(task, config, ctx)
		} 

		ctx.mutex.Lock() 
		for ctx.queue.Empty() && ctx.doneFlag == 1 { 
			ctx.cond.Wait()
		}
		ctx.mutex.Unlock()

		if ctx.queue.Empty() && atomic.LoadInt32(&ctx.doneFlag) == 0 { 
			break
		}

	}
	ctx.group.Done()
}


func producer(config *Config, ctx *Context) {
	for {
		t := &(NBqueue.Ticket{}) 
        if err := config.Decoder.Decode(t); err != nil {
            log.Println(err)
            return
        }
		if t.Command == "DONE" {
			ctx.mutex.Lock() 
			ctx.doneFlag = 0
			ctx.cond.Broadcast() 
			ctx.mutex.Unlock()
			return
		} 
		ctx.queue.Enqueue(t)
		ctx.mutex.Lock() 
		ctx.cond.Signal()
		ctx.mutex.Unlock()
	}
}


func doWork(t *NBqueue.Ticket, config *Config, ctx *Context) {
	switch t.Command {
	case "ADD":
		ctx.feed.Add(t.Body, t.Timestamp)
		config.Encoder.Encode(Response{ID: t.ID, Success: true})
	case "REMOVE":
		success := ctx.feed.Remove(t.Timestamp)
		config.Encoder.Encode(Response{ID: t.ID, Success: success})
	case "CONTAINS":
		success := ctx.feed.Contains(t.Timestamp)
		config.Encoder.Encode(Response{ID: t.ID, Success: success})
	case "FEED": 
		results := ctx.feed.Print()
		results.Id = int64(t.ID)
		config.Encoder.Encode(&results)
	}
}

