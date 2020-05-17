package littlelisp

import (
	"fmt"
)

var (
	library map[interface{}]func(x []interface{}) interface{}
	special map[interface{}]func(interface{}, *Context) (interface{}, error)
)

func init() {
	library = map[interface{}]func(x []interface{}) interface{}{
		"first": func(x []interface{}) interface{} {
			return x[0]
		},
		"rest": func(x []interface{}) interface{} {
			return x[1:]
		},
		"print": func(x []interface{}) interface{} {
			fmt.Println(x)
			return nil
		},
	}

	special = map[interface{}]func(interface{}, *Context) (interface{}, error){
		// A literal 1 is assumed to be true.
		"if": func(input interface{}, ctx *Context) (interface{}, error) {
			s, ok := input.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid input: %v", s)
			}
			n, err := Interpret(s[1], ctx)
			if err != nil {
				return nil, err
			}
			if n == int64(1) {
				return Interpret(s[2], ctx)
			} else {
				return Interpret(s[3], ctx)
			}
		},
	}
}

type Context struct {
	scope  map[interface{}]func(x []interface{}) interface{}
	parent *Context
}

func (c *Context) get(identifier interface{}) interface{} {
	v, exists := c.scope[identifier]
	if exists {
		return v
	} else if c.parent != nil {
		return c.parent.get(identifier)
	}
	return nil
}

func Interpret(input interface{}, ctx *Context) (interface{}, error) {
	if ctx == nil {
		return Interpret(input, &Context{
			scope:  library,
			parent: nil,
		})
	}
	s, ok := input.([]interface{})
	if ok {
		return interpretList(s, ctx)
	} else {
		a, ok := input.(Atom)
		if !ok {
			return nil, fmt.Errorf("input is not Atom: %v", input)
		}
		if a.TokenType == "identifier" {
			return ctx.get(a.Value), nil
		} else if a.TokenType == "literal" {
			return a.Value, nil
		}
	}

	return nil, nil
}

func interpretList(input []interface{}, ctx *Context) (interface{}, error) {
	fmt.Println(input[0])
	v, exists := special[input[0].(Atom).Value]
	if exists {
		return v(input, ctx)
	} else {
		// TODO
	}
	return nil, nil
}
