package data

import "strings"

// Reverse returns the reversed string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToUpperCase converts all characters to uppercase
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase converts all characters to lowercase
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// IsPalindrome checks if a string is a palindrome
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	return s == Reverse(s)
}

// CountVowels counts the number of vowels in a string
func CountVowels(s string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range s {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}
	return count
}

// TruncateString truncates a string to maxLen using Min
func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	truncLen := Min(len(runes), maxLen)
	return string(runes[:truncLen])
}

// RepeatString repeats a string count times, ensuring non-negative count with Max
func RepeatString(s string, count int) string {
	safeCount := Max(count, 0)
	return strings.Repeat(s, safeCount)
}

// CountDifference returns the absolute difference in lengths of two strings
func CountDifference(s1, s2 string) int {
	return Abs(len(s1) - len(s2))
}
