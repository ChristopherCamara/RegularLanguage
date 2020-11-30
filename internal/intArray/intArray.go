package intArray

import "fmt"

func IndexOf(find int, slice []int) int {
	for index, element := range slice {
		if find == element {
			return index
		}
	}
	return -1
}

func Remove(find int, slice *[]int) {
	for index, element := range *slice {
		if find == element {
			if index == len(*slice)-1 {
				*slice = (*slice)[:index]
			} else {
				*slice = append((*slice)[:index], (*slice)[index+1:]...)
			}
		}
	}
}

func Print(slice []int) {
	for index, element := range slice {
		if index == len(slice)-1 {
			fmt.Println(element)
		} else {
			fmt.Printf("%d, ", element)
		}
	}
}
