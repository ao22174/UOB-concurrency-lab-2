package main

import (
	"fmt"
	"github.com/ChrisGora/semaphore"
	"sync"

	//"github.com/ChrisGora/semaphore"
	"math/rand"
	"time"
)

type buffer struct {
	b                 []int
	size, read, write int
}

func newBuffer(size int) buffer {
	return buffer{
		b:     make([]int, size),
		size:  size,
		read:  0,
		write: 0,
	}
}

func (buffer *buffer) get() int {
	x := buffer.b[buffer.read]
	fmt.Println("Get\t", x, "\t", buffer)
	buffer.read = (buffer.read + 1) % len(buffer.b)
	return x
}

func (buffer *buffer) put(x int) {
	buffer.b[buffer.write] = x
	fmt.Println("Put\t", x, "\t", buffer)
	buffer.write = (buffer.write + 1) % len(buffer.b)
}

func producer(buffer *buffer, start, delta int, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex) {
	x := start
	for {
		spaceAvailable.Wait()
		mutex.Lock()
		buffer.put(x)
		x = x + delta
		workAvailable.Post()
		mutex.Unlock()
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func consumer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex) {
	for {
		workAvailable.Wait()
		mutex.Lock()
		_ = buffer.get()
		spaceAvailable.Post()
		mutex.Unlock()
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	buffer := newBuffer(5)
	mutex := new(sync.Mutex)

	spaceAvailable := semaphore.Init(5, 5)
	workAvailable := semaphore.Init(5, 0)

	go producer(&buffer, 1, 1, spaceAvailable, workAvailable, mutex)
	go producer(&buffer, 1000, -1, spaceAvailable, workAvailable, mutex)

	consumer(&buffer, spaceAvailable, workAvailable, mutex)
}
