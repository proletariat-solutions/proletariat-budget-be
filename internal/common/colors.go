package common

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

// GetPillColor generates a random color palette for tags based on a given seed
// and returns a slice of strings with the generated colors.
// Seed would be the name for the tag, to keep it pseudo-random instead of totally random.
func GetPillColor(seed string) (*[]string, error) {
	// transform the seed into an integer for the seed
	h := fnv.New32a()
	_, err := h.Write([]byte(seed)) // Fixed variable name from 's' to 'seed'
	if err != nil {
		return nil, err
	}
	seedInt := h.Sum32()

	// Create a new random source with the seed
	source := rand.NewSource(int64(seedInt))
	rng := rand.New(source)

	// Generate a random color in hex format (#RRGGBB)
	r := rng.Intn(256)
	g := rng.Intn(256)
	b := rng.Intn(256)
	color := fmt.Sprintf("#%02x%02x%02x", r, g, b)

	// Generate the opposite color for background
	oppositeR := 255 - r
	oppositeG := 255 - g
	oppositeB := 255 - b
	oppositeColor := fmt.Sprintf("#%02x%02x%02x", oppositeR, oppositeG, oppositeB)

	// Return both colors
	colors := []string{color, oppositeColor}
	return &colors, nil
}
