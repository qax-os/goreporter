package simpler

import (
	"testing"
)

func TestSimpler(t *testing.T) {
	path := make(map[string]string, 0)
	path["../../engine"] = "../../engine"
	Simpler(path)
}
