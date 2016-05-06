package cpu

import "testing"
import . "memory"

// swap C
func TestB_31_swap(t* testing.T) {
	SetC(0xA9)
	xCB_31_swap()
	testRegister(t, REG_C, 0x9A)
}

// swap [HL]
func TestB_36_swap(t* testing.T) {
	SetHL(0x00B9)
	Set(0x00B9, 0x9C)
	xCB_36_swap()

	if hl := Get(GetHL()); hl != 0xC9 {
		t.Errorf("TestB_31_swap() failed: expected [HL] = 0xC9, got 0x%02X", hl)
	}
}