package pipeline

import (
	"reflect"
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
	pairs := make([]reflectionPair, len(fs))

	for i, f := range fs {
		pair := reflectionPair{
			Type:  reflect.TypeOf(f),
			Value: reflect.ValueOf(f),
		}

		if pair.Type.Kind() != reflect.Func {
			return ErrNotAFunction
		}

		if len(pipe.funcs) != 0 {
			prev := pipe.funcs[len(pipe.funcs)-1]
			if !doParametersMatch(prev.Type, pair.Type) {
				return ErrParameterMismatch
			}
		}

		pairs[i] = pair
	}

	pipe.funcs = append(pipe.funcs, pairs...)

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

	if len(prevOutput) != len(nextInput) {
		return false
	}

	for i := range prevOutput {
		p := prevOutput[i]
		n := nextInput[i]

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
