package littlelisp

import (
	"log"
	"strings"
)

func Tokenize(input string) []string {
	ret := strings.ReplaceAll(input, "(", " ( ")
	ret = strings.ReplaceAll(ret, ")", " ) ")
	return strings.Fields(ret)
}

type Atom struct {
	TokenType string
	Value     string
}

func Parenthesize(tokens []string) (interface{}, []string) {
	if len(tokens) == 0 {
		log.Fatal("unexpected EOF while reading")
	}
	token := tokens[0]
	tokens = tokens[1:]
	if token == "(" {
		// It returns an empty slice instead of a nil slice. It is not essential.
		l := []interface{}{}
		for tokens[0] != ")" {
			s, out := Parenthesize(tokens)
			tokens = out
			l = append(l, s)
		}
		_ = tokens[0]
		tokens = tokens[1:]
		return l, tokens
	} else if token == ")" {
		log.Fatal("parenthesize error: unexpected )")
	} else {
		return Categorize([]rune(token)), tokens
	}
	return struct{}{}, nil
}

func Categorize(input []rune) Atom {
	if input[0] == '"' && input[len(input)-1] == '"' {
		return Atom{
			TokenType: "literal",
			Value:     string(input[1 : len(input)-1]),
		}
	} else {
		return Atom{
			TokenType: "identifier",
			Value:     string(input),
		}
	}
}
