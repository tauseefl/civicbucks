package main

import "testing"

func TestPalindrome(t *testing.T) {

	validPalindrome := 43234
	result := isPalindrome(validPalindrome)
	if !result {
		t.Errorf("Expected valid palindrome")
	}

	invalidPalindrome := 123456
	result2 := isPalindrome(invalidPalindrome)
	if result2 {
		t.Errorf("Expected invalid palindrome")
	}

}

func TestBinaryPalindrome(t *testing.T) {

	validPalindrome := 15351
	result := isBinaryPalindrome(validPalindrome)
	if !result {
		t.Errorf("Expected valid palindrome")
	}

	invalidPalindrome := 43234
	result2 := isBinaryPalindrome(invalidPalindrome)
	if result2 {
		t.Errorf("Expected invalid palindrome")
	}

}

func TestReverseString(t *testing.T) {

	forward := "ABCDEFGHIJKLMNOP"
	expected_reverse := "PONMLKJIHGFEDCBA"
	result := reverse(forward)
	if expected_reverse != result {
		t.Errorf("Expected reversed string")
	}

}
