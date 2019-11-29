package token

import "testing"

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"Empty", "", false},
		{"Space", " ", false},
		{"SpaceSuffix", "foo ", false},
		{"Number", "123", false},
		{"Keyword", "var", false},

		{"LettersASCII", "foo", true},
		{"MixedASCII", "bar123", true},
		{"UppercaseKeyword", "Func", true},
		{"LettersUnicode", "fóö", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsIdentifier(test.in); got != test.want {
				t.Fatalf("IsIdentifier(%q) = %t, want %v", test.in, got, test.want)
			}
		})
	}
}
