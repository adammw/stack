package stack

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

// Errorf saves stack trace and pass arguments to
// fmt.Errorf.
func Errorf(format string, a ...interface{}) error {
	return &withStack{
		fmt.Errorf(format, a...),
		callers(),
	}
}

// Origin returns unwrapped origin of the error.
func Origin(err error) error {
	if err == nil {
		return nil
	}

	for {
		e := err
		err = errors.Unwrap(err)
		if err == nil {
			return e
		}
	}
}

// Trace returns stack trace for error.
func Trace(err error) []runtime.Frame {
	stack := make([]runtime.Frame, 0, 30)

	for {
		if err == nil {
			break
		}

		e, ok := err.(*withStack)
		if !ok {
			err = errors.Unwrap(err)
			continue
		}

		stk := e.StackTrace()
		for i := len(stk) - 1; i >= 0; i-- {
			stack = append(stack, stk[i])
		}
		err = errors.Unwrap(err)
	}

	if len(stack) == 0 {
		return nil
	}

	return uniq(stack)
}

func uniq(stack []runtime.Frame) []runtime.Frame {
	seen := make(map[string]struct{}, len(stack))
	j := 0
	for _, frame := range stack {
		v := frame.File + ":" + strconv.Itoa(frame.Line)
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		stack[j] = frame
		j++
	}
	return stack[:j]
}
