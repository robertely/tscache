package main

import (
	"fmt"
	"time"
)

type point struct {
	Value interface{}
	Time  time.Time
	Next  *point
}

func (point *point) Update(value interface{}, timestamp time.Time) {
	// should...I be locking this?
	// bruuuup... DOnt thIK about it...just jst dont think about it mooooorrrty
	point.Value = value
	point.Time = timestamp
}

// Collection is a Circular Singly Linked List implimentation that contains time series data.
type Collection struct {
	head     *point
	tail     *point
	length   uint64
	capacity uint64
	// index map...
	// indexinterval uint64
}

func (collection *Collection) Write(value float64, timestamp time.Time) {
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

func (collection *Collection) Read(start time.Time, end time.Time) {
	//uhhhh
}

func (collection *Collection) printAll() {
	currpoint := collection.tail
	lastpoint := collection.head
	for currpoint != lastpoint {
		fmt.Println(*currpoint)
		currpoint = currpoint.Next
	}
	fmt.Println(*currpoint)
}

func main() {
	x := Collection{head: nil, tail: nil, length: 0, capacity: 1000}

	for i := 0; i < 1000000; i++ {
		x.Write(float64(i), time.Time{})
	}

	x.printAll()
	fmt.Println(x)

}
