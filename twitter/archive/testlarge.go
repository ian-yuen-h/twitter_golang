package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

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


// func main(){
// 	TestAllRequestsLarge()
// 	// TestAllRequestsXtraLarge()
// }


const evenParity = 1
const oddParity = 0

func getParity(numbers []int, parity int) []int {

	parityNums := make([]int, len(numbers)/2)
	var idx int
	for _, num := range numbers {
		if (parity == evenParity && num%2 == 0) ||
			(parity == oddParity && num%2 != 0) {
			parityNums[idx] = num
			idx++
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(parityNums)))
	return parityNums
}
func runAllRequests(threads, blockSize string, postInfo []int) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "twitter.go", threads, blockSize)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("<runTwitter>: error in getting stdout pipe: Contact Professor Samuels, if see this message.")
	}
	// //homebrew
	// var stdBuffer bytes.Buffer
	// mw := io.MultiWriter(os.Stdout, &stdBuffer)

	// cmd.Stderr = mw
	// //homebrew

	stdin, errIn := cmd.StdinPipe()
	if errIn != nil {
		fmt.Println("<runTwitter>: error in getting stdin pipe: Contact Professor Samuels, if see this message.")
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("<cmd.Start> error in executing Test: Contact Professor Samuels, if see this message.")
	}
	inDone := make(chan bool)
	outDone := make(chan bool)

	/***** First Wave: Add all posts and random Contains ***/
	wave1Done := make(chan bool)
	doneRequest := _TestDoneRequest{"DONE"}
	requestsAdd, responsesAdd, addIdx := createAdds(postInfo, 0)
	successSlice := make([]bool, len(postInfo))
	requestsContains, responsesContains, containsIdx := createContains(postInfo, successSlice, addIdx)

	/**** Second Wave: Remove all the even numbers *****/
	evenPosts := getParity(postInfo, evenParity)
	oddPosts := getParity(postInfo, oddParity)
	wave2Done := make(chan bool)
	removePosts := evenPosts
	successSliceRemove := make([]bool, len(removePosts))
	for idx, _ := range removePosts {
		successSliceRemove[idx] = true
	}
	requestsRemoves, responsesRemoves, removeIdx := createRemoves(removePosts, successSliceRemove, containsIdx)

	/**** Third Wave: Check that evens are removed and remove all the odds posts **/
	wave3Done := make(chan bool)
	containsOrder := evenPosts
	successSliceContains := make([]bool, len(containsOrder))
	requestsContains2, responsesContains2, containsIdx2 := createContains(containsOrder, successSliceContains, removeIdx)

	removePosts2 := oddPosts
	successSliceRemove2 := make([]bool, len(removePosts2))
	for idx, _ := range removePosts2 {
		successSliceRemove2[idx] = true
	}
	requestsRemoves2, responsesRemoves2, removeIdx2 := createRemoves(removePosts2, successSliceRemove2, containsIdx2)

	/**** Fourth Wave: check the feed is empty **/
	wave4Done := make(chan bool)
	order2 := []int{}
	requestFeed, responseFeedExpected, _ := createFeed(order2, removeIdx2)

	go func() {
		encoder := json.NewEncoder(stdin)
		for idx := 0; idx < len(postInfo); idx++ {
			requestAdd := requestsAdd[idx]
			if err := encoder.Encode(&requestAdd); err != nil {
				fmt.Println("<AllRequests> add cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
			}
		}
		fmt.Println("here 1")
		for idx := addIdx; idx < (addIdx + len(postInfo)); idx++ {
			requestContains := requestsContains[idx]
			if err := encoder.Encode(&requestContains); err != nil {
				fmt.Println("<AllRequests> request cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
			}
		}
		fmt.Println("here 2")
		<-wave1Done
		for idx := containsIdx; idx < (containsIdx + len(postInfo)/2); idx++ {
			requestRemove := requestsRemoves[idx]
			if err := encoder.Encode(&requestRemove); err != nil {
				fmt.Println("<AllRequests> remove cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
			}
		}
		fmt.Println("here 2.5")
		<-wave2Done

		for idx := removeIdx; idx < (removeIdx + len(postInfo)/2); idx++ {
			requestContains := requestsContains2[idx]
			if err := encoder.Encode(&requestContains); err != nil {
				fmt.Println("<AllRequests> contains cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
			}
		}
		fmt.Println("here 3")
		for idx := containsIdx2; idx < (containsIdx2 + len(postInfo)/2); idx++ {
			requestRemove := requestsRemoves2[idx]
			if err := encoder.Encode(&requestRemove); err != nil {
				fmt.Println("<AllRequests> remove cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
			}
		}
		fmt.Println("here 4")
		<-wave3Done
		if err := encoder.Encode(&requestFeed); err != nil {
			fmt.Println("<AllRequests> request cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
		}
		fmt.Println("here 5")
		<-wave4Done
		if err := encoder.Encode(&doneRequest); err != nil {
			fmt.Println("<AllRequests> done cmd.encode error in executing Test: Contact Professor Samuels, if see this message.")
		}
		fmt.Println("here 6")
		inDone <- true
	}()

	go func() {
		decoder := json.NewDecoder(stdout)
		var count int
		for {
			var response _TestNormalResponse
			if count < ((len(postInfo) * 2) + ((len(postInfo) / 2) * 3)) {
				if err := decoder.Decode(&response); err != nil {
					break
				}
			}
			if count >= 0 && count < (len(postInfo)*2) {
				if value, ok := responsesAdd[int(response.Id)]; ok {
					if value.Id != response.Id || value.Success != response.Success {
						fmt.Println("Add Request & Response id and success fields do not match. Got(Id=%v,Success=%v), Expected(Id=%v,Success=%v)",
							response.Id, response.Success, value.Id, value.Success)
					}
					count++
				} else if value, ok := responsesContains[int(response.Id)]; ok {
					if value.Id != response.Id {
						fmt.Println("Contains Request & Response id fields do not match. Got(%v), Expected(%v)",
							response.Id, value.Id)
					}
					count++
				} else {
					fmt.Println("Received an invalid id back from twitter.go. We only added ids [0,%v) but got(id):%v", len(postInfo), response.Id)
				}
				if count == (len(postInfo) * 2) {
					wave1Done <- true
					fmt.Println("wave 1 all received")
				}
			} else if count >= (len(postInfo)*2) && count < ((len(postInfo)*2)+(len(postInfo)/2)) {
				if value, ok := responsesRemoves[int(response.Id)]; ok {
					if value.Id != response.Id || value.Success != response.Success {
						fmt.Println("Remove Request & Response id and success fields do not match. Got(Id=%v,Success=%v), Expected(Id=%v,Success=%v)",
							response.Id, response.Success, value.Id, value.Success)
					}
					count++
				} else {
					fmt.Println("Received an invalid id back from twitter.go. We only removed the even ids but got(id):%v", response.Id)
				}
				if count == ((len(postInfo) * 2) + (len(postInfo) / 2)) {
					wave2Done <- true
					fmt.Println("wave 2 all received")
				}
			} else if count >= (len(postInfo)*2+(len(postInfo)/2)) && count < (len(postInfo)*2+(len(postInfo)/2)*3) {
				if value, ok := responsesContains2[int(response.Id)]; ok {
					if value.Id != response.Id || value.Success != response.Success {
						fmt.Println("Contains Request & Response id and success fields do not match. Got(Id=%v,Success=%v), Expected(Id=%v,Success=%v)",
							response.Id, response.Success, value.Id, value.Success)
					}
					count++
				} else if value, ok := responsesRemoves2[int(response.Id)]; ok {
					if value.Id != response.Id || value.Success != response.Success {
						fmt.Println("Remove Request & Response id and success fields do not match. Got(Id=%v,Success=%v), Expected(Id=%v,Success=%v)",
							response.Id, response.Success, value.Id, value.Success)
					}
					count++
				} else {
					fmt.Println("Received an invalid id back from twitter.go. We only removed the odd ids but got(id):%v", response.Id)
				}
				if count == (len(postInfo)*2 + (len(postInfo)/2)*3) {
					wave3Done <- true
					fmt.Println("wave 3 all received")
				}
			} else {
				var responseFeed _TestFeedResponse
				if err := decoder.Decode(&responseFeed); err != nil {
					fmt.Printf("Got an Error\n")
					break
				}
				if responseFeed.Id != responseFeedExpected.Id {
					fmt.Println("Feed Request & Response id fields do not match. Got(%v), Expected(%v)",
						responseFeed.Id, responseFeedExpected.Id)
				} else {
					if len(responseFeedExpected.Feed) != len(responseFeed.Feed) {
						fmt.Println("Feed Response number of posts not equal to each other. Got(%v), Expected(%v)",
							len(responseFeed.Feed), len(responseFeedExpected.Feed))
					}
					count++
				}
				wave4Done <- true
				fmt.Println("wave 4 all received")
				break
			}
		}
		if count != (len(postInfo)*2+(len(postInfo)/2)*3)+1 {
			fmt.Println("Did not receive the right amount of Add&Feed acknowledgements. Got:%v, Expected:%v", count, (len(postInfo)*2+(len(postInfo)/2)*3)+1)
		}
		outDone <- true
		fmt.Println("all received, Hallelujah")
	}()
	<-inDone
	<-outDone

	if err := cmd.Wait(); err != nil {
		fmt.Println("The automated test timed out. You may have a deadlock, starvation issue and/or you did not implement" +
			" the necessary code for passing this test.")
		fmt.Println("error:", err)
	}

	fmt.Println("exit func")
	// //homebrew
	// log.Println(stdBuffer.String())
	// //homebrew
}


func TestAllRequestsLarge() {
	threads := "32"
	blockSize := "1024"
	posts := generateSlice(25000)
	rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })
	runAllRequests(threads, blockSize, posts)
}

func TestAllRequestsXtraLarge() {
	threads := "32"
	blockSize := "10000"
	posts := generateSlice(75000)
	rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })
	runAllRequests(threads, blockSize, posts)
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