package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	//address := pubKey.Address()

	msg := []byte("Hello world")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)
	fmt.Println(sig)

	assert.True(t, sig.Verify(pubKey, msg))

}

func TestKeyPair_Sign_Verify(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("Hello world")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(pubKey, msg))

}
