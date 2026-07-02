package core

type Blockchain struct {
	headers   []*Header
	store     Storage
	validator *BlockValidator
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		headers: []*Header{},
	}
	bc.validator = NewBlockValidator(bc)
	return bc
}

func (bc Blockchain) height() uint32 {
	return uint32(len(bc.headers) - 1)
}
func (bc *Blockchain) AddBlock(b *Block) error {
	return nil
}
func (bc *Blockchain) SetBlockchainValidator(bv *BlockValidator) error {
	return nil
}
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height < bc.Height()
}

func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}
