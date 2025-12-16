package library

import (
	"fmt"
	"math/big"
	"testing"
)

/*
	TESTING string address --> Location
*/

func TestLocationFromValidString(t *testing.T) {
	location, err := LocationFromString(fmt.Sprintf("%s.2.1.12.30", big.NewInt(0).Text(36)))
	if err != nil {
		t.Errorf("failed to create location from: %v", location)
	}
	expectedLocation := Location{Hexagon: big.NewInt(0).Text(36), Wall: 2, Shelf: 1, Book: 12, Page: 30}
	if !location.Equals(expectedLocation) {
		t.Errorf("got %+v, want %+v", location, expectedLocation)
	}
}

func TestLocationFromLongAndShortAddressStrings(t *testing.T) {
	// Too short
	_, err := LocationFromString(fmt.Sprintf("%s.2.1", big.NewInt(0).Text(36))) // missing book and page
	if err == nil {
		t.Errorf("got nil, expected err")
	}

	// Too long
	_, err2 := LocationFromString(fmt.Sprintf("%s.2.1.12.30.34.39.29", big.NewInt(0).Text(36))) // missing book and page
	if err2 == nil {
		t.Errorf("got nil, expected err")
	}
}

func TestLocationFromInvalidHexagonString(t *testing.T) {
	invalidChars := "!@#$%^&*()_+-=[]{}|;':\"<>?,./~` "

	for _, ch := range invalidChars {
		input := fmt.Sprintf("hello%c.2.1.12.30", ch)
		_, err := LocationFromString(input)
		if err == nil {
			t.Errorf("got nil, expected err for input %s", input)
		}
	}
}

func TestLocationFromStringWithInvalidWallValue(t *testing.T) {
	// Too big
	_, err := LocationFromString(fmt.Sprintf("%s.30.1.12.30", big.NewInt(0).Text(36)))
	if err == nil {
		t.Errorf("got nil, expected err")
	}

	// Too small
	_, err2 := LocationFromString(fmt.Sprintf("%s.30.1.12.30", big.NewInt(0).Text(36)))
	if err2 == nil {
		t.Errorf("got nil, expected err")
	}
}

func TestLocationFromStringWithInvalidShelfValue(t *testing.T) {
	// Too big
	_, err := LocationFromString(fmt.Sprintf("%s.2.50.12.30", big.NewInt(0).Text(36)))
	if err == nil {
		t.Errorf("got nil, expected err")
	}

	// Too small
	_, err2 := LocationFromString(fmt.Sprintf("%s.2.-12.12.30", big.NewInt(0).Text(36)))
	if err2 == nil {
		t.Errorf("got nil, expected err")
	}
}

func TestLocationFromStringWithInvalidBookValue(t *testing.T) {
	// Too big
	_, err := LocationFromString(fmt.Sprintf("%s.2.1.64.30", big.NewInt(0).Text(36)))
	if err == nil {
		t.Errorf("got nil, expected err")
	}

	// Too small
	_, err2 := LocationFromString(fmt.Sprintf("%s.2.1.-12.30", big.NewInt(0).Text(36)))
	if err2 == nil {
		t.Errorf("got nil, expected err")
	}
}

func TestLocationFromStringWithInvalidPageValue(t *testing.T) {
	// Too big
	_, err := LocationFromString(fmt.Sprintf("%s.2.1.12.599", big.NewInt(0).Text(36)))
	if err == nil {
		t.Errorf("got nil, expected err")
	}

	// Too small
	_, err2 := LocationFromString(fmt.Sprintf("%s.2.1.12.-39", big.NewInt(0).Text(36)))
	if err2 == nil {
		t.Errorf("got nil, expected err")
	}
}

/*
	TESTING Location <---> big.Int conversion
*/

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
