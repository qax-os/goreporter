package depend

import (
	"fmt"
	"testing"
)

func Test_Depend(t *testing.T) {
	h := Depend("../../../goreporter", "")
	fmt.Println(h)
}
