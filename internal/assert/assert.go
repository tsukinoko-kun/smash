package assert

import (
	"cmp"
	"fmt"
)

func NotEqual[T comparable](a, b T, errorMessage string) {
	if a == b {
		panic(fmt.Sprintf("assertion failed: NotEqual(%v, %v): %s", a, b, errorMessage))
	}
}

func Equal[T comparable](a, b T, errorMessage string) {
	if a != b {
		panic(fmt.Sprintf("assertion failed: Equal(%v, %v): %s", a, b, errorMessage))
	}
}

func NotEmpty(s string, errorMessage string) {
	if len(s) == 0 {
		panic(fmt.Sprintf("assertion failed: NotEmpty(%q): %s", s, errorMessage))
	}
}

func GreaterThan[T cmp.Ordered](x, greaterThan T, errorMessage string) {
	if x <= greaterThan {
		panic(fmt.Sprintf("assertion failed: GreaterThan(%v, %v): %s", x, greaterThan, errorMessage))
	}
}

func SmallerThan[T cmp.Ordered](x, smallerThan T, errorMessage string) {
	if x >= smallerThan {
		panic(fmt.Sprintf("assertion failed: SmallerThan(%v, %v): %s", x, smallerThan, errorMessage))
	}
}

func NotNil[T any](x *T, errorMessage string) {
	if x == nil {
		panic(fmt.Sprintf("assertion failed: NotNil(%v): %s", x, errorMessage))
	}
}
