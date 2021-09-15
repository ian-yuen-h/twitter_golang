package main

import (
	"context"
	"encoding/json"
	"fmt"
	// "sort"

	//"io/ioutil"
	"math/rand"
	"os/exec"
	"strconv"
	// "testing"
	"time"
	// "io/ioutil"
	// "log"
	// "io"
	// "os"
	// "bytes"
)

func main(){
	TestSimpleFeedRequest()
}

func TestSimpleFeedRequest() {
	numOfThreadsStr := "16"
	blockSizeStr := "100"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "twitter.go", numOfThreadsStr, blockSizeStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("<SimpleFeedRequest>: error in getting stdout pipe: Contact Professor Samuels, if see this message.")
	}

	// var stdBuffer bytes.Buffer
	// mw := io.MultiWriter(os.Stdout, &stdBuffer)

	// cmd.Stdout = mw
	// cmd.Stderr = mw

	stdin, errIn := cmd.StdinPipe()
	if errIn != nil {
		fmt.Println("<SimpleFeedRequest>: error in getting stdin pipe: Contact Professor Samuels, if see this message.")
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("<SimpleFeedRequest> cmd.Start error in executing Test: Contact Professor Samuels, if see this message.")
	}
	inDone := make(chan bool)
	outDone := make(chan bool)

	doneRequest := _TestDoneRequest{"DONE"}

	numbers := []int{}
	request, responseExpected, _ := createFeed(numbers, 0)

	go func() {
		encoder := json.NewEncoder(stdin)
		time.Sleep(1 * time.Millisecond) // Wait a second before sending the feed request
		if err := encoder.Encode(&request); err != nil {
			fmt.Println("<SimpleFeedRequest> Feed cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
		}
		time.Sleep(1 * time.Millisecond) // Wait a second before sending the add request
		if err := encoder.Encode(&doneRequest); err != nil {
			fmt.Println("<SimpleFeedRequest> Done cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
		}
		inDone <- true
	}()

	go func() {
		decoder := json.NewDecoder(stdout)
		var count int
		for {
			// fmt.Println(1)
			var response _TestFeedResponse
			// fmt.Println("received object", response)
			// fmt.Println("received id", response.Id)
			// fmt.Println("expected id", responseExpected.Id)
			if err := decoder.Decode(&response); err != nil {
				// fmt.Println("received object", response)
				// fmt.Println("received object", response.Id)
				// fmt.Println("expected shit", responseExpected)
				// val := assert(responseExpected == response)
				// fmt.Println("bool value:", val)
				fmt.Println(err)
				break
			}
			// fmt.Println("here")
			if response.Id != responseExpected.Id {
				fmt.Println("Feed Request & Response id fields do not match. Got(%v), Expected(%v)",
					response.Id, responseExpected.Id)
			} else {
				if len(response.Feed) != 0 {
					fmt.Println("Feed Request was sent on an empty feed but the response return posts. Got(%v), Expected(%v)",
						len(response.Feed), 0)
				}
				count++
			}
			if count%5 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
		if count != 1 {
			fmt.Println("Did not receive the right amount of Feed Request acknowledgements. Got:%v, Expected:%v", count, 1)
		}
		outDone <- true
		
	}()

	// log.Println(stdBuffer.String())

	<-inDone
	<-outDone
	if err := cmd.Wait(); err != nil {
		fmt.Println("The automated test timed out. You may have a deadlock, starvation issue and/or you did not implement" +
			" the necessary code for passing this test.")
	}
}





//STRUCTS
type _TestAddRequest struct {
	Command   string  `json:"command"`
	Id        int64   `json:"id"`
	Timestamp float64 `json:"timestamp"`
	Body      string  `json:"body"`
}
type _TestRemoveRequest struct {
	Command   string  `json:"command"`
	Id        int64   `json:"id"`
	Timestamp float64 `json:"timestamp"`
}
type _TestContainsRequest struct {
	Command   string  `json:"command"`
	Id        int64   `json:"id"`
	Timestamp float64 `json:"timestamp"`
}
type _TestFeedRequest struct {
	Command string `json:"command"`
	Id      int64  `json:"id"`
}
type _TestDoneRequest struct {
	Command string `json:"command"`
}

type _TestNormalResponse struct {
	Success bool  `json:"success"`
	Id      int64 `json:"id"`
}

type _TestFeedResponse struct {
	Id   int64           `json:"id"`
	Feed []_TestPostData `json:"feed"`
}

type _TestPostData struct {
	Body      string  `json:"body"`
	Timestamp float64 `json:"timestamp"`
}

func generateSlice(size int) []int {

	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	return slice
}


///////
// Auxiliary functions needed for the tests.
//////
func createAdds(numbers []int, idx int) (map[int]_TestAddRequest, map[int]_TestNormalResponse, int) {

	requests := make(map[int]_TestAddRequest)
	responses := make(map[int]_TestNormalResponse)

	for _, number := range numbers {
		numberStr := strconv.Itoa(number)
		request := _TestAddRequest{"ADD", int64(idx), float64(number), numberStr}
		response := _TestNormalResponse{true, int64(idx)}
		requests[idx] = request
		responses[idx] = response
		idx++
	}
	return requests, responses, idx
}
func createContains(numbers []int, successes []bool, idx int) (map[int]_TestContainsRequest, map[int]_TestNormalResponse, int) {

	requests := make(map[int]_TestContainsRequest)
	responses := make(map[int]_TestNormalResponse)

	for i, number := range numbers {
		request := _TestContainsRequest{"CONTAINS", int64(idx), float64(number)}
		response := _TestNormalResponse{successes[i], int64(idx)}
		requests[idx] = request
		responses[idx] = response
		idx++
	}
	return requests, responses, idx
}
func createRemoves(numbers []int, successes []bool, idx int) (map[int]_TestRemoveRequest, map[int]_TestNormalResponse, int) {

	requests := make(map[int]_TestRemoveRequest)
	responses := make(map[int]_TestNormalResponse)

	for i, number := range numbers {
		request := _TestRemoveRequest{"REMOVE", int64(idx), float64(number)}
		response := _TestNormalResponse{successes[i], int64(idx)}
		requests[idx] = request
		responses[idx] = response
		idx++
	}
	return requests, responses, idx
}
func createFeed(numbers []int, idx int) (_TestFeedRequest, _TestFeedResponse, int) {

	postData := make([]_TestPostData, len(numbers))
	request := _TestFeedRequest{"FEED", int64(idx)}

	for i, number := range numbers {
		numberStr := strconv.Itoa(number)
		postData[i] = _TestPostData{numberStr, float64(number)}
	}
	response := _TestFeedResponse{int64(idx), postData}
	return request, response, idx + 1
}