# The Library of Babel

The library of Babel contains every possible 410 page book with 40 lines per page and 80 characters per line, Using 29 characters (space, a-z, comma and period).
Each book has a unique location (hexagon, wall, shelf, volume). The core algorithm converts text to location coordinates using base conversion.

## Context

The [Library of Babel](https://en.wikipedia.org/wiki/The_Library_of_Babel) is a short story by Jorge Luis Borges about an infinite library containing every possible book. Not just books that exist, but every combination of
characters that _could_ exist.

-   Your biography (that hasn't been written yet)
-   This exact documentation
-   Complete nonsense
-   Shakespeare's works

## Inspiration

-   Jonathan Basile's implementation at [libraryofbabel.info](https://libraryofbabel.info)
-   Addressing system was inspired by [@tdjsnelling](https://github.com/tdjsnelling)'s [TypeScript implementation](https://github.com/tdjsnelling/babel)

## Usage

Core API allows a user to:

-   Text Search -> Search for locations where user text exists in the library
-   Browse -> View the page contents of a given location
-   Random -> View a page from a random location in the library

You can interact with the program via [web app](https://babel.c12i.xyz) and a very simple [CLI](./cmd/cli/main.go)

contact: [`hello@collinsmuriuki.xyz`](mailto:hello@collinsmuriuki.xyz)
