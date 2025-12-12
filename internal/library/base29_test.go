package library

import (
	"fmt"
	"strings"
	"testing"
)

func TestLibraryBase29EncodeAndDecode(t *testing.T) {
	library := NewLibrary()
	testString := "hello world"
	num, err := library.Base29Encode(testString)
	if err != nil {
		t.Errorf("failed to encode text: %v", err)
	}
	decodedString := library.Base29Decode(num)
	decodedString = strings.TrimSpace(decodedString)

	if testString != decodedString {
		t.Errorf("decoded base29 number does not match original string: original %v, decodedString: %v", testString, decodedString)
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
