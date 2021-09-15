package main

import (
  "fmt"
  "encoding/json"
  // // "log"
  // "os"
  "reflect"
)

var (
  payloadOne = `{"command":"ADD","id":2,"body":"hello","timestamp": 64}`
  payloadTwo = `{"command": "REMOVE", "id": 3, "timestamp": 43242423}`
)

func main() {

  // sample = '{"command": "ADD", "id": 1, "body": "just setting up my twttr", "timestamp": 43242423}'

  // sample := "{\"command\": \"REMOVE\", \"id\": 1, \"body\": \"just setting up my twttr\", \"timestamp\": 43242423}"
  // fmt.Println(sample)

  // samples := []byte(sample)

  // var m Requests

  // err := json.Unmarshal(samples, &m)

  // fmt.Println(err)


  // obj := addRequest{Command: "ADD", Id:2, Body:"hello", Timestamp: 64}

  // fmt.Println("original", obj)

  // b, _ := json.Marshal(obj)

  // log.Println("marshalled", string(b))

  // var m addRequest

  // err := json.Unmarshal(b, &m)

  // fmt.Println("unmarshalled", m)
  // fmt.Println(err)

  // matchingObj, _ := regexp.Compile("\"command\": \"(?P<Year>[a-zA-Z]+)\", \"id\": (?P<Year>[0-9]+), \"body\": \"(?P<Year>[^\t\n\f\r]+)\", \"timestamp\": (?P<Year>[0-9]+)")
  
  
	// theResult := matchingObj.FindStringSubmatch(sample)

  // command := theResult[1]
  // fmt.Println(command)

  // id := theResult[2]
  // fmt.Println(id)

  // body := theResult[3]
  // fmt.Println(body)

  // timeStamp := theResult[4]
  // fmt.Println(timeStamp)

  // x := []byte(`{"command":"ADD","id":2,"body":"hello","timestamp": 64}`)
  // var q addRequest
  // ee := json.Unmarshal(x, &q)
  // fmt.Println("unmarshalled", q)
  // fmt.Println("irr", ee)

  // dec := json.NewDecoder(os.Stdin)
  // // enc := json.NewEncoder(os.Stdout)
  
  // for {
  //   var list []*Requests
  //   var v Requests
  //   if err := dec.Decode(&v); err != nil {
  //       list = append(list, &v)
  //       fmt.Println(v)
  //   }
  // }

  aa := decodeJSON(payloadOne)
  fmt.Println("return", aa)  // fmt.Println(*aa.Command)

  xType := reflect.TypeOf(aa)
  fmt.Println(xType)

  f, ok := aa.(*addRequest)
  if !ok {
      // baz was not of type *foo. The assertion failed
  }else{
    qType := reflect.TypeOf(f)
    fmt.Println(qType)
    fmt.Println("converted", f)
    fmt.Println(f.Command)
  }
}


func decodeJSON(payload string) interface{} {
  fmt.Println(payload)
  command := struct{
    Command string
  }{}

  piece1 := []byte(payload)

  if err := json.Unmarshal(piece1, &command); err != nil {
    panic(err)
  }

  fmt.Println(command)

  switch command.Command{
	case "ADD":

		var obj addRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "REMOVE":

		var obj removeRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "CONTAINS":

		var obj containsRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "FEED":

		var obj feedRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj

	case "DONE":
		var obj doneRequest
		err := json.Unmarshal(piece1, &obj)

		_ = err

		return &obj
	}
  return nil
}




type Requests interface{
	getCommand() string
	getID()	int
}

type addRequest struct{
	Command string
	Id int
	Body string
	Timestamp int
}


type removeRequest struct{
	Command string
	Id int
	Timestamp int
}


type containsRequest struct{
	Command string
	Id int
	Timestamp int
}

type feedRequest struct{
	Command string
	Id int
}

type doneRequest struct{
	Command string
	Id int
}

func (r addRequest) getCommand() string {
    return r.Command
}
func (r removeRequest) getCommand() string {
    return r.Command
}
func (r containsRequest) getCommand() string {
    return r.Command
}
func (r feedRequest) getCommand() string {
    return r.Command
}
func (r doneRequest) getCommand() string {
    return r.Command
}

func (r addRequest) getID() int {
    return r.Id
}
func (r removeRequest) getID()int {
    return r.Id
}
func (r containsRequest) getID() int {
    return r.Id
}
func (r feedRequest) getID() int {
    return r.Id
}
func (r doneRequest) getID() int {
    return r.Id
}

func commandType(g Requests) string{
	return g.getCommand()
}

func identify(g Requests) int{
	return g.getID()
}