package library

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
)

const (
	pagesPerBook    = 410
	booksPerShelf   = 32
	shelvesPerWall  = 5
	wallsPerHexagon = 4
	linesPerPage    = 40
	charsPerLine    = 80
	charsPerPage    = linesPerPage * charsPerLine // 3200
)

type Location struct {
	Hexagon string
	Wall    int
	Shelf   int
	Book    int
	Page    int
}

// Get Location from a period separated string: "<hexagon>.<wall>.<shelf>.<book>.<page>"
func LocationFromString(address string) (*Location, error) {
	parts := strings.Split(address, ".")
	if partsLen := len(parts); partsLen != 5 {
		return nil, fmt.Errorf("address is not of valid length, expected %d, got %d", 5, partsLen)
	}

	// Validate hexagon
	hexagon := parts[0]
	if _, ok := new(big.Int).SetString(hexagon, 36); !ok {
		return nil, fmt.Errorf("invalid hexagon: must be valid base-36 string")
	}

	// Parse and validate numeric parts
	wall, err := parseAndValidate(parts[1], "wall", 0, wallsPerHexagon-1)
	if err != nil {
		return nil, err
	}

	shelf, err := parseAndValidate(parts[2], "shelf", 0, shelvesPerWall-1)
	if err != nil {
		return nil, err
	}

	book, err := parseAndValidate(parts[3], "book", 0, booksPerShelf-1)
	if err != nil {
		return nil, err
	}

	page, err := parseAndValidate(parts[4], "page", 1, pagesPerBook)
	if err != nil {
		return nil, err
	}

	return &Location{
		Hexagon: hexagon,
		Wall:    wall,
		Shelf:   shelf,
		Book:    book,
		Page:    page,
	}, nil
}

func parseAndValidate(s string, name string, min, max int) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %v", name, err)
	}
	if val < min || val > max {
		return 0, fmt.Errorf("%s must be between %d and %d, got %d", name, min, max, val)
	}
	return val, nil
}

func (l Location) Equals(other Location) bool {
	return l.Hexagon == other.Hexagon &&
		l.Wall == other.Wall &&
		l.Shelf == other.Shelf &&
		l.Book == other.Book &&
		l.Page == other.Page
}

type Library struct {
	charset     string
	charToIndex map[rune]int
	base        *big.Int
}

// Build the Library
func NewLibrary() *Library {
	charset := " abcdefghijklmnopqrstuvwxyz,."
	charToIndex := map[rune]int{}

	for i, char := range charset {
		charToIndex[char] = i
	}

	return &Library{
		charset:     charset,
		charToIndex: charToIndex,
		base:        big.NewInt(29),
	}
}

// Converts a given text into a base29 number
func (l Library) Base29Encode(text string) (*big.Int, error) {
	if text == "" {
		return nil, errors.New("text should not be empty")
	}

	if len(text) > charsPerPage {
		return nil, errors.New("text exceeds 3200 character limit")
	}

	result := big.NewInt(0)

	// Pad original string with spaces
	// XXX: Final implementation will need the text padded with random characters
	var builder strings.Builder
	builder.WriteString(text)
	for range charsPerPage - len(text) {
		builder.WriteRune(' ')
	}

	for _, char := range strings.ToLower(builder.String()) {
		result.Mul(result, l.base)
		index, exists := l.charToIndex[char]

		if !exists {
			return nil, fmt.Errorf("text contains invalid characters, supported charset: %v", l.charset)
		}

		result.Add(result, big.NewInt(int64(index)))
	}

	return result, nil
}

// Convert base29 number back to a string
func (l Library) Base29Decode(n *big.Int) string {
	temp := new(big.Int).Abs(n)
	quotient, remainder := new(big.Int), new(big.Int)
	runes := []rune{}

	for temp.Sign() > 0 {
		quotient, remainder = quotient.DivMod(temp, l.base, remainder)
		charByte := l.charset[remainder.Int64()]
		runes = append(runes, rune(charByte))
		temp.Set(quotient)
	}
	slices.Reverse(runes)

	return string(runes)
}

// Given a big int, get the Location
func LocationFromBigInt(n *big.Int) *Location {
	temp := new(big.Int).Abs(n)
	quotient := new(big.Int)

	// Get page
	page := new(big.Int)
	quotient, page = quotient.DivMod(temp, big.NewInt(pagesPerBook), page)
	temp.Set(quotient)

	// Get book
	book := new(big.Int)
	quotient, book = quotient.DivMod(temp, big.NewInt(booksPerShelf), book)
	temp.Set(quotient)

	// Get shelf
	shelf := new(big.Int)
	quotient, shelf = quotient.DivMod(temp, big.NewInt(shelvesPerWall), shelf)
	temp.Set(quotient)

	// Get wall
	wall := new(big.Int)
	quotient, wall = quotient.DivMod(temp, big.NewInt(wallsPerHexagon), wall)
	temp.Set(quotient)

	return &Location{
		// Whatever is left from the quotient is the hexagon identifier
		Hexagon: quotient.Text(36),
		Wall:    int(wall.Int64()),
		Shelf:   int(shelf.Int64()),
		Book:    int(book.Int64()),
		Page:    int(page.Int64()) + 1,
	}
}

// Get Location from a big.Int
func (l Location) BigIntFromLocation() (*big.Int, error) {
	hexagon, ok := new(big.Int).SetString(l.Hexagon, 36)
	if !ok {
		return nil, errors.New("invalid hexagon string format")
	}

	// Build up from hexagon
	result := new(big.Int).Set(hexagon)

	// Add wall
	result.Mul(result, big.NewInt(wallsPerHexagon))
	result.Add(result, big.NewInt(int64(l.Wall)))

	// Add shelf
	result.Mul(result, big.NewInt(shelvesPerWall))
	result.Add(result, big.NewInt(int64(l.Shelf)))

	// Add book
	result.Mul(result, big.NewInt(booksPerShelf))
	result.Add(result, big.NewInt(int64(l.Book)))

	// Add page
	result.Mul(result, big.NewInt(pagesPerBook))
	result.Add(result, big.NewInt(int64(l.Page)-1))

	return result, nil
}
