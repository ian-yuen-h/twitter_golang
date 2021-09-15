package main

import (
	"context"
	"encoding/json"
	"fmt"
	// "sort"

	//"io/ioutil"
	"math/rand"
	"os/exec"
	// "strconv"
	// "testing"
	"time"
	// "io/ioutil"
	"log"
	"io"
	"os"
	"bytes"
)

func main(){
	TestSimpleDone()
}

func TestSimpleDone() {

	numOfThreadsStr := "4"
	blockSizeStr := "1"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "twitter.go", numOfThreadsStr, blockSizeStr)

	//bootleg code
	// stdout, errIn := cmd.StdoutPipe()
	// b, _ := cmd.Output()

	// fmt.Println(1)
	// a, _ := cmd.StdoutPipe()
	// fmt.Println(2)

	// b, _ := ioutil.ReadAll(a)
	// fmt.Println(2)


	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	stdin, errIn := cmd.StdinPipe()
	if errIn != nil {
		fmt.Println("<TestSimpleDone>: stdin error in getting stdin pipe: Contact Professor Samuels, if see this message.")
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("<TestSimpleDone> cmd.Start error in executing test: Contact Professor Samuels, if see this message.")
	}

	done := make(chan bool)
	// fmt.Println(string(stdout))
	// fmt.Println(1)

	// fmt.Println(string(b))

	go func() {

		encoder := json.NewEncoder(stdin)

		request := _TestDoneRequest{"DONE"}

		if err := encoder.Encode(&request); err != nil {
			fmt.Println("<TestSimpleDone> cmd.encode in executing test: Contact Professor Samuels, if see this message.")
		}
		done <- true
	}()

	log.Println(stdBuffer.String())

	<-done
	if err := cmd.Wait(); err != nil {
		fmt.Println("The automated test timed out. You may have a deadlock, starvation issue and/or you did not implement" +
			"the necessary code for passing this test.")
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