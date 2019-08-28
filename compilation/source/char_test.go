package source

import (
	"testing"
)

const (
	specials      = "{}[]<>'\""
	numbers       = "123456789"
	lowerLetters  = "abcdefghijklmnopqrstuvwxyz"
	upperLetters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letters       = lowerLetters + upperLetters
	alphanumerics = numbers + letters
)

// TestIsInRange tests whether the IsInRange method is inclusive.
func TestIsInRange(test *testing.T) {
	const lower, upper = 'a', 'c'

	for _, letter := range "abc" {
		if Char(letter).IsInRange(lower, upper) {
			continue
		}
		test.Errorf("%c is not in range (%c..%c)", letter, lower, upper)
	}
}

func TestSpecialsAreNotAlphanumeric(test *testing.T) {
	for _, special := range specials {
		if !Char(special).IsAlphanumeric() {
			continue
		}
		test.Errorf("%c is considered alphanumeric", special)
	}
}

func TestNumbersAndLettersAreAlphanumeric(test *testing.T) {
	for _, alphanumeric := range alphanumerics {
		if Char(alphanumeric).IsAlphanumeric() {
			continue
		}
		test.Errorf("%c is not considered alphanumeric", alphanumeric)
	}
}

// TestLetterIsNotNumeric calls the IsNumeric method for every letter in the
// alphabet and fails of the call returns true.
func TestLetterIsNotNumeric(test *testing.T) {
	for _, letter := range letters {
		if !Char(letter).IsNumeric() {
			continue
		}
		test.Errorf("letter (%c) is considered numeric", letter)
	}
}

func TestNumberIsNumeric(test *testing.T) {
	for _, number := range numbers {
		if Char(number).IsNumeric() {
			continue
		}
		test.Errorf("(%d) is not considered numeric", number)
	}
}

func TestNumberIsNotAlphabetic(test *testing.T) {
	for _, number := range numbers {
		if !Char(number).IsAlphabetic() {
			continue
		}
		test.Errorf("number mistaken for alphabetic char %c", number)
	}
}
