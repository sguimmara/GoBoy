// this package is automatically generated by instruction_generator.py
package cpu

import . "memory"


// dec  A - A = A-1
func x3D_dec() int {
	original := GetA()
    value := original - 1
    SetA(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  B - B = B-1
func x05_dec() int {
	original := GetB()
    value := original - 1
    SetB(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  C - C = C-1
func x0D_dec() int {
	original := GetC()
    value := original - 1
    SetC(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  D - D = D-1
func x15_dec() int {
	original := GetD()
    value := original - 1
    SetD(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  E - E = E-1
func x1D_dec() int {
	original := GetE()
    value := original - 1
    SetE(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  H - H = H-1
func x25_dec() int {
	original := GetH()
    value := original - 1
    SetH(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  L - L = L-1
func x2D_dec() int {
	original := GetL()
    value := original - 1
    SetL(value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 1
}

// dec  (HL) - (HL)=(HL)-1
func x35_dec() int {
	original := Get(GetHL())
    value := original - 1
    Set(GetHL(), value)

    hc := IsSubHalfCarry(original, uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_8bit)

	return 3
}

// dec  BC - BC = BC-1
func x0B_dec() int {
	original := GetBC()
    value := original - 1
    SetBC(value)

    hc := IsSubHalfCarry(getHighBits(original), uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_16bit)

	return 2
}

// dec  DE - DE = DE-1
func x1B_dec() int {
	original := GetDE()
    value := original - 1
    SetDE(value)

    hc := IsSubHalfCarry(getHighBits(original), uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_16bit)

	return 2
}

// dec  HL - HL = HL-1
func x2B_dec() int {
	original := GetHL()
    value := original - 1
    SetHL(value)

    hc := IsSubHalfCarry(getHighBits(original), uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_16bit)

	return 2
}

// dec  SP - SP = SP-1
func x3B_dec() int {
	original := GetSP()
    value := original - 1
    SetSP(value)

    hc := IsSubHalfCarry(getHighBits(original), uint8(1))

    SetFlags(int(value), F_SET_IF, F_SET_1, hc, F_IGNORE, F_16bit)

	return 2
}

