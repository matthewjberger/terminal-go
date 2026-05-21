// Package parse turns a raw player input line into a normalized
// verb and argument list.
package parse

import "strings"

type Token struct {
	Verb   string
	Args   []string
	Object string
}

func Tokenize(line string) Token {
	fields := strings.Fields(strings.ToLower(line))
	if len(fields) == 0 {
		return Token{}
	}
	verb, consumed := normalizeVerb(fields)

	args := make([]string, 0, len(fields)-consumed)
	for _, word := range fields[consumed:] {
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

func normalizeVerb(fields []string) (verb string, consumed int) {
	if len(fields) >= 2 {
		switch fields[0] + " " + fields[1] {
		case "pick up":
			return "take", 2
		case "put down":
			return "drop", 2
		}
	}
	switch fields[0] {
	case "l":
		return "look", 1
	case "i", "inv":
		return "inventory", 1
	case "x":
		return "examine", 1
	case "n":
		return "north", 1
	case "s":
		return "south", 1
	case "e":
		return "east", 1
	case "w":
		return "west", 1
	case "u":
		return "up", 1
	case "d":
		return "down", 1
	case "g":
		return "go", 1
	case "get", "grab", "pick":
		return "take", 1
	case "put":
		return "drop", 1
	case "q":
		return "quit", 1
	case "h", "?":
		return "help", 1
	}
	return fields[0], 1
}

func isFiller(word string) bool {
	switch word {
	case "a", "an", "the", "at", "on", "in", "to":
		return true
	}
	return false
}
