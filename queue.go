package main

import "fmt"

// We use Ring buffer here for queue

const QUEUE_SIZE = 255 * 10000

type Queue struct {
	// tracking the next element to pop
	head uint32
	// tracking the next element to push
	tail        uint32
	underArr    []byte
	elementSize []byte
}

func (q *Queue) init() {
	q.head = 0
	q.tail = 0
	q.underArr = make([]byte, QUEUE_SIZE)
	q.elementSize = make([]byte, QUEUE_SIZE/255)
}

// assume all data <= 255 bytes
// each slot in queue is 255 even if data is less than 255 bytes
// the slot will be marked by one integer (like the idx of the element)
// => we know which one to pop later
// get the current tail
func (q *Queue) push(data []byte) {
	copy(q.underArr[q.tail:q.tail+uint32(len(data))], data)
	q.elementSize[q.tail] = byte(len(data))
	// increase tail pointer
	q.tail += 255
	// wrap the tail pointer if exceed the ring buffer size
	q.tail %= QUEUE_SIZE
}

func (q *Queue) pop() []byte {
	msg := q.underArr[q.head : q.head+uint32(q.elementSize[q.head])]
	q.head += 255
	// wrap the head pointer if excedd the ring buffer size
	q.head %= QUEUE_SIZE
	return msg
}

func (q *Queue) debug() {
	fmt.Println("//////////////////////////////////////// Debug queue")
	cur := q.head
	for {
		fmt.Printf("Msg in queue %s\n", q.underArr[cur:cur+uint32(q.elementSize[cur])])
		cur += 255
		cur %= QUEUE_SIZE
		if cur == q.tail {
			break
		}
	}
}
