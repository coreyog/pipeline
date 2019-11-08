package pipeline

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

type reflectionPair struct {
	Value reflect.Value
	Type  reflect.Type
}

type Pipeline struct {
	funcs []reflectionPair
}

func New() (pipe *Pipeline) {
	return &Pipeline{
		funcs: []reflectionPair{},
	}
}

func (pipe *Pipeline) PushFunc(fs ...interface{}) (err error) {
	temp := append(pipe.funcs[:0:0], pipe.funcs...)

	for _, f := range fs {
		pair := reflectionPair{
			Type:  reflect.TypeOf(f),
			Value: reflect.ValueOf(f),
		}

		if pair.Type.Kind() != reflect.Func {
			return ErrNotAFunction
		}

		if len(temp) != 0 {
			prev := temp[len(temp)-1]
			if !doParametersMatch(prev.Type, pair.Type) {
				return ErrParameterMismatch
			}
		}

		temp = append(temp, pair)
	}

	pipe.funcs = temp

	return nil
}

func (pipe *Pipeline) Length() int {
	return len(pipe.funcs)
}

func (pipe *Pipeline) Reset() {
	pipe.funcs = []reflectionPair{}
}

func (pipe *Pipeline) Call(args ...interface{}) (results []interface{}, err error) {
	values := make([]reflect.Value, len(args))

	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}

	for _, f := range pipe.funcs {
		if f.Type.NumIn() == len(values)-1 {
			err := values[len(values)-1]

			values = values[:len(values)-1]

			if !err.IsNil() {
				return valueToInterface(values), err.Interface().(error)
			}
		}

		values = f.Value.Call(values)
	}

	results = valueToInterface(values)

	return results, nil
}

func (pipe *Pipeline) String() (s string) {
	names := make([]string, len(pipe.funcs))
	for i, f := range pipe.funcs {
		names[i] = runtime.FuncForPC(f.Value.Pointer()).Name()
	}

	return fmt.Sprintf("[%s]", strings.Join(names, ", "))
}

func doParametersMatch(prev reflect.Type, next reflect.Type) (ok bool) {
	prevOutput := make([]reflect.Type, prev.NumOut())
	nextInput := make([]reflect.Type, next.NumIn())

	for i := range prevOutput {
		prevOutput[i] = prev.Out(i)
	}

	for i := range nextInput {
		nextInput[i] = next.In(i)
	}

	if prev.NumOut() == next.NumIn()+1 {
		potentialErr := prevOutput[len(prevOutput)-1]
		if !potentialErr.Implements(errType) {
			// parameter count mismatch, can't add
			return false
		}

		prevOutput = prevOutput[:len(prevOutput)-1]
	}

	variadicIn := next.IsVariadic()
	if len(prevOutput) != len(nextInput) && !variadicIn {
		return false
	}

	for i := range prevOutput {
		p := prevOutput[i]

		var n reflect.Type
		if i < len(nextInput) {
			n = nextInput[i]
			if i == len(nextInput)-1 && variadicIn {
				// variadic types indicate they're an array of the type,
				// we want the single type, call Elem()
				n = n.Elem()
			}
		} else if variadicIn {
			n = nextInput[len(nextInput)-1].Elem()
		} else {
			// Probably can't happen, playing it safe.
			// This branch only happens if there's a differing
			// count of parameters and they're not variadic
			// which should be caught before this loop.
			return false
		}

		if !p.AssignableTo(n) {
			return false
		}
	}

	return true
}

func valueToInterface(values []reflect.Value) (results []interface{}) {
	results = make([]interface{}, len(values))
	for i, v := range values {
		results[i] = v.Interface()
	}

	return results
}
