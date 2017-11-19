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

//  i think there may be an obo in here...
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

func (collection *Collection) search(start time.Time, end time.Time) (ResultTail, ResultHead *point) {
	ResultTail = nil
	ResultHead = nil

	// validate here.
	if start.After(end) {
		return
	}
	if start.After(collection.head.Time) {
		return
	}

	// Find start point
	if start.Before(collection.tail.Time) {
		ResultTail = collection.tail
	} else {
		// find the real start
		curr := collection.tail
		for start.After(curr.Time) || start.Equal(curr.Time) {
			ResultTail = curr
			curr = curr.Next
		}
	}
	// Find end point
	if end.After(collection.head.Time) {
		ResultHead = collection.head
	} // find the real end
	fmt.println("DONE")
	return
}

func (collection *Collection) Read(start time.Time, end time.Time) {
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
	x := Collection{head: nil, tail: nil, length: 0, capacity: 1000}

	for i := int64(10); i <= 20; i++ {
		x.Write(float64(i), time.Unix(i, 0))
	}

	x.printAll()
	fmt.Println("XXXXXXXXXX")
	fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(14, 0)))
	fmt.Println(x.search(time.Unix(12, 0), time.Unix(14, 0)))

}
