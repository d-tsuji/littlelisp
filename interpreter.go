package littlelisp

import (
	"fmt"
)

var (
	library map[interface{}]func(x interface{}) (interface{}, error)
	special map[interface{}]func(interface{}, *Context) (interface{}, error)
)

func init() {
	library = map[interface{}]func(x interface{}) (interface{}, error){
		"first": func(x interface{}) (interface{}, error) {
			xs, ok := x.([]interface{})
			if !ok {
				return nil, fmt.Errorf("fail to type assertion: %v (%T)", x, x)
			}
			return xs[0], nil
		},
		"rest": func(x interface{}) (interface{}, error) {
			xs, ok := x.([]interface{})
			if !ok {
				return nil, fmt.Errorf("fail to type assertion: %v (%T)", x, x)
			}
			return xs[1:], nil
		},
		"print": func(x interface{}) (interface{}, error) {
			return x, nil
		},
	}

	special = map[interface{}]func(interface{}, *Context) (interface{}, error){
		"let": func(input interface{}, ctx *Context) (interface{}, error) {
			s, ok := input.([]interface{})
			if !ok {
				return nil, fmt.Errorf("fail to type assertion: %v (%T)", input, input)
			}
			ss, ok := s[1].([]interface{})
			if !ok {
				return nil, fmt.Errorf("fail to type assertion: %v (%T)", s[1], s[1])
			}
			sp := make(map[interface{}]func(x interface{}) (interface{}, error))
			for _, l := range ss {
				ls, ok := l.([]interface{})
				if !ok {
					return nil, fmt.Errorf("fail to type assertion: %v (%T)", ls, ls)
				}
				sp[ls[0].(Atom).Value] = func(x interface{}) (interface{}, error) {
					return Interpret(x, ctx)
				}
			}
			letCtx := &Context{
				scope:  sp,
				parent: ctx,
			}

			return Interpret(s[2], letCtx)
		},

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
	scope  map[interface{}]func(x interface{}) (interface{}, error)
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
	v, exists := special[input[0].(Atom).Value]
	if exists {
		return v(input, ctx)
	} else {
		list := make([]interface{}, len(input))
		for i, v := range input {
			tmp, err := Interpret(v, ctx)
			if err != nil {
				return nil, err
			}
			list[i] = tmp
		}
		f, ok := list[0].(func(x interface{}) (interface{}, error))
		if ok {
			return f(list[1])
		} else {
			return list, nil
		}
	}
}
