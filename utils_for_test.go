// Utilities for test evaluation

package goblockly

import (
	"runtime"
	"testing"
)

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
