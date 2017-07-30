package staticcheck

import (
	"testing"
)

func TestStaticCheck(t *testing.T) {
	path := make(map[string]string, 0)
	path["../../engine"] = "../../engine"
	StaticCheck(path)
}
