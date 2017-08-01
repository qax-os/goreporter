package aligncheck

import (
	"testing"
)

func Test_AlignCheck(t *testing.T) {
	LinterAligncheck{}.ComputeMetric("net/http")
}
