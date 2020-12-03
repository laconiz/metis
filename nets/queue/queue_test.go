package queue

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue_Push(t *testing.T) {

	const cap = 10
	queue := New(cap)

	for i := 0; i < cap; i++ {
		if queue.Push(i) != nil {
			t.FailNow()
		}
	}

	if queue.Push(cap) != ErrOverFlow {
		t.FailNow()
	}

	queue.Close()

	for i := 0; i < cap; i++ {
		if queue.Push(cap) != ErrClosed {
			t.FailNow()
		}
	}
}

func TestQueue_Pop(t *testing.T) {

	queue := New(0)

	read, write := int32(0), int32(0)

	var readCh sync.WaitGroup

	readCh.Add(1)
	go func() {
		for {
			if event, closed := queue.Pop(); closed {
				if event != nil {
					t.FailNow()
				}
				readCh.Done()
				break
			} else {
				if event == nil {
					t.FailNow()
				}
				read++
			}
		}
	}()

	now := time.Now()

	var writeCh sync.WaitGroup

	for i := 0; i < 100; i++ {
		writeCh.Add(1)
		go func() {
			for i := 0; i < 1000000; i++ {
				if err := queue.Push(i); err != nil {
					t.FailNow()
				} else {
					atomic.AddInt32(&write, 1)
				}
			}
			writeCh.Done()
		}()
	}

	writeCh.Wait()
	log.Println(time.Since(now))
	queue.Close()
	log.Println(time.Since(now))
	readCh.Wait()
	log.Println(time.Since(now))
	log.Printf("read: %d, write: %d", read, write)
}

func TestQueue_Close(t *testing.T) {

	const cap = 10
	queue := New(cap)

	ch := make(chan bool, 1)

	go func() {
		for {
			if _, closed := queue.Pop(); closed {
				ch <- true
				break
			}
		}
	}()

	if queue.Close() != nil {
		t.FailNow()
	}

	if queue.Close() != ErrClosed {
		t.FailNow()
	}

	<-ch
}
