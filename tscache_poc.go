package main

import (
	"fmt"
	"sync"
	"time"
)

type point struct {
	Value interface{}
	Time  time.Time
	Next  *point
}

func (point *point) Update(value interface{}, timestamp time.Time) {
	point.Value = value
	point.Time = timestamp
}

// Collection is a Circular Singly Linked List implimentation that contains time series data.
type Collection struct {
	head     *point
	tail     *point
	length   uint64
	capacity uint64
	mutex    sync.Mutex
	// index map...
	// indexinterval uint64
}

//  i think there may be an obo in here...
func (collection *Collection) Write(value interface{}, timestamp time.Time) {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()

	// If passed 0 use `now` for timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	// So there are actually two states this can be in. "Growing and Full."
	// When "Growing" it acts like a linked list.
	// When "Full" It acts more like a ring buffer.
	if collection.length < collection.capacity {
		// Create our new point.
		newpoint := &point{value, timestamp, nil}
		// First entry
		if collection.length == 0 {
			newpoint.Next = newpoint
			collection.head = newpoint
			collection.tail = newpoint
		} else {
			// Growing...
			newpoint.Next = collection.tail
			collection.head.Next = newpoint
			collection.head = newpoint
		}
		collection.length++
	} else {
		// Bump tail ahead one
		collection.tail = collection.tail.Next
		// Update the old tail
		collection.head.Next.Update(value, timestamp)
		// Make the old tail the new head.
		collection.head = collection.head.Next
	}
}

// Length returns length of a Collection
func (collection *Collection) Length() uint64 {
	return collection.length
}

// Capacity returns the capacity of a Collection
func (collection *Collection) Capacity() uint64 {
	return collection.capacity
}

func (collection *Collection) search(start, end time.Time) (ResultTail, ResultHead *point) {
	// Validation
	if start.After(end) {
		return
	}
	if start.After(collection.head.Time) {
		return
	}

	// Find start point
	ResultTail = collection.tail
	for start.After(ResultTail.Time) {
		ResultTail = ResultTail.Next
	}

	// Find end point
	ResultHead = ResultTail
	for end.After(ResultHead.Time) && ResultHead.Next != collection.tail {
		fmt.Println(ResultHead.Time)
		ResultHead = ResultHead.Next
	}

	return
}

func (collection *Collection) Read(start *point, end *point) {
	//uhhhh
}

func (collection *Collection) printAll() {
	currpoint := collection.tail
	lastpoint := collection.head
	for currpoint != lastpoint {
		// fmt.Println(currpoint.Time.UnixNano(), ":", currpoint.Value)
		fmt.Println(*currpoint)
		currpoint = currpoint.Next
	}
}

func main() {
	x := Collection{head: nil, tail: nil, length: 0, capacity: 50000000}

	for i := int64(0); i <= 50000000; i++ {
		x.Write(i, time.Now())
	}

	fmt.Println("XXXXXXXXXX")
	// x.printAll()
	// x.Write(20.800000000, time.Unix(20, 800000000))
	// fmt.Println(x.search(time.Unix(12, 0), time.Unix(14, 0)))
	// fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(14, 0)))
	// fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(100, 0)))
	// fmt.Println(x.search(time.Unix(12, 0), time.Unix(20, 800000000)))
	time.Sleep(time.Second * 30)

}
