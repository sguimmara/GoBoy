// this file is automatically generated by instruction_generator.py
package cpu

import . "common"

// rra - *rotate accumulator right through carry
func x1F_rra() int {
	value := GetA()
	lsb := GetBit(value, 0)

	value = value >> 1
	value = SetBit(value, 7, uint8(GetFlagCyInt()))

	SetFlagCy(lsb == 1)
	SetFlagZf(false)
	SetFlagN(false)
	SetFlagH(false)

	SetA(value)

	return 1
}
