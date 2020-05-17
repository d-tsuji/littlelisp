package littlelisp

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"normal",
			`((lambda (x) x) "Lisp")`,
			[]string{"(", "(", "lambda", "(", "x", ")", "x", ")", `"Lisp"`, ")"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tokenize(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParenthesize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  interface{}
	}{
		{
			"sample",
			`((lambda (x) x) "Lisp")`,
			[]interface{}{
				[]interface{}{Atom{"identifier", "lambda"}, []interface{}{Atom{"identifier", "x"}}, Atom{"identifier", "x"}},
				Atom{"literal", "Lisp"},
			},
		},
		{
			"should lex a single atom",
			`a`,
			Atom{"identifier", "a"},
		},
		{
			"should lex an atom in a list",
			`()`,
			[]interface{}{},
		},
		{
			"should lex multi atom list",
			`(hi you)`,
			[]interface{}{Atom{"identifier", "hi"}, Atom{"identifier", "you"}},
		},
		{
			"should lex list containing list",
			`((x))`,
			[]interface{}{[]interface{}{Atom{"identifier", "x"}}},
		},
		{
			"should lex list containing list",
			`(x (x))`,
			[]interface{}{Atom{"identifier", "x"}, []interface{}{Atom{"identifier", "x"}}},
		},
		{
			"should lex list containing list",
			`(x y)`,
			[]interface{}{Atom{"identifier", "x"}, Atom{"identifier", "y"}},
		},
		{
			"should lex list containing list",
			"(x (y) z)",
			[]interface{}{Atom{"identifier", "x"}, []interface{}{Atom{"identifier", "y"}}, Atom{"identifier", "z"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Parenthesize(Tokenize(tt.input))
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Parenthesize() differs: (-got +want)\n%s", diff)
			}
		})
	}
}
