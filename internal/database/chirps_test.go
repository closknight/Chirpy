package database

import "testing"

func TestContains(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	if !Contains(nums, 1) {
		t.Error("nums should contain 1")
	}
	if Contains(nums, 100) {
		t.Error("nums should not contain 100")
	}
}

func TestRemoveProfanity(t *testing.T) {
	st1 := "What a kerfuffle"
	exp1 := "What a ****"

	out1 := RemoveProfanity(st1)
	if out1 != exp1 {
		t.Errorf("Output: %s; Expected: %s", out1, exp1)
	}
}
