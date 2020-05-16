package littlelisp

import "strings"

func Tokenize(input string) []string {
	ret := strings.ReplaceAll(input, "(", " ( ")
	ret = strings.ReplaceAll(ret, ")", " ) ")
	return strings.Fields(ret)
}
