package datastructures

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Signed | constraints.Float](v T) T {
	return max(-v, v)
}

func Pow[B constraints.Integer | constraints.Float, E constraints.Integer](b B, e E) B {
	if e < 0 {
		return pow(1/b, -e)
	}
	return pow(b, e)
}

func pow[B constraints.Integer | constraints.Float, E constraints.Integer](b B, e E) B {
	if e < 0 {
		panic(fmt.Sprintf("This function should only be called with positive integer exponent, found %d", e))
	}

	if e == 0 {
		return 1
	}
	if e == 1 {
		return b
	}

	halfPow := pow(b, e/2)
	if e%2 == 0 {
		return halfPow * halfPow
	}
	return b * halfPow * halfPow
}
