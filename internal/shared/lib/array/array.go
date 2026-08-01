package array

func ToSet[Element comparable](array []Element) map[Element]bool {
	set := make(map[Element]bool, len(array))
	for _, element := range array {
		set[element] = true
	}
	return set
}

func GetDistinctCount[Element comparable](array []Element) int {
	set := make(map[Element]bool, len(array))
	for _, element := range array {
		set[element] = true
	}
	return len(set)
}
