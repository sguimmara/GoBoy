// this file is automatically generated by instruction_generator.py
package cpu

import . "common"

// rlca - *rotate accumulator left
func x07_rlca() int {
	value := GetA()
	msb := GetBit(value, 7)

	value = value << 1
	value = SetBit(value, 0, uint8(GetFlagCyInt()))

	SetFlagCy(msb == 1)
	SetFlagZf(false)
	SetFlagN(false)
	SetFlagH(false)

	SetA(value)

	return 1
}

