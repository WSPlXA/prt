package core

import "prt/crypto"

type Transaction struct {
	Data []byte

	PublicKey crypto.PublicKey
	Signature crypto.Signature
}
