package rest

import "testing"

func Test_trimDuplicate(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{
			input:  "",
			output: "",
		},
		{
			input:  "/",
			output: "/",
		},
		{
			input:  "//",
			output: "/",
		},
		{
			input:  "///",
			output: "/",
		},
		{
			input:  "/api/v0/resource",
			output: "/api/v0/resource",
		},
		{
			input:  "/api//v0/resource",
			output: "/api/v0/resource",
		},
		{
			input:  "//api/v0/resource",
			output: "/api/v0/resource",
		},
		{
			input:  "//api//v0//resource/",
			output: "/api/v0/resource/",
		},
	}

	for _, c := range cases {
		output := trimDuplicateRune(c.input, '/')
		if output != c.output {
			t.Errorf("input: %s, expect: %s, actual: %s",
				c.input, c.output, output)
		}
	}
}
