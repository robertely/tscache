package tscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestZeroLengthRead(t *testing.T) {
	assert := require.New(t)

	c := NewCollection(0)

	head, tail, length := c.Search(time.Unix(12, 0), time.Unix(20, 0))
	assert.Nil(head)
	assert.Nil(tail)
	assert.Zero(length)
}

func TestOneLengthSearch(t *testing.T) {
	assert := require.New(t)

	c := NewCollection(1)
	c.Write(3.14, time.Now())

	head, tail, length := c.Search(time.Unix(0, 0), time.Now())

	assert.NotNil(head)
	assert.NotNil(tail)
	assert.Equal(int(1), int(length))
}

func TestManySearch(t *testing.T) {
	assert := require.New(t)

	c := NewCollection(5000)

	for i := 0; i < c.Capacity(); i++ {
		c.Write(i, time.Unix(int64(i), 0))
	}

	head, tail, length := c.Search(time.Unix(0, 0), time.Now())

	assert.NotNil(head)
	assert.NotNil(tail)
	assert.Equal(int(5000), int(length))
}
