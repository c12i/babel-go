package main

import (
	"fmt"
)

func getDigits(num int) []int {
	// Handle the special case for zero
	if num == 0 {
		return []int{0}
	}

	var digits []int
	temp := num
	// Loop until the number becomes 0
	for temp != 0 {
		// Extract the last digit using modulo 10
		digit := temp % 10
		// Prepend the digit to the slice to maintain original order
		// Note: prepending to a slice can be less efficient for very large numbers
		digits = append([]int{digit}, digits...)
		// Truncate the last digit using integer division by 10
		temp /= 10
	}

	return digits
}

func main() {
	number := 12345
	digitsSlice := getDigits(number)
	fmt.Println("Original number:", number)
	fmt.Println("Digits:", digitsSlice)
	// Output:
	// Original number: 12345
	// Digits: [1 2 3 4 5]
}
