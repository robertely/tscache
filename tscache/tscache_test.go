package tscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNodeUpdate(t *testing.T) {
	n := node{0, time.Time{}, nil}
	assert.Equal(t, n.Value, 0)
	n.update(1, time.Time{})
	assert.Equal(t, n.Value, 1)
}

func TestCollectionCreate(t *testing.T) {
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

func TestSearchBasic(t *testing.T) {
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

func TestSearchSingleReturn(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	start, end, len := c.Search(time.Unix(10, 0), time.Unix(10, 0))
	assert.Equal(t, 1, len, "Result length incorrect")
	assert.Equal(t, 10, start.Value, "Result start value incorrect")
	assert.Equal(t, 10, end.Value, "Result end value incorrect")
}

func TestSearchZeroStartZeroEnd(t *testing.T) {
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

func TestSearchZeroStart(t *testing.T) {
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

func TestSearchZeroEnd(t *testing.T) {
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

func TestSearchBadQuery(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	rstart, rend, rlen := c.Search(time.Unix(10, 0), time.Unix(5, 0))
	assert.Equal(t, 0, rlen)
	assert.Nil(t, rend)
	assert.Nil(t, rstart)

	rstart, rend, rlen = c.Search(time.Unix(20, 0), time.Unix(30, 0))
	assert.Equal(t, 0, rlen)
	assert.Nil(t, rend)
	assert.Nil(t, rstart)

}

func TestSearchEmptyCollection(t *testing.T) {
	c := NewCollection(10)
	rstart, rend, rlen := c.Search(time.Unix(5, 0), time.Unix(10, 0))
	assert.Equal(t, 0, rlen)
	assert.Nil(t, rend)
	assert.Nil(t, rstart)

}

func TestReadBasic(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	rstart, rend, rlen := c.Search(time.Unix(5, 0), time.Unix(10, 0))
	reply := c.Read(rstart, rend, rlen)
	assert.Equal(t, 6, len(reply))
	assert.Equal(t, 5, reply[0].Value)
	assert.Equal(t, 10, reply[5].Value)
}

func TestReadSinglePoint(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	rstart, rend, rlen := c.Search(time.Unix(5, 0), time.Unix(5, 0))
	reply := c.Read(rstart, rend, rlen)
	assert.Equal(t, 1, len(reply))
	assert.Equal(t, 5, reply[0].Value)
}

func TestReadNilSearch(t *testing.T) {
	c := NewCollection(10)
	// Fill the collection
	for i := 0; i < 15; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
	rstart, rend, rlen := c.Search(time.Unix(5, 0), time.Unix(4, 0))
	reply := c.Read(rstart, rend, rlen)
	assert.Equal(t, 0, len(reply))
}

func BenchmarkCollectionWriteBigCache(b *testing.B) {
	c := NewCollection(b.N)
	b.ResetTimer()
	// Fill the collection
	for i := 0; i < b.N; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
}

func BenchmarkCollectionWriteSmallCache(b *testing.B) {
	c := NewCollection(1000)
	b.ResetTimer()
	// Fill the collection
	for i := 0; i < b.N; i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}
}
