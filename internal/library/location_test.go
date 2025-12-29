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
	expectedLocation := Location{
		Hexagon: big.NewInt(0).Text(36),
		Wall:    2,
		Shelf:   1,
		Book:    12,
		Page:    30,
	}
	if !location.Equals(expectedLocation) {
		t.Errorf("got %+v, want %+v", location, expectedLocation)
	}
}

func TestLocationFromLongAndShortAddressStrings(t *testing.T) {
	// Too short
	_, err := LocationFromString(
		fmt.Sprintf("%s.2.1", big.NewInt(0).Text(36)),
	) // missing book and page
	if err == nil {
		t.Errorf("got nil, expected err")
	}

	// Too long
	_, err2 := LocationFromString(
		fmt.Sprintf("%s.2.1.12.30.34.39.29", big.NewInt(0).Text(36)),
	) // missing book and page
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
	_, err2 := LocationFromString(fmt.Sprintf("%s.-30.1.12.30", big.NewInt(0).Text(36)))
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

	originalNum, err := library.generateBase29Number("Hello world", 0)
	if err != nil {
		t.Errorf("failed to base29 encode: %v", err)
	}

	location := locationFromBase29Number(originalNum)
	number, err := location.ToBigInt()
	if err != nil {
		t.Errorf("location to big.Int conversion failed: %v", err)
	}

	if originalNum.Cmp(number) != 0 {
		t.Errorf("failed to get back original number")
	}
}

func TestGetLocationWithInvalidHexagonString(t *testing.T) {
	library := NewLibrary()
	number, err := library.generateBase29Number("Hello world", 0)
	if err != nil {
		t.Errorf("failed to base29 encode: %v", err)
	}
	location := locationFromBase29Number(number)
	location.Hexagon = "invalid base32 string"
	_, err2 := location.ToBigInt()
	if err2 == nil {
		t.Errorf("invalid base32 string succeeded when it was expected to fail: %v", err2)
	}
}

/*
	TESTING Next() and Previous() navigation
*/

func TestLocationNext_SimplePage(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	next := location.Next()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    2,
	}

	if !next.Equals(expected) {
		t.Errorf("got %+v, want %+v", next, expected)
	}
}

func TestLocationNext_PageOverflowToBook(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    410, // max page
	}

	next := location.Next()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    1,
		Page:    1,
	}

	if !next.Equals(expected) {
		t.Errorf("got %+v, want %+v", next, expected)
	}
}

func TestLocationNext_BookOverflowToShelf(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    31, // max book (0-31)
		Page:    410,
	}

	next := location.Next()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   1,
		Book:    0,
		Page:    1,
	}

	if !next.Equals(expected) {
		t.Errorf("got %+v, want %+v", next, expected)
	}
}

func TestLocationNext_ShelfOverflowToWall(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   4, // max shelf (0-4)
		Book:    31,
		Page:    410,
	}

	next := location.Next()
	expected := Location{
		Hexagon: "0",
		Wall:    1,
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	if !next.Equals(expected) {
		t.Errorf("got %+v, want %+v", next, expected)
	}
}

func TestLocationNext_WallOverflowToHexagon(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    3, // max wall (0-3)
		Shelf:   4,
		Book:    31,
		Page:    410,
	}

	next := location.Next()
	expected := Location{
		Hexagon: "1",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	if !next.Equals(expected) {
		t.Errorf("got %+v, want %+v", next, expected)
	}
}

func TestLocationNext_HexagonIncrement(t *testing.T) {
	location := Location{
		Hexagon: "abc",
		Wall:    3,
		Shelf:   4,
		Book:    31,
		Page:    410,
	}

	next := location.Next()
	expected := Location{
		Hexagon: "abd",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	if !next.Equals(expected) {
		t.Errorf("got %+v, want %+v", next, expected)
	}
}

func TestLocationPrevious_SimplePage(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    2,
	}

	prev := location.Previous()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	if !prev.Equals(expected) {
		t.Errorf("got %+v, want %+v", prev, expected)
	}
}

func TestLocationPrevious_PageUnderflowToBook(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    1,
		Page:    1, // min page
	}

	prev := location.Previous()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    410,
	}

	if !prev.Equals(expected) {
		t.Errorf("got %+v, want %+v", prev, expected)
	}
}

func TestLocationPrevious_BookUnderflowToShelf(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   1,
		Book:    0, // min book
		Page:    1,
	}

	prev := location.Previous()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   0,
		Book:    31,
		Page:    410,
	}

	if !prev.Equals(expected) {
		t.Errorf("got %+v, want %+v", prev, expected)
	}
}

func TestLocationPrevious_ShelfUnderflowToWall(t *testing.T) {
	location := Location{
		Hexagon: "0",
		Wall:    1,
		Shelf:   0, // min shelf
		Book:    0,
		Page:    1,
	}

	prev := location.Previous()
	expected := Location{
		Hexagon: "0",
		Wall:    0,
		Shelf:   4,
		Book:    31,
		Page:    410,
	}

	if !prev.Equals(expected) {
		t.Errorf("got %+v, want %+v", prev, expected)
	}
}

func TestLocationPrevious_WallUnderflowToHexagon(t *testing.T) {
	location := Location{
		Hexagon: "1",
		Wall:    0, // min wall
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	prev := location.Previous()
	expected := Location{
		Hexagon: "0",
		Wall:    3,
		Shelf:   4,
		Book:    31,
		Page:    410,
	}

	if !prev.Equals(expected) {
		t.Errorf("got %+v, want %+v", prev, expected)
	}
}

func TestLocationPrevious_HexagonDecrement(t *testing.T) {
	location := Location{
		Hexagon: "abd",
		Wall:    0,
		Shelf:   0,
		Book:    0,
		Page:    1,
	}

	prev := location.Previous()
	expected := Location{
		Hexagon: "abc",
		Wall:    3,
		Shelf:   4,
		Book:    31,
		Page:    410,
	}

	if !prev.Equals(expected) {
		t.Errorf("got %+v, want %+v", prev, expected)
	}
}

func TestLocationNextPrevious_RoundTrip(t *testing.T) {
	// Test that Next().Previous() returns to original
	location := Location{
		Hexagon: "3a7f",
		Wall:    2,
		Shelf:   3,
		Book:    15,
		Page:    204,
	}

	roundTrip := location.Next().Previous()
	if !roundTrip.Equals(location) {
		t.Errorf("Next().Previous() roundtrip failed: got %+v, want %+v", roundTrip, location)
	}

	// Test that Previous().Next() returns to original
	roundTrip2 := location.Previous().Next()
	if !roundTrip2.Equals(location) {
		t.Errorf("Previous().Next() roundtrip failed: got %+v, want %+v", roundTrip2, location)
	}
}
