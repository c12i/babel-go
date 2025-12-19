package library

import (
	"fmt"
	"strings"
	"testing"
)

func TestLibraryBase29EncodeAndDecode(t *testing.T) {
	library := NewLibrary()
	input := "hello world"
	num, err := library.Base29Encode(input)
	if err != nil {
		t.Errorf("failed to encode text: %v", err)
	}
	pageContent := library.Base29Decode(num)

	if !strings.Contains(pageContent, input) {
		t.Errorf(
			"input string \"%s\" not in page content",
			input,
		)
	}
}

func TestLibraryBase29EmptyString(t *testing.T) {
	library := NewLibrary()
	_, error := library.Base29Encode("")

	if error == nil {
		t.Errorf("empty strings should error")
	}
}

func TestLibraryBase29EncodeLongText(t *testing.T) {
	text := strings.Repeat("hello", 1000)
	library := NewLibrary()
	_, error := library.Base29Encode(text)

	if error == nil {
		t.Errorf("empty strings should error")
	}
}

func TestLibraryBase29EncodeInvalidChars(t *testing.T) {
	library := NewLibrary()
	invalidChars := "!@#$%^&*()_+-=[]{}|;':\"<>?/~`"

	for _, char := range invalidChars {
		_, err := library.Base29Encode(fmt.Sprintf("hello%c", char))
		if err == nil {
			t.Errorf("encoded with invalid character: %c", char)
		}
	}
}
