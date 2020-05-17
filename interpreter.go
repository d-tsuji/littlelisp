package littlelisp

import (
	"fmt"
)

type Context struct {
	scope  map[interface{}]func(x []interface{}) interface{}
	parent *Context
}

var library = map[interface{}]func(x []interface{}) interface{}{
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
		// TODO
		fmt.Println(s)
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
