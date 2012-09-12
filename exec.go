package main

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	zero      reflect.Value
)

type variables map[string]reflect.Value
type functions map[string]reflect.Value

type state struct {
	funcs  functions
	vars   variables
	output io.Writer
	t      *ListNode
}

func (s *state) recover(errp *error) {
	if e := recover(); e != nil {
		*errp = fmt.Errorf("%s", e)
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
	}
}

func (s *state) errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func (s *state) walkNode(n Node) reflect.Value {
	switch n.Type() {
	case NodeCall:
		return s.makeCall(n.(*CallNode))
	}

	return zero
}

func (s *state) makeCall(n *CallNode) reflect.Value {
	// Get the func in the index
	f, ok := s.funcs[n.Name]
	if !ok {
		panic(fmt.Errorf("function not defined %s", n.Name))
	}

	// Analyze its type
	t := f.Type()
	numArgs := t.NumIn()

	// Check if the number of args it's correct
	if t.IsVariadic() {
		numArgs -= 1
		if len(n.Args) < numArgs {
			s.errorf("wrong number of args for %s: want at least %d, got %d", n.Name,
				numArgs, len(n.Args))
		}
	} else if len(n.Args) != numArgs {
		s.errorf("wrong number of args for %s: want %d, got %d", n.Name, numArgs, len(n.Args))
	}

	// Check if the function return it's correct
	s.checkFuncReturn(n.Name, t)

	// Prepare the arguments array
	nargs := numArgs
	if t.IsVariadic() {
		nargs += len(n.Args) - nargs
	}
	args := make([]reflect.Value, nargs)

	// Add the fixed arguments
	i := 0
	for ; i < numArgs; i++ {
		args[i] = s.evalArg(t.In(i), n.Args[i])
	}

	// Add the variadic arguments
	if t.IsVariadic() {
		argType := t.In(numArgs).Elem()
		for ; i < len(n.Args); i++ {
			args[i] = s.evalArg(argType, n.Args[i])
		}
	}

	res := f.Call(args)

	if len(res) == 2 && !res[1].IsNil() {
		s.errorf("error calling %s: %s", n.Name, res[1].Interface().(error))
	}

	return res[0]
}

func (s *state) checkFuncReturn(name string, t reflect.Type) {
	// Check the number of return values
	if t.NumOut() == 1 || (t.NumOut() == 2 && t.Out(1) == errorType) {
		return
	}

	s.errorf("can't handle multiple returns from function %s", name)
}

func (s *state) evalArg(t reflect.Type, n Node) reflect.Value {
	// If the arg it's a subtree, execute it first
	switch n.Type() {
	case NodeCall:
		return s.makeCall(n.(*CallNode))
	}

	// Return the correct value depending on the needed type
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, ok := n.(*NumberNode); ok && n.IsInt {
			v := reflect.New(t).Elem()
			v.SetInt(n.Int64)
			return v
		}
		s.errorf("expected integer; found %s", n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, ok := n.(*NumberNode); ok && n.IsUint {
			v := reflect.New(t).Elem()
			v.SetUint(n.Uint64)
			return v
		}
		s.errorf("expected unsigned integer; found %s", n)

	case reflect.String:
		v := reflect.New(t).Elem()
		v.SetString(n.(*StringNode).Text)
		return v

	case reflect.Interface:
		return s.evalEmptyInterface(t, n)
	}

	// Can't handle that type of arguments
	s.errorf("can't handle %+v for arg of type %s", n, t)
	panic("not reached")
}

func (s *state) evalEmptyInterface(t reflect.Type, n Node) reflect.Value {
	// Depending on the node type, try to guess the best arg
	switch n := n.(type) {
	case *NumberNode:
		return s.idealConstant(n)

	case *StringNode:
		return reflect.ValueOf(n.Text)
	}

	// Can't handle this kind of node
	s.errorf("can't handle assignment of %s to empty interface argument", n)
	panic("not reached")
}

// Try to parse nums as ideal constant. Note that an unsigned integer ideal
// constant should be an integer one too.
func (s *state) idealConstant(c *NumberNode) reflect.Value {
	switch {
	case c.IsInt:
		n := int(c.Int64)
		if int64(n) != c.Int64 {
			s.errorf("%s overflows int", c.Text)
		}
		return reflect.ValueOf(n)

	case c.IsUint:
		s.errorf("unsigned integers are not supported: %s", c.Text)
	}

	return zero
}

func (s *state) print(v reflect.Value) {
	if v.Kind() == reflect.String {
		fmt.Fprint(s.output, v.Interface())
		return
	}

	fmt.Fprintln(s.output, v.Interface())
}

func Exec(output io.Writer, tree *ListNode, funcs map[string]interface{}) (err error) {
	// Convert the functions to reflect values
	f := map[string]reflect.Value{}
	for name, fn := range funcs {
		f[name] = reflect.ValueOf(fn)
	}

	// Build the environment
	s := &state{
		vars:   make(variables),
		funcs:  f,
		output: output,
		t:      tree,
	}

	// Hook up the recover
	defer s.recover(&err)

	// Start the execution
	for _, n := range s.t.Nodes {
		s.print(s.walkNode(n))
	}

	return nil
}
