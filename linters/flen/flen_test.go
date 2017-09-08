package flen

import (
	"os"
	"reflect"
	"testing"
)

var (
	sampleFuncLens *FuncLens
	sampleCode     = `
func single() {
	println("hello single")
}

func double() {
	println("hello double")
	println("hello double")
}

func trouble() {
	println("hello trouble")
	println("hello trouble")
	println("hello trouble")
}
`
	lengAllWant  = []int{3, 0, 1, 2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0}
	funcLengWant = [][]string{
		[]string{"0", "Len", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "81", "1"},
		[]string{"1", "Swap", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "95", "1"},
		[]string{"2", "rangeAsked", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "50", "1"},
		[]string{"3", "Less", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "82", "11"},
		[]string{"4", "getPkgPath", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "205", "17"},
		[]string{"5", "FuncLen", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "52", "20"},
		[]string{"6", "computeHistogram", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "98", "26"},
		[]string{"7", "GenerateFuncLens", "/Users/wangguoliang/Documents/wgliang/src/github.com/360EntSecGroup-Skylar/goreporter/linters/flen/flen.go", "130", "70"},
	}
)

func TestFuncLensLen(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Receiver.
		flens *FuncLens
		// Expected results.
		want int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.flens.Len(); got != tt.want {
			t.Errorf("%q. FuncLens.Len() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFuncLensLess(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Receiver.
		flens *FuncLens
		// Parameters.
		i int
		j int
		// Expected results.
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.flens.Less(tt.i, tt.j); got != tt.want {
			t.Errorf("%q. FuncLens.Less() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFuncLensSwap(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Receiver.
		flens *FuncLens
		// Parameters.
		i int
		j int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.flens.Swap(tt.i, tt.j)
	}
}

func TestFuncLensComputeHistogram(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Receiver.
		flens *FuncLens
		// Expected results.
		want []int
	}{}
	for _, tt := range tests {
		if got := tt.flens.computeHistogram(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. FuncLens.computeHistogram() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGenerateFuncLens(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		pkg     string
		options *Options
		// Expected results.
		want    FuncLens
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, _, err := GenerateFuncLens(tt.pkg, tt.options)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GenerateFuncLens() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GenerateFuncLens() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetPkgPath(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		pkgname string
		// Expected results.
		want  string
		want1 *os.PathError
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, got1 := getPkgPath(tt.pkgname)
		if got != tt.want {
			t.Errorf("%q. getPkgPath() got = %v, want %v", tt.name, got, tt.want)
		}
		if !reflect.DeepEqual(got1, tt.want1) {
			t.Errorf("%q. getPkgPath() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}

func TestFuncLen(t *testing.T) {
	lengAll, funcLeng := FuncLen("github.com/360EntSecGroup-Skylar/goreporter/linters/flen")
	if !reflect.DeepEqual(lengAll, lengAllWant) {
		t.Errorf("want %v, but got %v", lengAll, lengAllWant)
	}

	if !reflect.DeepEqual(funcLeng, funcLengWant) {
		t.Errorf("want %v, but got %v", funcLeng, funcLengWant)
	}
}
