package main

import "github.com/alecthomas/kong"

var CLI struct {
	Search struct {
		Text string `arg:"" name:"text" help:"Text to search in the library"`
	} `cmd:"" help:"Search for text in the library of Babel"`
	Browse struct {
		Address string `arg:"" name:"address" help:"Period separated string of the address to browse in the library: <hexagon>.<wall>.<shelf>.<book>.<page>"`
	} `cmd:"" help:"Browse a page of a book in the library given its address"`
}

func main() {
	ctx := kong.Parse(&CLI)

	switch ctx.Command() {
	case "search <text>":
	case "browse <address>":
	default:
		panic(ctx.Command())
	}
}
