package inject_test

import "testing"

// testLogger is useful when debugging tests.
type testLogger struct {
	t *testing.T
}

func (t testLogger) Debugf(f string, args ...interface{}) {
	t.t.Logf(f, args...)
}
