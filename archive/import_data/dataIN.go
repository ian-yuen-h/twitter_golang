package main

import (
	"proj1/requests"
	"encoding/json"
	"bufio"
	"os"
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	result := importData()
	fmt.Println(result)
}


func importData() []string{

	slice1 := make([]string, 3)

	//buffer to read data into
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fileLine := scanner.Text()
		slice1 = append(slice1, fileLine)
	}

	return slice1
}

func unpack(slice1 []string) []*requests.Requests{

	slice2 := make([]*requests.Requests, 3)

	for _, piece := range slice1 {
		command := identify(piece)
		obj := convertFromJSON(command, piece)
		slice2 = append(slice2, &obj)

	}

	return slice2

}

func identify(piece string) string{
	matchingCommand, _ := regexp.Compile("\"command\": \"(?P<Year>[a-zA-Z]+)\"")
  
  
	theResult := matchingCommand.FindStringSubmatch(piece)

	return theResult[1]
}

func convertFromJSON(command string, piece string) *requests.Requests{

	piece1 := []byte(piece)

	switch command{
	case "ADD":

		var obj requests.addRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "REMOVE":

		var obj requests.removeRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "CONTAINS":

		var obj requests.containsRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "FEED":

		var obj requests.feedRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "DONE":
		var obj requests.doneRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj
	}
}

func convertToJSON(obj *requests.Requests) string{

	pBytes, err := json.Marshal(obj)
	_ = err

	return string(pBytes)
}

//decode from json into task objects

//encode from task objects into json

//decoded objects, give to server, add to task queue

//task queue manager: spawn threads

//link queue manger thread to the feed

