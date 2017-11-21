package main

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	Value interface{}
	Time  time.Time
	Next  *node
}

func (node *node) update(value interface{}, timestamp time.Time) {
	node.Value = value
	node.Time = timestamp
}

// Collection is a Circular Singly Linked List implimentation that contains time series data.
type Collection struct {
	head     *node
	tail     *node
	length   uint64
	capacity uint64
	mutex    sync.Mutex
}

//  Writes a node into a collection as the newest node(head.)
func (collection *Collection) Write(value interface{}, timestamp time.Time) {
	// Aquire Lock
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
		// Create our new node.
		newnode := &node{value, timestamp, nil}
		// First entry
		if collection.length == 0 {
			newnode.Next = newnode
			collection.head = newnode
			collection.tail = newnode
		} else {
			// Growing...
			newnode.Next = collection.tail
			collection.head.Next = newnode
			collection.head = newnode
		}
		collection.length++
	} else {
		// Bump tail ahead one
		collection.tail = collection.tail.Next
		// Update the old tail
		collection.head.Next.update(value, timestamp)
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

// TODO: Not happy with this make do better
func (collection *Collection) search(start, end time.Time) (ResultTail, ResultHead *node, length uint) {
	// Aquire Lock
	collection.mutex.Lock()
	defer collection.mutex.Unlock()

	// Validation
	if start.After(end) {
		return
	}
	if start.After(collection.head.Time) {
		return
	}

	// Find start node
	ResultTail = collection.tail
	for start.After(ResultTail.Time) {
		ResultTail = ResultTail.Next
	}

	// Find end node
	// breaks with fractional seconds....
	ResultHead = ResultTail
	for end.After(ResultHead.Time) && ResultHead.Next != collection.tail {
		ResultHead = ResultHead.Next
		length++
	}
	length++
	return
}

// Point is the external version of node{}
type Point struct {
	Value interface{}
	Time  time.Time
}

// Can i return [length]Point ?
func (collection *Collection) Read(start, end *node, length uint) []Point {
	// Aquire Lock
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	// I guess this is redundant, but what else would i call "start" ?
	currnode := start
	// Build response
	result := make([]Point, length)
	for i := 0; currnode != end.Next; i++ {
		result[i] = Point{Value: currnode.Value, Time: currnode.Time}
		currnode = currnode.Next
	}
	return result
}

// weird things happen in circles..
func (collection *Collection) printAll() {
	currnode := collection.tail
	fmt.Println(*currnode)
	for currnode != collection.head {
		fmt.Println(*currnode.Next)
		currnode = currnode.Next
	}
}

func main() {
	x := Collection{head: nil, tail: nil, length: 0, capacity: 50000}

	for i := int64(10); i <= 22; i++ {
		x.Write(i, time.Unix(i, 0))
	}
	fmt.Println("XXXXXXXXXX")
	x.printAll()
	// fmt.Println(x)
	// x.Write(20.800000000, time.Unix(20, 800000000))
	// fmt.Println(x.search(time.Unix(12, 0), time.Unix(14, 0)))
	// fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(14, 0)))
	// fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(100, 0)))
	resStart, resEnd, length := x.search(time.Unix(12, 0), time.Unix(20, 0))
	fmt.Println(resStart, resEnd, length)
	fmt.Println("XXXXXXXXXX")
	for _, i := range x.Read(x.search(time.Unix(12, 0), time.Unix(20, 0))) {
		fmt.Println(i)
	}
}
