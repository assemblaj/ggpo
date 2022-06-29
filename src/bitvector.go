package ggthx

const BitVectorNibbleSize int = 8

type BitVector []byte

func (b BitVector) SetBit(offset *int) {
	b[*offset/8] |= (1 << (*offset % 8))
	*offset++
}

func (b BitVector) ClearBit(offset *int) {
	b[(*offset)/8] = b[*offset/8] & ^(1 << ((*offset) % 8))
	*offset += 1
}

func (b BitVector) WriteNibblet(nibble int, offset *int) {
	Assert(nibble < (1 << BitVectorNibbleSize))
	for i := 0; i < BitVectorNibbleSize; i++ {
		if nibble&(1<<i) == 1 {
			b.SetBit(offset)
		} else {
			b.ClearBit(offset)
		}
	}
}

func (b BitVector) ReadBit(offset *int) bool {
	retval := (b[(*offset)/8] & (1 << ((*offset) % 8))) != 0
	*offset += 1
	return retval
}

func (b BitVector) ReadNibblet(offset *int) int {
	nibblet := 0
	for i := 0; i < BitVectorNibbleSize; i++ {
		var result int
		if b.ReadBit(offset) {
			result = 1
		} else {
			result = 0
		}
		nibblet |= (result << i)
	}
	return nibblet
}

/*
func BitVector_SetBit(vector []byte, offset *int) {
	vector[*offset/8] |= (1 << (*offset % 8))
	*offset++
}

func BitVector_ClearBit(vector []byte, offset *int) {
	vector[(*offset)/8] = vector[*offset/8] & ^(1 << ((*offset) % 8))
	*offset += 1
}

func BitVector_WriteNibblet(vector []byte, nibble int, offset *int) error {
	if nibble < (1 << BITVECTOR_NIBBLE_SIZE) {
		for i := 0; i < BITVECTOR_NIBBLE_SIZE; i++ {
			if nibble&(1<<i) == 1 {
				BitVector_SetBit(vector, offset)
			} else {
				BitVector_ClearBit(vector, offset)
			}
		}
		return nil
	} else {
		return errors.New("BitVector_WriteNibblet error")
	}
}

func BitVector_ReadBit(vector []byte, offset *int) bool {
	retval := (vector[(*offset)/8] & (1 << ((*offset) % 8))) != 0
	*offset += 1
	return retval
}

func BitVector_ReadNibblet(vector []byte, offset *int) int {
	nibblet := 0
	for i := 0; i < BITVECTOR_NIBBLE_SIZE; i++ {
		var result int
		if BitVector_ReadBit(vector, offset) {
			result = 1
		} else {
			result = 0
		}
		nibblet |= (result << i)
	}
	return nibblet
}
*/
