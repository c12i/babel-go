package main

import (
	"fmt"

	"github.com/c12i/babel-go/internal/library"
)

func main() {
	library := library.NewLibrary()
	fmt.Printf("%v", library)
}
