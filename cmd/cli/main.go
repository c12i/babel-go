package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/c12i/babel-go/internal/library"
)

var CLI struct {
	Search SearchCmd `cmd:"" help:"Search for text in the library of Babel"`
	Browse BrowseCmd `cmd:"" help:"Browse a page of a book in the library given its address"`
}

type Context struct {
	Library *library.Library
}

type SearchCmd struct {
	Text string `arg:"" name:"text" help:"Text to search in the library"`
}

func (s *SearchCmd) Run(ctx *Context) error {
	location, err := ctx.Library.Search(s.Text)
	if err != nil {
		return err
	}
	fmt.Printf("Found at:\n")
	fmt.Printf("  Hexagon: %s\n", location.Hexagon)
	fmt.Printf("  Wall:    %d\n", location.Wall)
	fmt.Printf("  Shelf:   %d\n", location.Shelf)
	fmt.Printf("  Book:    %d\n", location.Book)
	fmt.Printf("  Page:    %d\n", location.Page)
	return nil
}

type BrowseCmd struct {
	Address string `arg:"" name:"address" help:"Period separated string of the address to browse in the library: <hexagon>.<wall>.<shelf>.<book>.<page>"`
}

func (s *BrowseCmd) Run(ctx *Context) error {
	location, err := library.LocationFromString(s.Address)
	if err != nil {
		return err
	}
	pageContent, err := ctx.Library.Browse(location)
	if err != nil {
		return err
	}
	fmt.Printf("%s", pageContent)
	return nil
}

func main() {
	ctx := kong.Parse(
		&CLI,
		kong.Name("babel"),
		kong.Description("Library of Babel CLI - Search and browse the infinite library"),
		kong.UsageOnError(),
	)
	err := ctx.Run(&Context{Library: library.NewLibrary()})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
