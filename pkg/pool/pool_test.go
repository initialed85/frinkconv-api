package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	item1 := "item 1"
	item2 := "item 2"
	item3 := "item 3"

	p := New(2)

	item, err := p.GetTimeout(time.Millisecond)
	assert.NotNil(t, err)
	assert.Nil(t, item)

	err = p.PutTimeout(item1, time.Millisecond)
	assert.Nil(t, err)

	item, err = p.GetTimeout(time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, item1, item)

	err = p.PutTimeout(item1, time.Millisecond)
	assert.Nil(t, err)

	err = p.PutTimeout(item2, time.Millisecond)
	assert.Nil(t, err)

	err = p.PutTimeout(item3, time.Millisecond)
	assert.NotNil(t, err)

	item, err = p.GetTimeout(time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, item1, item)

	item, err = p.GetTimeout(time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, item2, item)
}
