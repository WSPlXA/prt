package core

import (
	"fmt"
	"prt/types"
	"testing"
	"time"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Height:    height,
		Timestamp: time.Now().UnixNano(),
	}
	tx := Transaction{
		Data: []byte("foo"),
	}
	return NewBlock(header, []Transaction{tx})
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(0)
	var hasher Hasher[*Block] = &BlockHasher{}
	fmt.Println(b.Hash(hasher))
}
