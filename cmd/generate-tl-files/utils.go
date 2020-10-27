package main

import (
	"strings"
	"unicode"
)

func goify(name string) string {
	repl := []string{
		"_", "",
		" ", "",
		".", "",
		"P2p", "P2P",
	}

	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
	for i, r := range runes {
		if r == '_' || r == ' ' || r == '.' {
			if i+1 == len(runes) {
				break
			}
			runes[i+1] = unicode.ToUpper(runes[i+1])
		}
	}

	return strings.NewReplacer(repl...).Replace(string(runes))
}
