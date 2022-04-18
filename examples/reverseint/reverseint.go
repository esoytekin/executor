package reverseint

import "math"

func Reverse(x int32) int32 {
	isNegative := false

	if x < 0 {
		isNegative = true
		x = -x
	}

	var reverse int64 = 0

	for x > 0 {
		reverse = reverse*10 + int64(x%10)
		x /= 10
	}

	if reverse > math.MaxInt32 {
		return 0
	}

	if isNegative {
		return int32(-reverse)
	}

	return int32(reverse)
}
