package jsonfile

import (
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string with trailing comment",
			input:    `"foo bar" // some comment`,
			expected: `"foo bar"`,
		},
		{
			name: "object with comments",
			input: `{
				"foo": 123, // bla bla bla
				"bar": "baz",
				"bla": 23423 // asd
			}`,
			expected: `{"foo":123,"bar":"baz","bla":23423}`,
		},
		{
			name: "leading comment",
			input: `// leading comment
			{
				"foo": 123, // bla bla bla
				"bar": "baz",
				"bla": 23423 // asd
			}`,
			expected: `{"foo":123,"bar":"baz","bla":23423}`,
		},
		{
			name: "trailing comment",
			input: `{
				"foo": 123,
				"bar": "baz",
				"bla": 23423
			}// trailing comment
			// trailing comment`,
			expected: `{"foo":123,"bar":"baz","bla":23423}`,
		},
		{
			name: "multiple comment lines inside object",
			input: `
			{
				"foo": 123,
				"bar": "baz",
				// asd
				// foo
				//
				//bas
				"bla": 23423
			}`,
			expected: `{"foo":123,"bar":"baz","bla":23423}`,
		},
		{
			name:     "more than two slashes",
			input:    `"asd asd asd" ////////`,
			expected: `"asd asd asd"`,
		},
	}

	for _, test := range tests {
		in := &commentFilter{file: strings.NewReader(test.input)}
		out := new(strings.Builder)

		_, err := io.Copy(out, in)
		if err != nil {
			t.Error("test", test.name, "failed:", err)
		}

		output := out.String()
		if output != test.expected {
			t.Error("test", test.name, "failed: expected", test.expected, "but got", output)
		}
	}
}
