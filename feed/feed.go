package feed

import(
	"math"
	"proj1/lock"
	// "fmt"
)

//Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	Add(string, float64)
	Contains(float64) bool
	Remove(float64) bool
	Print() *ResponseFeed
}

//feed is the internal representation of a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation. You can assume the feed will not have duplicate posts
type feed struct {
	start *post // a pointer to the beginning post
	end *post
	lock *lock.ReadWriteLock
}

//post is the internal representation of a post on a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation.
type post struct {
	body      string // the text of the post
	timestamp float64  // Unix timestamp of the post
	next      *post  // the next post in the feed
}


//NewPost creates and returns a new post value given its body and timestamp
func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

//NewFeed creates a empy user feed
func NewFeed() Feed {

	newFeed := feed{}
	start := newPost("nil", 0, nil)
	end := newPost("nil", math.Inf(-1), nil)

	start.next = end

	newFeed.start = start
	newFeed.end = end

	newLock := lock.NewReadWriteLock()
	newFeed.lock = newLock

	return &newFeed
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *feed) Add(body string, timestamp float64) {

	defer f.lock.Unlock()

	f.lock.Lock()
	nPost := newPost(body, timestamp, nil)

	pred := f.start
	curr := pred.next
	// fmt.Println(curr.timestamp)
	// fmt.Println(timestamp)

	for (curr.timestamp > timestamp) {
		pred = curr
		curr = curr.next
	} 
	
	nPost.next = curr		//posts at same time?
	pred.next = nPost
	

}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp float64) bool {

	defer f.lock.Unlock()

	f.lock.Lock()
	pred := f.start
	curr := pred.next

	for (curr.timestamp > timestamp) {
		pred = curr
		curr = curr.next
	}

	if timestamp == curr.timestamp{
		pred.next = curr.next
		return true
	}else{
		return false
	}
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp float64) bool {

	defer f.lock.RUnlock()

	f.lock.RLock()
	pred := f.start
	curr := pred.next

	for (curr.timestamp > timestamp) {
		pred = curr
		curr = curr.next
	}

	if timestamp == curr.timestamp{
		return true
	}else{
		return false
	}
}

func (f *feed) Print() *ResponseFeed{

	defer f.lock.RUnlock()

	f.lock.RLock()
	pred := f.start
	curr := pred.next

	results := make([]Post, 0)
	for (curr.timestamp != math.Inf(-1)){
		temp := Post{Body: curr.body, Timestamp: curr.timestamp}
		results = append(results, temp)
		curr = curr.next
		
	}


	results2 := ResponseFeed{Feed: results}

	return &results2
}



func Adds(r Feed, body string, timestamp float64){
	r.Add(body, timestamp)
}

func Contianss(r Feed, timestamp float64) bool{
	return r.Contains(timestamp)
}

func Removess(r Feed, timestamp float64) bool{
	return r.Remove(timestamp)
}

func Prints(r Feed) *ResponseFeed{
	return r.Print()
}

type ResponseFeed struct{
	Id int64
	Feed []Post
}

type Post struct{
	Body string
	Timestamp float64
}