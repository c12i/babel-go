# The library of Babel

The library of Babel contains every possible 410 page book with 40 lines per page and 80 characters per line, Using 29 characters (space, a-z, comma and period).
Each book has a unique location (hexagon, wall, shelf, volume). The core algorithm converts text to location coordinates using base conversion.

## Context

The Library of Babel is a short story by Jorge Luis Borges about an infinite library containing every possible book. Not just books that exist, but every combination of
characters that _could_ exist.

- Your biography (that hasn't been written yet)
- This exact documentation
- Complete nonsense
- Shakespeare's works

## Implementation

The genius of the original implementation (by Jonathan Basile at [libraryofbabel.info](https://libraryofbabel.info) is that you don't store all these books.
This would require infinite storage. Instead, you use math to calculate

- Given some text -> what is its location in the library?
- Given a location -> what's the text that's there?

It's deterministic: the same text _always_ maps to the same location.

## Key components

### Character Set

- Only 29 characters: space + a-z + comma + period

### Book Structure

- 410 pages per book
- 40 lines per page
- 80 characters per line
- =3,200 characters per page
- =1,312,000 characters per book

### Location System

The library is organized hierarchically:

- Hexagon: A room in the library (can be huge numbers)
- Wall: 4 walls (0-3)
- Shelf: 5 shelves per wall (0-4)
- Volume: 32 books per shelf (0-31)
- Page: 410 pages per book (1-410)

### The Core Algorithm

This is the mathematical magic

- Text -> Number: Treat text as a base29 number
- Number -> Location: Break that huge number into coordinates (hexagon, wall, shelf, volume, page)

#### User Flows

1. Search: User has text and wants to find its location

- User has text and wants to find its location
- Text -> `Base29Encode` -> `*big.Int`
- Show the user the coordinates

2. Browse: User has coordinates and wants to see what text is there

- `Location` -> `LocationToBigInt` -> `*big.Int`
- `*big.Int` -> `Base29Decode` -> `Text`

### Autogenerating Text

Deterministically generating the full 3200 characters of a page based on a Location's coordinates

#### How it works

The key principle: **Deterministic Randomness**

Use the location coordinates as a seed for pseudo-random generation:

- Same location -> Always same "random" content
- Different location -> different content
- No storage needed - can be generated on the fly

#### The algorithm

- Convert location to a seed number
- Initialize random generator with that seed number
- Generate 3200 random indices (0 - 28)
- Map indices to characters
- Result -> deterministic random page
