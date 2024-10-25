# Fizzgig Search

Fizzgig Search is a simple fuzzy search library, written in Go, that calculates similarities between a search term and a list of possible answers, given a certain threshold.

## Usage

```
import (
    "fmt"
    "fizzgig-search/search"
)

func main() {
	matches := search.Search("apple", []string{"pineapple", "kiwi"}, 0.5)
	fmt.Printf("Matches: %v", matches)
}

// Matches: {apple [{pineapple 0.56}] 1}
```