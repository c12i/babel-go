package library

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
)

const (
	PagesPerBook    = 410
	BooksPerShelf   = 32
	ShelvesPerWall  = 5
	WallsPerHexagon = 4
	LinesPerPage    = 40
	CharsPerLine    = 80
	CharsPerPage    = LinesPerPage * CharsPerLine // 3200
)

type Location struct {
	Hexagon string
	Wall    int
	Shelf   int
	Book    int
	Page    int
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

	if len(text) > CharsPerPage {
		return nil, errors.New("text exceeds 3200 character limit")
	}

	result := big.NewInt(0)

	// Pad original string with spaces
	// XXX: Final implementation will need the text padded with random characters
	var builder strings.Builder
	builder.WriteString(text)
	for range CharsPerPage - len(text) {
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
	quotient := new(big.Int)
	remainder := new(big.Int)

	runeSlice := []rune{}

	for temp.Sign() > 0 {
		quotient, remainder = quotient.DivMod(temp, l.base, remainder)
		charByte := l.charset[remainder.Int64()]
		runeSlice = append(runeSlice, rune(charByte))
		temp.Set(quotient)
	}

	slices.Reverse(runeSlice)

	return string(runeSlice)
}

// Given a big int, get the Location
func LocationFromBigInt(n *big.Int) *Location {
	temp := new(big.Int).Abs(n)
	quotient := new(big.Int)

	// Get page
	page := new(big.Int)
	quotient, page = quotient.DivMod(temp, big.NewInt(PagesPerBook), page)
	temp.Set(quotient)

	// Get book
	book := new(big.Int)
	quotient, book = quotient.DivMod(temp, big.NewInt(BooksPerShelf), book)
	temp.Set(quotient)

	// Get shelf
	shelf := new(big.Int)
	quotient, shelf = quotient.DivMod(temp, big.NewInt(ShelvesPerWall), shelf)
	temp.Set(quotient)

	// Get wall
	wall := new(big.Int)
	quotient, wall = quotient.DivMod(temp, big.NewInt(WallsPerHexagon), wall)
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
	result.Mul(result, big.NewInt(WallsPerHexagon))
	result.Add(result, big.NewInt(int64(l.Wall)))

	// Add shelf
	result.Mul(result, big.NewInt(ShelvesPerWall))
	result.Add(result, big.NewInt(int64(l.Shelf)))

	// Add book
	result.Mul(result, big.NewInt(BooksPerShelf))
	result.Add(result, big.NewInt(int64(l.Book)))

	// Add page
	result.Mul(result, big.NewInt(PagesPerBook))
	result.Add(result, big.NewInt(int64(l.Page)-1))

	return result, nil
}
