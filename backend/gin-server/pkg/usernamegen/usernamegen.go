package usernamegen

import (
	_ "embed"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//go:embed data/adverbs.txt
var adverbsData string

//go:embed data/adjectives.txt
var adjectivesData string

//go:embed data/animals.txt
var animalsData string

var (
	adverbs    []string
	adjectives []string
	animals    []string
)

func init() {
	adverbs = strings.Split(strings.TrimSpace(adverbsData), "\n")
	adjectives = strings.Split(strings.TrimSpace(adjectivesData), "\n")
	animals = strings.Split(strings.TrimSpace(animalsData), "\n")
}

// Generate returns a unique username like "swift-sassy-panther-1712938493"
func Generate() string {
	adverb := adverbs[rand.Intn(len(adverbs))]
	adjective := adjectives[rand.Intn(len(adjectives))]
	animal := animals[rand.Intn(len(animals))]
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%s-%s-%d", adverb, adjective, animal, timestamp)
}
