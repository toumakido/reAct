package data

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference between two integers
func Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}

// Divide returns the quotient of two integers
// Returns 0 if divisor is 0
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// Factorial returns the factorial of n
// Returns 1 for n <= 1
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// AbsoluteAdd returns the absolute value of the sum of two integers
func AbsoluteAdd(a, b int) int {
	return Abs(Add(a, b))
}

// ClampedMultiply multiplies two integers and clamps the result
func ClampedMultiply(a, b, min, max int) int {
	return Clamp(Multiply(a, b), min, max)
}

// SafeFactorial returns factorial using Max to ensure non-negative input
func SafeFactorial(n int) int {
	return Factorial(Max(n, 0))
}
