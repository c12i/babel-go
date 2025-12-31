package library

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
)

const (
	maxHexagonCharSize = 3004
	base36Chars        = "0123456789abcdefghijklmnopqrstuvwxyz"
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

	// validate hexagon
	hexagon := parts[0]
	if _, ok := new(big.Int).SetString(hexagon, 36); !ok {
		return nil, fmt.Errorf("invalid hexagon: must be valid base-36 string")
	}

	// parse and validate numeric parts
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

// Determine a Location's given its big Int representation
func locationFromBase29Number(n *big.Int) *Location {
	temp, quotient := new(big.Int).Abs(n), new(big.Int)

	// get page
	page := new(big.Int)
	quotient, page = quotient.DivMod(temp, big.NewInt(pagesPerBook), page)
	temp.Set(quotient)

	// get book
	book := new(big.Int)
	quotient, book = quotient.DivMod(temp, big.NewInt(booksPerShelf), book)
	temp.Set(quotient)

	// get shelf
	shelf := new(big.Int)
	quotient, shelf = quotient.DivMod(temp, big.NewInt(shelvesPerWall), shelf)
	temp.Set(quotient)

	// get wall
	wall := new(big.Int)
	quotient, wall = quotient.DivMod(temp, big.NewInt(wallsPerHexagon), wall)
	temp.Set(quotient)

	return &Location{
		// whatever is left from the quotient is the hexagon identifier
		Hexagon: quotient.Text(36),
		Wall:    int(wall.Int64()),
		Shelf:   int(shelf.Int64()),
		Book:    int(book.Int64()),
		Page:    int(page.Int64()) + 1,
	}
}

func parseAndValidate(s string, name string, min, max int) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %v", name, err)
	}
	if num < min || num > max {
		return 0, fmt.Errorf("%s must be between %d and %d, got %d", name, min, max, num)
	}
	return num, nil
}

// Get Location from a big.Int
func (l Location) ToBigInt() (*big.Int, error) {
	hexagon, ok := new(big.Int).SetString(l.Hexagon, 36)
	if !ok {
		return nil, errors.New("invalid hexagon string format")
	}

	// build up from hexagon
	result := new(big.Int).Set(hexagon)

	// add wall
	result.Mul(result, big.NewInt(wallsPerHexagon))
	result.Add(result, big.NewInt(int64(l.Wall)))

	// add shelf
	result.Mul(result, big.NewInt(shelvesPerWall))
	result.Add(result, big.NewInt(int64(l.Shelf)))

	// add book
	result.Mul(result, big.NewInt(booksPerShelf))
	result.Add(result, big.NewInt(int64(l.Book)))

	// add page
	result.Mul(result, big.NewInt(pagesPerBook))
	result.Add(result, big.NewInt(int64(l.Page)-1))

	return result, nil
}

func (l Location) Equals(other Location) bool {
	return l.Hexagon == other.Hexagon &&
		l.Wall == other.Wall &&
		l.Shelf == other.Shelf &&
		l.Book == other.Book &&
		l.Page == other.Page
}

func (l Location) String() string {
	return fmt.Sprintf("%s.%d.%d.%d.%d", l.Hexagon, l.Wall, l.Shelf, l.Book, l.Page)
}

// Next returns the next page location
func (l Location) Next() *Location {
	next := Location{
		Hexagon: l.Hexagon,
		Wall:    l.Wall,
		Shelf:   l.Shelf,
		Book:    l.Book,
		Page:    l.Page,
	}

	// increment page
	if next.Page < pagesPerBook {
		next.Page++
		return &next
	}

	// page is at max, increment book
	next.Page = 1
	if next.Book < booksPerShelf-1 {
		next.Book++
		return &next
	}

	// book is at max, increment shelf
	next.Book = 0
	if next.Shelf < shelvesPerWall-1 {
		next.Shelf++
		return &next
	}

	// shelf is at max, increment wall
	next.Shelf = 0
	if next.Wall < wallsPerHexagon-1 {
		next.Wall++
		return &next
	}

	// wall is at max, increment hexagon
	next.Wall = 0
	hexInt, ok := new(big.Int).SetString(next.Hexagon, 36)
	if !ok {
		return &next
	}
	hexInt.Add(hexInt, big.NewInt(1))
	next.Hexagon = hexInt.Text(36)

	return &next
}

// Previous returns the previous page location
func (l Location) Previous() *Location {
	prev := Location{
		Hexagon: l.Hexagon,
		Wall:    l.Wall,
		Shelf:   l.Shelf,
		Book:    l.Book,
		Page:    l.Page,
	}

	// decrement page
	if prev.Page > 1 {
		prev.Page--
		return &prev
	}

	// page is at min, decrement book
	prev.Page = pagesPerBook
	if prev.Book > 0 {
		prev.Book--
		return &prev
	}

	// book is at min, decrement shelf
	prev.Book = booksPerShelf - 1
	if prev.Shelf > 0 {
		prev.Shelf--
		return &prev
	}

	// shelf is at min, decrement wall
	prev.Shelf = shelvesPerWall - 1
	if prev.Wall > 0 {
		prev.Wall--
		return &prev
	}

	// wall is at min, decrement hexagon
	prev.Wall = wallsPerHexagon - 1
	hexInt, ok := new(big.Int).SetString(prev.Hexagon, 36)
	if !ok || hexInt.Sign() <= 0 {
		return &prev
	}
	hexInt.Sub(hexInt, big.NewInt(1))
	prev.Hexagon = hexInt.Text(36)

	return &prev
}

// Random generates a random location in the library
func RandomLocation() *Location {
	// Generate random base36 string of random length
	hexagonLen := rand.Intn(maxHexagonCharSize) + 1 //nolint: gosec
	hexagon := make([]byte, hexagonLen)
	for i := range hexagon {
		hexagon[i] = base36Chars[rand.Intn(36)] //nolint:gosec
	}

	return &Location{
		Hexagon: string(hexagon),
		Wall:    rand.Intn(wallsPerHexagon),  //nolint: gosec
		Shelf:   rand.Intn(shelvesPerWall),   //nolint: gosec
		Book:    rand.Intn(booksPerShelf),    //nolint: gosec
		Page:    rand.Intn(pagesPerBook) + 1, // nolint: gosec
	}
}
