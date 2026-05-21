// Package parse turns a raw player input line into a normalized
// verb and argument list.
package parse

import "strings"

// Token is the result of parsing one input line. Verb is the
// normalized lowercased canonical verb. Args / Object are
// lowercased for matching against item names. RawArgs /
// RawObject preserve the original case for things like
// filenames passed to save/load.
type Token struct {
	Verb      string
	Args      []string
	Object    string
	RawArgs   []string
	RawObject string
}

func Tokenize(line string) Token {
	raw := strings.Fields(line)
	if len(raw) == 0 {
		return Token{}
	}
	lower := make([]string, len(raw))
	for i, word := range raw {
		lower[i] = strings.ToLower(word)
	}
	verb, consumed := normalizeVerb(lower)

	args := make([]string, 0, len(raw)-consumed)
	rawArgs := make([]string, 0, len(raw)-consumed)
	for i := consumed; i < len(raw); i++ {
		if isFiller(lower[i]) {
			continue
		}
		args = append(args, lower[i])
		rawArgs = append(rawArgs, raw[i])
	}
	return Token{
		Verb:      verb,
		Args:      args,
		Object:    strings.Join(args, " "),
		RawArgs:   rawArgs,
		RawObject: strings.Join(rawArgs, " "),
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
	case "get", "grab":
		return "take", 1
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
