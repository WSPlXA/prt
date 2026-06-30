package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"prt/types"
)

var _ Hasher[*Block] = (*BlockHasher)(nil)

type BlockHasher struct{}

func (bh *BlockHasher) Hash(b *Block) types.Hash {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(b.Header); err != nil {
		panic(err)
	}
	h := sha256.Sum256(buf.Bytes())
	return types.Hash(h)
}
