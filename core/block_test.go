package core

import (
	"prt/types"
	"testing"
	"time"
)

func randomBlock(height uint32) {
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

}
