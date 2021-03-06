package globals

import (
	"fmt"
)

type intFunc func(a, b int) int

var (
	intFuncs = map[string]intFunc{
		"plus":   func(a, b int) int { return a + b },
		"minus":  func(a, b int) int { return a - b },
		"times":  func(a, b int) int { return a * b },
		"divide": func(a, b int) int { return a / b },
		"modulo": func(a, b int) int { return a % b },
	}
)

func op(name string, args []interface{}) (interface{}, error) {
	// These two operator can be used as a sign changer for numbers
	// (so only one argument). Anyway, who uses the + as a unary function??
	if len(args) == 1 && (name == "plus" || name == "minus") {
		args = append([]interface{}{0}, args...)
	}

	// Check that two arguments are provided at least
	if len(args) < 2 {
		return 0, fmt.Errorf("at least two params are needed for the %s operator", name)
	}

	// Try to sum the integers
	ac, ok := args[0].(int)
	if ok {
		f := intFuncs[name]
		for _, arg := range args[1:] {
			n := arg.(int)
			if name == "divide" && n == 0 {
				return ac, fmt.Errorf("division by zero")
			}

			ac = f(ac, n)
		}
		return ac, nil
	}

	return 0, fmt.Errorf("%s operator can't handle this kind of numbers", name)
}

func Plus(args ...interface{}) (interface{}, error) {
	return op("plus", args)
}

func Minus(args ...interface{}) (interface{}, error) {
	return op("minus", args)
}

func Times(args ...interface{}) (interface{}, error) {
	return op("times", args)
}

func Divide(args ...interface{}) (interface{}, error) {
	return op("divide", args)
}

func Modulo(args ...interface{}) (interface{}, error) {
	return op("modulo", args)
}

func GreaterThan(a, b int) bool {
	return a > b
}

func GreaterEqual(a, b int) bool {
	return a >= b
}

func LessThan(a, b int) bool {
	return a < b
}

func LessEqual(a, b int) bool {
	return a <= b
}

func Equal(a, b interface{}) bool {
	return a == b
}

func Not(a bool) bool {
	return !a
}
