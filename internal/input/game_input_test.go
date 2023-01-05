package input_test

import (
	"testing"

	"github.com/assemblaj/ggpo/internal/input"
)

func TestNewGameInput(t *testing.T) {
	_, err := input.NewGameInput(0, []byte{1, 2, 3, 4}, 4)
	if err != nil {
		t.Errorf("Creating a perfectly valid input return nil.")
	}
}

func TestGameInputEqual(t *testing.T) {
	gi, _ := input.NewGameInput(0, []byte{1, 2, 3, 4}, 4)
	gi2, _ := input.NewGameInput(0, []byte{1, 2, 3, 4}, 4)
	isEqual, err := gi.Equal(&gi2, false)
	if err != nil {
		t.Errorf("Checking for equality on inputs created from the exact same data created an error.")
	}

	if !isEqual {
		t.Errorf("Inputs created from the exact same parameters are said to not be equal.")
	}
}

func TestGameInputEqualBitsOnly(t *testing.T) {
	gi, _ := input.NewGameInput(1, []byte{1, 2, 3, 4}, 4)
	gi2, _ := input.NewGameInput(0, []byte{1, 2, 3, 4}, 4)
	isEqual, err := gi.Equal(&gi2, true)
	if err != nil {
		t.Errorf("Checking for equality created an error.")
	}

	if !isEqual {
		t.Errorf("Inputs created with the same bytes but with different frames are said to be unequal even though the equal checks for bits only.")
	}
}

func TestGameInputEqualDifferentSizes(t *testing.T) {
	gi, _ := input.NewGameInput(1, []byte{1, 2, 3, 4}, 4)
	gi2, _ := input.NewGameInput(0, []byte{1, 2, 3, 4, 5, 6}, 6)
	isEqual, err := gi.Equal(&gi2, false)
	if err != nil {
		t.Errorf("Checking for equality on inputs that are a different size returned an error instead of just being false..")
	}
	if isEqual {
		t.Errorf("Inputs with different sizes returned as equal. ")
	}
}

func TestGameInputEqualZeroSize(t *testing.T) {
	gi, _ := input.NewGameInput(1, []byte{1, 2, 3, 4}, 4)
	gi2 := input.GameInput{}
	_, err := gi.Equal(&gi2, false)
	if err == nil {
		t.Errorf("Checking for equality on inputs with sizes <= 0 didn't return an error.")
	}
}

func TestGameInputErase(t *testing.T) {
	gi, _ := input.NewGameInput(1, []byte{1, 2, 3, 4}, 4)
	gi.Erase()
	for i := 0; i < len(gi.Bits); i++ {
		if gi.Bits[0] != 0 {
			t.Errorf("Erased inputs aren't zero'd.")
		}
	}
}

func TestGameInputString(t *testing.T) {
	gi, _ := input.NewGameInput(0, []byte{1, 2, 3, 4}, 4)
	gi2, _ := input.NewGameInput(0, []byte{1, 2, 3, 4}, 4)

	if gi.String() != gi2.String() {
		t.Errorf("Inputs created from the exact same parameters did not created the same strings.")
	}
}

func TestGameZeroSize(t *testing.T) {
	_, err := input.NewGameInput(0, []byte{1, 2, 3, 4}, 0)

	if err == nil {
		t.Errorf("Inputs with sizes <= 0 should create an error.")
	}
}

func TestGameInputIsNull(t *testing.T) {
	input, _ := input.NewGameInput(input.NullFrame, []byte{1, 2, 3, 4}, 4)
	want := true
	got := input.IsNull()
	if want != got {
		t.Errorf("want %t got %t.", want, got)
	}
}
