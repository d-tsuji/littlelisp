package littlelisp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInterpretAtom(t *testing.T) {
	type args struct {
		input string
		ctx   *Context
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"should return string atom", args{input: `"a"`, ctx: nil}, "a", false},
		//{"should return string with space atom", args{input: `"a b"`, ctx: nil}, "a b", false},
		{"should return number atom", args{input: `123`, ctx: nil}, int64(123), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Parser(tt.args.input)
			got, err := Interpret(input, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Interpret() differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestInterpretIf(t *testing.T) {
	type args struct {
		input string
		ctx   *Context
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"should choose the right branch", args{input: `(if 1 42 4711)`, ctx: nil}, int64(42), false},
		{"should choose the right branch", args{input: `(if 0 42 4711)`, ctx: nil}, int64(4711), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Parser(tt.args.input)
			got, err := Interpret(input, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Interpret() differs: (-got +want)\n%s", diff)
			}
		})
	}
}
