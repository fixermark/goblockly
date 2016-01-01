// Utilities for test evaluation

package goblockly

import (
	"runtime"
	"testing"
)

// checkFailInterpreter is a testing utility to verify that the interpreter
// failed while running test. Interpreter failure is represented by a panic, so
// we're just catching panics here.
func checkFailInterpreter(t *testing.T, test func()) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	test()

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "UNKNOWN"
		line = 0
	}
	t.Errorf("%s:%d - Failure not called on interpreter.", file, line)
}
