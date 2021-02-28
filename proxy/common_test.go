package proxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseStringList(t *testing.T) {
	origValues := []string{"item1", "item2", "item3", "item4"}
	newValues := append([]string{}, origValues...)

	ReverseStringList(newValues)
	assert.NotEqual(t, origValues, newValues)
	assert.Equal(t, newValues, []string{"item4", "item3", "item2", "item1"})
}

