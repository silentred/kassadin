package util

// InSliceInt checks if needle is in the stack
func InSliceInt(needle int, stack []int) bool {
	for _, item := range stack {
		if needle == item {
			return true
		}
	}
	return false
}
