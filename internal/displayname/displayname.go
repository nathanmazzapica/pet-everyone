package displayname

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/TwiN/go-away"
)

const MaxLength = 20

var (
	errEmptyName      = errors.New("display name cannot be empty")
	errTooLong        = fmt.Errorf("display name must be at most %d characters", MaxLength)
	errInvalidChars   = errors.New("display name may only contain letters, digits, or underscores")
	errProfaneOrSlurs = errors.New("name cannot contain profanity or slurs")

	allowedPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

	adjectives = []string{"brave", "calm", "swift", "bright", "happy", "kind", "quick", "sharp", "sunny", "lucky"}
	nouns      = []string{"Fox", "Cat", "Dog", "Hawk", "Panda", "Goat", "Wolf", "Lynx", "Seal", "Crab"}
)

// Generate returns a display name in the format adjective + Noun + 4-digit number.
// If rng is nil, a new source seeded with the current time is used.
func Generate(rng *rand.Rand) string {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	for {
		adjective := adjectives[rng.Intn(len(adjectives))]
		noun := nouns[rng.Intn(len(nouns))]
		number := fmt.Sprintf("%04d", rng.Intn(10000))

		candidate := adjective + noun + number
		if len(candidate) <= MaxLength {
			return candidate
		}
	}
}

// Validate checks length and allowed characters for a display name.
func Validate(name string) error {
	if name == "" {
		return errEmptyName
	}
	if len(name) > MaxLength {
		return errTooLong
	}
	if !allowedPattern.MatchString(name) {
		return errInvalidChars
	}
	if containsProfanity(name) {
		return errProfaneOrSlurs
	}
	return nil
}

// Normalize lower-cases a display name for case-insensitive comparisons.
func Normalize(name string) string {
	return strings.ToLower(name)
}

func containsProfanity(name string) bool {
	return goaway.IsProfane(name)
}
