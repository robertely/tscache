package tscache

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
	length   int
	capacity int
	mutex    sync.Mutex
}

// NewCollection returns initialized Collection
func NewCollection(capacity int) *Collection {
	return &Collection{
		capacity: capacity,
	}
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
func (collection *Collection) Length() int {
	return collection.length
}

// Capacity returns the capacity of a Collection
func (collection *Collection) Capacity() int {
	return collection.capacity
}

// TODO: Not happy with this make do better
// Search js jesus christ i dono
func (collection *Collection) Search(start, end time.Time) (ResultTail, ResultHead *node, length int) {
	// Aquire Lock
	collection.mutex.Lock()
	defer collection.mutex.Unlock()

	// Validation
	// Range error
	if start.After(end) && !end.IsZero() {
		return
	}

	// Passed an empty collection?
	if collection.head == nil {
		return
	}

	// Range error
	if start.After(collection.head.Time) {
		return
	}
	// return Complete set
	if start.IsZero() && end.IsZero() {
		ResultHead = collection.head
		ResultTail = collection.tail
		length = collection.Length()
		return
	}

	// Find start node
	ResultTail = collection.tail
	if !start.IsZero() {
		for start.After(ResultTail.Time) {
			ResultTail = ResultTail.Next
		}
	}

	// Find end node
	// breaks with fractional seconds.... # THIS PART IS GROSS
	if end.IsZero() {
		ResultHead = ResultTail
		for ResultHead.Next != collection.tail {
			ResultHead = ResultHead.Next
			length++
		}
	} else {
		ResultHead = ResultTail
		for end.After(ResultHead.Time) && ResultHead.Next != collection.tail {
			ResultHead = ResultHead.Next
			length++
		}
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
func (collection *Collection) Read(start, end *node, length int) []Point {
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
func (collection *Collection) PrintAll() {
	currnode := collection.tail
	fmt.Println(*currnode)
	for currnode != collection.head {
		fmt.Println(*currnode.Next)
		currnode = currnode.Next
	}
}
