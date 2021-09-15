// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import (
	"sync"
)

type ReadWriteLock struct {
	mutex	*sync.Mutex
	cond	*sync.Cond	
	activeReaders	int
	activeWriters	bool
	waitingWriters	int
}

func NewReadWriteLock() *ReadWriteLock{
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)
	newLock := ReadWriteLock{mutex: &mutex, cond: cond}

	return &newLock
}

func (theLock *ReadWriteLock) RLock() {

	defer theLock.mutex.Unlock()
	theLock.mutex.Lock()

	for ((theLock.waitingWriters > 0) || (theLock.activeWriters) || theLock.activeReaders >= 64){
		theLock.cond.Signal()
		theLock.cond.Wait()
	}
	theLock.activeReaders += 1
}

func (theLock *ReadWriteLock) RUnlock() {

	defer theLock.mutex.Unlock()

	theLock.mutex.Lock()

	theLock.activeReaders -= 1

	if theLock.activeReaders < 64{
		theLock.cond.Signal()
	} else if theLock.activeReaders == 0 {
		theLock.cond.Signal()
	}
}

func (theLock *ReadWriteLock) Lock() {
	defer theLock.mutex.Unlock()
	theLock.mutex.Lock()

	theLock.waitingWriters += 1

	for ((theLock.activeReaders > 0) || (theLock.activeWriters)){
		theLock.cond.Signal()
		theLock.cond.Wait()
	}

	theLock.waitingWriters -= 1

	theLock.activeWriters = true
}

func (theLock *ReadWriteLock) Unlock() {
	defer theLock.mutex.Unlock()

	theLock.mutex.Lock()

	theLock.activeWriters = false

	theLock.cond.Signal()
}