package core

import "prt/types"

type Hasher[T any] interface {
	Hash(T) types.Hash
}
