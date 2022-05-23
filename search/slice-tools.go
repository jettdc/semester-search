package search

import "math"

const ProximityRadius = 5

func removeFromSlice(slice []string, term string) []string {
	res := make([]string, 0)
	for _, e := range slice {
		if e != term {
			res = append(res, e)
		}
	}
	return res
}

func sliceContains(slice []string, term string) bool {
	for _, e := range slice {
		if e == term {
			return true
		}
	}
	return false
}

func WordsAreInProximity(searchTerms []string, searching []string, index int) bool {
	nowLookingFor := removeFromSlice(searchTerms, searching[index])

	// If we've found all elements we are looking for, then we are done searching and have found a solution
	if len(nowLookingFor) == 0 {
		return true
	}

	// See if any matches on the left side
	isInLeft := false
	leftRange := int(math.Max(float64(index-ProximityRadius), 0))
	leftFoundIndex := -1
	for i := leftRange; i < index; i++ {
		if sliceContains(nowLookingFor, searching[i]) {
			leftFoundIndex = i
			break
		}
	}

	// A search term was found in the left portion
	if leftFoundIndex != -1 {
		isInLeft = WordsAreInProximity(nowLookingFor, searching, leftFoundIndex)
	}

	// See if any matches on the right side
	isInRight := false
	rightRange := int(math.Min(float64(index+ProximityRadius), float64(len(searching)-1)))
	rightFoundIndex := -1
	for i := index + 1; i <= rightRange; i++ {
		if sliceContains(nowLookingFor, searching[i]) {
			rightFoundIndex = i
			break
		}
	}

	// A search term was found in the left portion
	if rightFoundIndex != -1 {
		isInRight = WordsAreInProximity(nowLookingFor, searching, rightFoundIndex)
	}

	return isInLeft || isInRight
}
