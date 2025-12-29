package library

import (
	"fmt"
	"strings"
	"testing"
)

const searchText = "hello world"

/*
TESTING core library of babel API
*/

func TestLibrarySearchStream(t *testing.T) {
	library := NewLibrary()
	results, err := library.SearchStream(searchText)
	if err != nil {
		t.Errorf("search stream failed: %v", err)
	}
	locations := []*Location{}
	limit := 100
	for range limit {
		location := <-results
		locations = append(locations, location)
	}
	if l := len(locations); l != limit {
		t.Errorf("expected %d locations, got %d", limit, l)
	}
	if err := assertSearchedLocations(library, locations); err != nil {
		t.Errorf("locations assertion failed: %v", err)
	}
}

func TestLibrarySearchPagintated(t *testing.T) {
	library := NewLibrary()
	limit, offset := 50, 0
	results, err := library.SearchPaginated(searchText, offset, limit)
	if err != nil {
		t.Errorf("search stream failed: %v", err)
	}
	if l := len(results); l != limit {
		t.Errorf("expected %d results, got %d", limit, l)
	}
	if err := assertSearchedLocations(library, results); err != nil {
		t.Errorf("locations assertion failed: %v", err)
	}
	offset = 50
	results2, err := library.SearchPaginated(searchText, offset, limit)
	if err != nil {
		t.Errorf("search stream failed: %v", err)
	}
	if l := len(results2); l != limit {
		t.Errorf("expected %d results, got %d", limit, l)
	}
	if err := assertSearchedLocations(library, results); err != nil {
		t.Errorf("locations assertion failed: %v", err)
	}
}

func TestLibrarySearchPaginatedWithInvalidLimit(t *testing.T) {
	library := NewLibrary()
	limit, offset := -50, 0
	_, err := library.SearchPaginated(searchText, offset, limit)
	if err == nil {
		t.Errorf("expected err, found nil")
	}
}

func TestLibrarySearchPaginatedWithInvalidOffset(t *testing.T) {
	library := NewLibrary()
	limit, offset := 50, -50
	_, err := library.SearchPaginated(searchText, offset, limit)
	if err == nil {
		t.Errorf("expected err, found nil")
	}
}

func assertSearchedLocations(library *Library, locations []*Location) error {
	for _, location := range locations {
		pageContent, err := library.Browse(location)
		if err != nil {
			return fmt.Errorf("failed to browse location: %v", err)
		}
		if !strings.Contains(pageContent, searchText) {
			return fmt.Errorf("location page does not contain search text: %s", searchText)
		}
	}
	return nil
}

/*
TESTING Base29 encode and decoding
*/
func TestLibraryBase29EncodeAndDecode(t *testing.T) {
	library := NewLibrary()
	num, err := library.generateBase29Number(searchText, 0)
	if err != nil {
		t.Errorf("failed to encode text: %v", err)
	}
	pageContent := library.base29NumberToString(num)

	if !strings.Contains(pageContent, searchText) {
		t.Errorf(
			"input string \"%s\" not in page content",
			searchText,
		)
	}
}

func TestLibraryBase29EmptyString(t *testing.T) {
	library := NewLibrary()
	_, error := library.generateBase29Number("", 0)

	if error == nil {
		t.Errorf("empty strings should error")
	}
}

func TestLibraryBase29EncodeLongText(t *testing.T) {
	text := strings.Repeat("hello", 1000)
	library := NewLibrary()
	_, error := library.generateBase29Number(text, 0)

	if error == nil {
		t.Errorf("empty strings should error")
	}
}

func TestLibraryBase29EncodeInvalidChars(t *testing.T) {
	library := NewLibrary()
	invalidChars := "!@#$%^&*()_+-=[]{}|;':\"<>?/~`"

	for _, char := range invalidChars {
		_, err := library.generateBase29Number(fmt.Sprintf("hello%c", char), 0)
		if err == nil {
			t.Errorf("encoded with invalid character: %c", char)
		}
	}
}
