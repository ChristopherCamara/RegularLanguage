package stringArray

func IndexOf(test string, slice []string) int {
	for index, element := range slice {
		if test == element {
			return index
		}
	}
	return -1
}
