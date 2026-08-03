package main

import (
	"strings"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"empty", "", true},
		{"single", "x", true},
		{"classic", "A man a plan a canal Panama", true},
		{"not palindrome", "go", false},
		{"spaces only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.in); got != tt.want {
				t.Fatalf("IsPalindrome(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	input := "A man a plan a canal Panama"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsPalindrome(input)
	}
}

func FuzzIsPalindrome(f *testing.F) {
	f.Add("aba")
	f.Add("not")
	f.Fuzz(func(t *testing.T, s string) {
		got := IsPalindrome(s)
		compact := stringsToLowerNoSpace(s)
		want := compact == reverse(compact)
		if got != want {
			t.Fatalf("IsPalindrome(%q)=%v, independent check=%v", s, got, want)
		}
	})
}

func stringsToLowerNoSpace(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == ' ' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		b.WriteRune(r)
	}
	return b.String()
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}