package structtag_test

import (
	"testing"

	"github.com/facebookgo/structtag"
)

func TestExtract(t *testing.T) {
	cases := []struct {
		Name  string // Input name
		Tag   string // Input tag
		Found bool   // Expected found status
		Value string // Expected value
		Error bool   // Indicates if an error is expected
	}{
		{
			Name:  "inject",
			Tag:   `inject:`,
			Error: true,
		},
		{
			Name:  "inject",
			Tag:   `inject:"`,
			Error: true,
		},
		{
			Name:  "inject",
			Tag:   `inject:""`,
			Found: true,
		},
		{
			Name:  "inject",
			Tag:   `inject:"a"`,
			Found: true,
			Value: "a",
		},
		{
			Name:  "inject",
			Tag:   ` inject:"a"`,
			Found: true,
			Value: "a",
		},
		{
			Name: "inject",
			Tag:  `  `,
		},
		{
			Name:  "inject",
			Tag:   `inject:"\"a"`,
			Found: true,
			Value: `"a`,
		},
	}

	for _, e := range cases {
		found, value, err := structtag.Extract(e.Name, e.Tag)

		if !e.Error && err != nil {
			t.Fatalf("unexpected error %s for case %+v", err, e)
		}
		if e.Error && err == nil {
			t.Fatalf("did not get expected error for case %+v", e)
		}

		if found != e.Found {
			if e.Found {
				t.Fatalf("did not find value when expecting to %+v", e)
			} else {
				t.Fatalf("found value when not expecting to %+v", e)
			}
		}

		if value != e.Value {
			t.Fatalf(`found unexpected value "%s" for %+v`, value, e)
		}
	}
}
