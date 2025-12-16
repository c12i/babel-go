package library

import "testing"

func TestGetLocationFromBigInt(t *testing.T) {
	library := NewLibrary()

	originalNum, err := library.Base29Encode("Hello world")
	if err != nil {
		t.Errorf("failed to base29 encode: %v", err)
	}

	location := LocationFromBigInt(originalNum)
	number, err := location.BigIntFromLocation()
	if err != nil {
		t.Errorf("location to big.Int conversion failed: %v", err)
	}

	if originalNum.Cmp(number) != 0 {
		t.Errorf("failed to get back original number")
	}
}

func TestGetLocationWithInvalidHexagonString(t *testing.T) {
	library := NewLibrary()
	number, err := library.Base29Encode("Hello world")
	if err != nil {
		t.Errorf("failed to base29 encode: %v", err)
	}
	location := LocationFromBigInt(number)
	location.Hexagon = "invalid base32 string"
	_, err2 := location.BigIntFromLocation()
	if err2 == nil {
		t.Errorf("invalid base32 string succeeded when it was expected to fail: %v", err2)
	}
}
