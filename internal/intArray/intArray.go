package intArray

func IndexOf(test int, slice []int) int {
	for index, element := range slice {
		if test == element {
			return index
		}
	}
	return -1
}
