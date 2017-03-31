package depend

import (
	"fmt"
	"testing"
)

func Test_Dep(t *testing.T) {
	h := Dep("../gotest", "")
	fmt.Println(h)
}
