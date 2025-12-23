package library

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const (
	pagesPerBook         = 410
	booksPerShelf        = 32
	shelvesPerWall       = 5
	wallsPerHexagon      = 4
	linesPerPage         = 40
	charsPerLine         = 80
	charsPerPage         = linesPerPage * charsPerLine // 3200
	exponentialDecayRate = 1.05
)

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

// Deprecated: Search is deprecated. Use SearchStream or SearchPaginated instead.
func (l Library) Search(text string) (*Location, error) {
	bigInt, err := l.generateBase29Number(text, 0)
	if err != nil {
		return nil, err
	}
	location := locationFromBase29Number(bigInt)
	return location, nil
}

func (l Library) SearchStream(text string) (<-chan *Location, error) {
	totalCount := l.GetOccurrenceCount(text)
	// location and job worker channel
	locationChan, workerChan := make(chan *Location, 100), make(chan int, 100)
	// start fixed number of workers
	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Go(func() {
			// each worker processes multiple variants
			for variant := range workerChan {
				bigInt, err := l.generateBase29Number(text, variant)
				if err != nil {
					continue
				}
				location := locationFromBase29Number(bigInt)
				locationChan <- location
			}
		})
	}

	// send jobs and close channels
	go func() {
		defer close(workerChan)
		for variant := 0; variant < totalCount; variant++ {
			workerChan <- variant
		}
	}()

	// close results when workers finish
	go func() {
		defer close(locationChan)
		wg.Wait()
	}()

	return locationChan, nil
}

func (l Library) SearchPaginated(text string, offset, limit int) ([]*Location, error) {
	totalCount := l.GetOccurrenceCount(text)

	// Validate parameters
	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}
	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	if offset >= totalCount {
		return []*Location{}, nil
	}

	// calculate actual results to return
	endIndex := min(offset+limit, totalCount)
	actualLimit := endIndex - offset

	locations := make([]*Location, 0, actualLimit)

	// generate locations from offset to endIndex
	for variant := offset; variant < endIndex; variant++ {
		bigInt, err := l.generateBase29Number(text, variant)
		if err != nil {
			return nil, fmt.Errorf("error generating location for variant %d: %w", variant, err)
		}

		location := locationFromBase29Number(bigInt)
		locations = append(locations, location)
	}

	return locations, nil
}

// Deterministically determine the occurrence rate of a given text in the library
// using exponential decay
func (l Library) GetOccurrenceCount(text string) int {
	textLen := len(text)

	// decay initial max count exponentially by length
	maxCount := 100_000_000 // 100M for single characters
	baseCount := maxCount / int(math.Pow(exponentialDecayRate, float64(textLen-1)))

	baseCount = max(1, baseCount)

	hash := sha256.Sum256([]byte(strings.ToLower(text)))
	seed := int64(binary.BigEndian.Uint64(hash[:8])) //nolint:gosec
	rng := rand.New(rand.NewSource(seed))            //nolint:gosec

	variation := max(1, baseCount/4) // ±25% variation, minimum 1
	adjustment := rng.Intn(2*variation) - variation

	return max(1, baseCount+adjustment)
}

func (l Library) Browse(location *Location) (string, error) {
	bigInt, err := location.ToBigInt()
	if err != nil {
		return "", err
	}
	pageContent := l.base29NumberToString(bigInt)
	return pageContent, nil
}

// Converts a given text into a base29 number.
func (l Library) generateBase29Number(text string, variant int) (*big.Int, error) {
	if text == "" {
		return nil, errors.New("text should not be empty")
	}
	if len(text) > charsPerPage {
		return nil, errors.New("text exceeds 3200 character limit")
	}

	result := big.NewInt(0)
	pageChars := l.seedPageChars(text, variant)

	for _, char := range pageChars {
		result.Mul(result, l.base)
		index, exists := l.charToIndex[char]

		if !exists {
			return nil, fmt.Errorf(
				"text contains invalid characters, supported charset: %v", l.charset,
			)
		}

		result.Add(result, big.NewInt(int64(index)))
	}

	return result, nil
}

// A deterministic seed based on the hash of the input text is used to generate the position
// The text will appear in the page, the same seed is used to populate the page contents
func (l Library) seedPageChars(text string, variant int) string {
	input := fmt.Sprintf("%s\x00%d", strings.ToLower(text), variant)
	textHash := sha256.Sum256([]byte(input))
	textSeed := int64(binary.BigEndian.Uint64(textHash[:8])) //nolint:gosec // overflow acceptable
	rng := rand.New(rand.NewSource(textSeed))                //nolint:gosec // crypto not needed

	// Generate position from seeded rng
	maxPosition := charsPerPage - len(text)
	position := rng.Intn(maxPosition + 1)

	pageChars := make([]byte, charsPerPage)
	for i := range charsPerPage {
		pageChars[i] = l.charset[rng.Intn(len(l.charset))]
	}

	// insert text at determined position
	copy(pageChars[position:], []byte(strings.ToLower(text)))

	return string(pageChars)
}

// Convert base29 number back to a string
func (l Library) base29NumberToString(n *big.Int) string {
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

// Determine a Location's given its big Int representation
func locationFromBase29Number(n *big.Int) *Location {
	temp, quotient := new(big.Int).Abs(n), new(big.Int)

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
