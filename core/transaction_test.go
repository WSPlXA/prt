package core

import (
	"fmt"
	"prt/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrasation_Sign(t *testing.T) {
	Data := []byte("hello world")
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: Data,
	}
	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
	fmt.Printf("Signed data is %#v", tx.Signature)
}
