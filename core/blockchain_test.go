package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockchian(t *testing.T) {
	bc := NewBlockchain()
	assert.NotNil(t, bc.validator)

}
