package main

import "testing"

func TestSanitize(t *testing.T) {
	st1 := "What a kerfuffle"
	exp1 := "What a ****"

	out1 := sanitize(st1)
	if out1 != exp1 {
		t.Errorf("Output: %s; Expected: %s", out1, exp1)
	}
}
