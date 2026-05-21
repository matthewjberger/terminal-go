// Package parse turns a raw player input line into a normalized
// verb and argument list.
package parse

import "strings"

// Token is the result of tokenizing one input line. Verb is
// lowercased and normalized through the alias table (so "n",
// "look", "x", "i" all reduce to canonical verbs). Args is the
// remaining words with articles and prepositions stripped.
// Object is Args joined by a single space.
type Token struct {
	Verb   string
	Args   []string
	Object string
}

// Tokenize lowercases the input, splits on whitespace, normalizes
// the verb, and strips filler words ("a", "an", "the", "at",
// "on", "in", "to") from the remaining arguments.
func Tokenize(line string) Token {
	fields := strings.Fields(strings.ToLower(line))
	if len(fields) == 0 {
		return Token{}
	}
	verb := normalizeVerb(fields[0])

	args := make([]string, 0, len(fields)-1)
	for _, word := range fields[1:] {
		if isFiller(word) {
			continue
		}
		args = append(args, word)
	}
	return Token{
		Verb:   verb,
		Args:   args,
		Object: strings.Join(args, " "),
	}
}

func normalizeVerb(word string) string {
	switch word {
	case "l":
		return "look"
	case "i", "inv":
		return "inventory"
	case "x":
		return "examine"
	case "n":
		return "north"
	case "s":
		return "south"
	case "e":
		return "east"
	case "w":
		return "west"
	case "u":
		return "up"
	case "d":
		return "down"
	case "g":
		return "go"
	case "get", "grab", "pick":
		return "take"
	case "put":
		return "drop"
	case "q":
		return "quit"
	case "h", "?":
		return "help"
	}
	return word
}

func isFiller(word string) bool {
	switch word {
	case "a", "an", "the", "at", "on", "in", "to":
		return true
	}
	return false
}
