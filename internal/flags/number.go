package flags

// Clamp limit min and max values
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Batches split num to slice of nums for each cc
func Batches(num, cc int) []int {
	if cc == 1 {
		return []int{num}
	}

	result := make([]int, cc)
	if num == 0 {
		return result
	}

	v := num / cc
	if v == 0 {
		result[0] = num
		return result
	}

	for i := 0; i < len(result); i++ {
		if i == len(result)-1 {
			result[i] = num
		} else {
			result[i] = v
		}
		num -= v
	}

	return result
}
