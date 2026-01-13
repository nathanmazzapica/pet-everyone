package displayname

import (
	"math/rand"
	"regexp"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
		desc  string
	}{
		{"BraveFox0001", true, "letters digits underscore allowed"},
		{"", false, "empty not allowed"},
		{"thisiswaytoolongforthisfield", false, "length limit enforced"},
		{"bad-name", false, "hyphen not allowed"},
	}

	for _, tc := range cases {
		err := Validate(tc.name)
		if tc.valid && err != nil {
			t.Fatalf("expected valid for %s: %v", tc.desc, err)
		}
		if !tc.valid && err == nil {
			t.Fatalf("expected error for %s", tc.desc)
		}
	}
}

func TestGenerateProducesValidName(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	name := Generate(rng)

	if len(name) > MaxLength {
		t.Fatalf("generated name too long: %d", len(name))
	}

	pattern := regexp.MustCompile(`^[a-z]+[A-Z][a-zA-Z]+\d{4}$`)
	if !pattern.MatchString(name) {
		t.Fatalf("generated name has wrong format: %s", name)
	}
}
