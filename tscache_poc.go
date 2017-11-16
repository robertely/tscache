package main

import (
	"fmt"
	"time"
	//"sync"
)

type Point struct {
	Value float64
	Time  time.Time
	Next  *Point
}

func (point *Point) Update(value float64, timestamp time.Time) {
	// should...I be locking this?
	point.Value = value
	point.Time = timestamp
}

// Circular singly linked list
type Collection struct {
	Head   *Point
	Tail   *Point
	Length uint64
	Limit  uint64
	// mux   sync.Mutex
	// BucketInterval uint64
}

func (collection *Collection) Write(value float64, timestamp time.Time) {
	// If passed 0 use `now` for timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	// So there are actually two states this can be in. "Growing and Full."
	// When "Growing" it acts like a linked list.
	// When "Full" It acts more like a ring buffer.
	if collection.Length < collection.Limit {
		// Create our new point.
		newpoint := &Point{value, timestamp, nil}
		// First entry
		if collection.Length == 0 {
			newpoint.Next = newpoint
			collection.Head = newpoint
			collection.Tail = newpoint
		} else {
			// Growing...
			newpoint.Next = collection.Tail
			collection.Head.Next = newpoint
			collection.Head = newpoint
		}
		collection.Length++
	} else {
		// Bump tail ahead one
		collection.Tail = collection.Tail.Next
		// Update the old tail
		collection.Head.Next.Update(value, timestamp)
		// Make the old tail the new head.
		collection.Head = collection.Head.Next
	}
}

func (collection *Collection) Read(timestamp time.Time) {
	//start at....tail
	// scan.....
	// return &
}

func (collection *Collection) PrintAll() {
	currpoint := collection.Tail
	lastpoint := collection.Head
	for currpoint != lastpoint {
		fmt.Println(*currpoint)
		currpoint = currpoint.Next
	}
	fmt.Println(*currpoint)
}

func main() {
	// nil head, nil tail, length of 0, max length of 15
	x := Collection{nil, nil, 0, 15}
	for i := 0; i < 22; i++ {
		x.Write(float64(i), time.Unix(0, 0))
	}

	x.PrintAll()
	fmt.Println(x)

}
