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
	"log"
	"io"
	"os"
	"bytes"
)

func main(){
	TestSimpleAddAndFeedRequest()
}

func TestSimpleAddAndFeedRequest() {
	numOfThreadsStr := "3"
	blockSizeStr := "1"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "twitter.go", numOfThreadsStr, blockSizeStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("<SimpleAddAndFeedRequest>: error in getting stdout pipe: Contact Professor Samuels, if see this message.")
	}

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	stdin, errIn := cmd.StdinPipe()
	if errIn != nil {
		fmt.Println("<SimpleAddAndFeedRequest>: error in getting stdin pipe: Contact Professor Samuels, if see this message.")
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("<SimpleAddAndFeedRequest> cmd.Start error in executing Test: Contact Professor Samuels, if see this message.")
	}
	inDone := make(chan bool)
	outDone := make(chan bool)
	addDone := make(chan bool)

	doneRequest := _TestDoneRequest{"DONE"}
	postInfo := []int{1, 2, 18, 9, 8, 20, 16, 10, 6, 14, 17, 15, 19, 5, 13, 11, 7, 4, 3, 12}
	requestsAdd, responsesAdd, addIdx := createAdds(postInfo, 0)
	order := []int{20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	requestFeed, responseFeedExpected, _ := createFeed(order, addIdx)

	go func() {
		encoder := json.NewEncoder(stdin)
		for idx := 0; idx < len(postInfo); idx++ {
			requestAdd := requestsAdd[idx]
			if err := encoder.Encode(&requestAdd); err != nil {
				fmt.Println("<SimpleAddAndFeedRequest> add cmd.encode error in executing test: Contact Professor Samuels, if see this message.")
			}
		}
		<-addDone
		if err := encoder.Encode(&requestFeed); err != nil {
			fmt.Println("<SimpleAddAndFeedRequest> feed cmd.encode error in executing test: Contact Professor Samuels, if see this message.")
		}
		if err := encoder.Encode(&doneRequest); err != nil {
			fmt.Println("<SimpleAddAndFeedRequest> done cmd.encode errror in executing test: Contact Professor Samuels, if see this message.")
		}
		inDone <- true
	}()

	go func() {
		decoder := json.NewDecoder(stdout)
		var count int
		for {
			var response _TestNormalResponse
			if err := decoder.Decode(&response); err != nil {
				break
			}
			if value, ok := responsesAdd[int(response.Id)]; ok {
				if value.Id != response.Id || value.Success != response.Success {
					fmt.Println("Add request & response id and success fields do not match. Got(Id=%v,Success=%v), Expected(Id=%v,Success=%v)",
						response.Id, response.Success, value.Id, value.Success)
				}
				count++
			} else {
				fmt.Println("Received an invalid id back from twitter.go. We only added ids [0,30] but got(id):%v", response.Id)
			}
			if count == 20 {
				addDone <- true
				var responseFeed _TestFeedResponse
				if err := decoder.Decode(&responseFeed); err != nil {
					break
				}
				if responseFeed.Id != responseFeedExpected.Id {
					fmt.Println("Feed request & response id fields do not match. Got(%v), Expected(%v)",
						responseFeed.Id, responseFeedExpected.Id)
				} else {
					if len(responseFeedExpected.Feed) != len(responseFeed.Feed) {
						fmt.Println("Feed response number of posts not equal to each other. Got(%v), Expected(%v)",
							len(responseFeed.Feed), len(responseFeedExpected.Feed))
					}
					for idx, post := range responseFeed.Feed {
						if post.Body != responseFeedExpected.Feed[idx].Body || post.Timestamp != responseFeedExpected.Feed[idx].Timestamp {
							fmt.Println("Feed response post data does not match. This is checking that the order returned is correct. Got(Body:%v, TimeStamp:%v), Expected(Body:%v, TimeStamp:%v)",
								post.Body, post.Timestamp, responseFeedExpected.Feed[idx].Body, responseFeedExpected.Feed[idx].Timestamp)
						}
					}
					count++
				}
			}
		}
		if count != len(requestsAdd)+1 {
			fmt.Println("Did not receive the right amount of Add&Feed acknowledgements. Got:%v, Expected:%v", count, len(requestsAdd)+1)
		}
		outDone <- true
	}()

	log.Println(stdBuffer.String())

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