package main

import (
	"fmt"
	"testing"
)

func TestPerm(t *testing.T) {
	for n, tc := range []struct {
		s      string
		length int
		want   []string
	}{
		{"", 0, []string{""}},
		{"", 1, []string{""}},

		{"a", 0, []string{""}},
		{"a", 1, []string{"a"}},
		{"a", 2, []string{"aa"}},
		{"a", 3, []string{"aaa"}},

		{"ab", 0, []string{""}},
		{"ab", 1, []string{"a", "b"}},
		{"ab", 2, []string{"aa", "ab", "ba", "bb"}},
		{"ab", 3, []string{"aaa", "aab", "aba", "abb", "baa", "bab", "bba", "bbb"}},
	} {
		var got []string

		perm(tc.s, tc.length, func(s string) {
			got = append(got, s)
		})

		if fmt.Sprint(got) != fmt.Sprint(tc.want) {
			t.Errorf("[%d] got %v, want %v", n, got, tc.want)
		}
	}
}
