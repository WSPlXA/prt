package core

var _ Validator[*Block] = (*BlockValidator)(nil)

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {

	return &BlockValidator{
		bc: bc,
	}
}

func (bv *BlockValidator) Validate(b *Block) error {

	return nil
}
