package library

import (
	"fmt"
	"strings"
	"testing"
)

func TestLibraryBase29EncodeAndDecode(t *testing.T) {
	library := NewLibrary()
	input := "hello world"
	num, err := library.generateBase29BigInt(input, 0)
	if err != nil {
		t.Errorf("failed to encode text: %v", err)
	}
	pageContent := library.getStringFromBase29BigInt(num)

	if !strings.Contains(pageContent, input) {
		t.Errorf(
			"input string \"%s\" not in page content",
			input,
		)
	}
}

func TestLibraryBase29EmptyString(t *testing.T) {
	library := NewLibrary()
	_, error := library.generateBase29BigInt("", 0)

	if error == nil {
		t.Errorf("empty strings should error")
	}
}

func TestLibraryBase29EncodeLongText(t *testing.T) {
	text := strings.Repeat("hello", 1000)
	library := NewLibrary()
	_, error := library.generateBase29BigInt(text, 0)

	if error == nil {
		t.Errorf("empty strings should error")
	}
}

func TestLibraryBase29EncodeInvalidChars(t *testing.T) {
	library := NewLibrary()
	invalidChars := "!@#$%^&*()_+-=[]{}|;':\"<>?/~`"

	for _, char := range invalidChars {
		_, err := library.generateBase29BigInt(fmt.Sprintf("hello%c", char), 0)
		if err == nil {
			t.Errorf("encoded with invalid character: %c", char)
		}
	}
}
