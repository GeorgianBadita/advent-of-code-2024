package datastructures

import (
	"testing"
)

func TestPow(t *testing.T) {
	r := Pow(2, 0)
	if r != 1 {
		t.Fatalf("Expected 2^0 to be 1, it was %d instead", r)
	}

	r = Pow(3, 5)
	if r != 243 {
		t.Fatalf("Expected 3^5 to be 243, it was %d instead", r)
	}

	r1 := Pow(3.0, -5)
	if Abs(r1-0.00411) > 1e-5 {
		t.Fatalf("Expected 3^-5 to be 0.00411, it was %f instead", r1)
	}
}
