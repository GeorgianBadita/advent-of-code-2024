package datastructures

import "golang.org/x/exp/constraints"

func Abs[T constraints.Signed | constraints.Float](v T) T {
	return max(-v, v)
}
