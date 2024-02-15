package blockchain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test the block creation
func TestBlock(t *testing.T) {
	b := NewBlock("test", []byte{})
	assert.NotNil(t, b, "Block creation failed")
	fmt.Println(b.String())
}