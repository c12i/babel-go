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
	"strings"
	"sync"
)

const (
	pagesPerBook    = 410
	booksPerShelf   = 32
	shelvesPerWall  = 5
	wallsPerHexagon = 4
	linesPerPage    = 40
	charsPerLine    = 80
	charsPerPage    = linesPerPage * charsPerLine // 3200
	// the rate at which search results in the library continue to reduce as the length of the search text increases
	exponentialDecayRate = 1.10
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
		for variant := range totalCount {
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

	// validate parameters
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
	maxCount := 1_000_000_000 // 100M for single characters
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
