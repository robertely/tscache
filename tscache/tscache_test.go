package tscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateNode(t *testing.T) {
	n := node{0, time.Time{}, nil}
	assert.Equal(t, n.Value, 0)
	n.update(1, time.Time{})
	assert.Equal(t, n.Value, 1)
}

func TestCreateCollection(t *testing.T) {
	c := NewCollection(10)
	assert.Equal(t, c.Capacity(), 10)
	assert.Equal(t, c.Length(), 0)
}

func TestCollectionWrite(t *testing.T) {
	c := NewCollection(3)
	// Fill the collection
	for i := 0; i < 3; i++ {
		c.Write(i, time.Time{})
		assert.Equal(t, i+1, c.Length(), "Failed to get Length")
	}
	// Overrun the collection
	c.Write(4, time.Time{})
	assert.Equal(t, 3, c.Length(), "Collection Overrun failed")
	assert.Equal(t, 1, c.tail.Value, "Result length incorrect")
	assert.Equal(t, 4, c.head.Value)
}

func TestCollectionSearch(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	start, end, len := c.Search(time.Unix(8, 0), time.Unix(12, 0))
	assert.Equal(t, 5, len, "Result length incorrect")
	assert.Equal(t, 8, start.Value, "Result start value incorrect")
	assert.Equal(t, 12, end.Value, "Result end value incorrect")
}

func TestCollectionZeroSearch(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	start, end, len := c.Search(time.Time{}, time.Time{})
	assert.Equal(t, 10, len, "Result length incorrect")
	assert.Equal(t, 5, start.Value, "Result start value incorrect")
	assert.Equal(t, 14, end.Value, "Result end value incorrect")
}

func TestCollectionZeroStartSearch(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	start, end, len := c.Search(time.Time{}, time.Unix(8, 0))
	// c.PrintAll()
	assert.Equal(t, 4, len, "Result length incorrect")
	assert.Equal(t, 5, start.Value, "Result start value incorrect")
	assert.Equal(t, 8, end.Value, "Result end value incorrect")
}

func TestCollectionZeroEndSearch(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	start, end, len := c.Search(time.Unix(8, 0), time.Time{})
	assert.Equal(t, 7, len, "Result length incorrect")
	assert.Equal(t, 8, start.Value, "Result start value incorrect")
	assert.Equal(t, 14, end.Value, "Result end value incorrect")
}
