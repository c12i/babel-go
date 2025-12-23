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
	Text   string `arg:"" help:"Text to search for"`
	Offset int    `help:"Starting position" default:"0"`
	Limit  int    `help:"Number of results" default:"10"`
}

func (s *SearchCmd) Run(ctx *Context) error {
	lib := ctx.Library

	totalCount := lib.GetOccurrenceCount(s.Text)
	locations, err := lib.SearchPaginated(s.Text, s.Offset, s.Limit)
	if err != nil {
		return err
	}

	fmt.Printf("Text '%s' appears in %d locations. Showing %d results starting from %d:\n\n",
		s.Text, totalCount, len(locations), s.Offset+1)

	for i, location := range locations {
		fmt.Printf("  %d. %s\n", s.Offset+i+1, location.String())
	}

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
