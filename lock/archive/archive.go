// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

//implement condition variable, with mutex, atomics, waitgroup
//coarse grain implementation

//limit readers to 64
//every reader, increment count
//if counter >= 64, set flag

//two locks


import (
	"fmt"
	"sync"
)

type ReadWriteLock struct {
	mutex	*sync.Mutex
	cond	*sync.Cond	
	activeReaders	int
	activeWriters	bool
	waitingReaders	int
	waitingWriters	int
}

func NewReadWriteLock(mutex *sync.Mutex, cond	*sync.Cond) *ReadWriteLock{
	newLock := ReadWriteLock{mutex: mutex, cond: cond}

	return &newLock
}

func (theLock *ReadWriteLock) BeginRead() {
	theLock.mutex.Lock()

	for ((theLock.waitingWriters > 0) || (theLock.activeWriters) || theLock.activeReaders >= 64){
		theLock.cond.Wait()
	}

	theLock.activeReaders += 1

	theLock.mutex.Unlock()

}

func (theLock *ReadWriteLock) EndRead() {
	theLock.mutex.Lock()

	theLock.activeReaders -= 1

	if theLock.activeReaders < 64{
		theLock.cond.Signal()
	} else if theLock.activeReaders == 0 {
		theLock.cond.Broadcast()
	}

	theLock.mutex.Unlock()
}

func (theLock *ReadWriteLock) BeginWrite() {
	theLock.mutex.Lock()

	theLock.waitingWriters += 1

	for ((theLock.activeReaders > 0) || (theLock.activeWriters)){
		theLock.cond.Wait()
	}

	theLock.waitingWriters -= 1

	theLock.activeWriters = true

	theLock.mutex.Unlock()
}

func (theLock *ReadWriteLock) EndWrite() {

	theLock.mutex.Lock()

	theLock.activeWriters = false

	theLock.cond.Broadcast()

	theLock.mutex.Unlock()
}