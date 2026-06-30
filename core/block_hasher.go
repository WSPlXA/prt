package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"prt/types"
)

var _ Hasher[any] = (*BlockHasher)(nil)

type BlockHasher struct{}

func (bh *BlockHasher) Hash(T any) types.Hash {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(T.(*Block).Header); err != nil {
		panic(err)
	}
	h := sha256.Sum256(buf.Bytes())
	return types.Hash(h)
}
