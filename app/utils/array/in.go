package array

func In(arr []string, val string) bool {
	rel := map[string]int{}
	for _, s := range arr {
		rel[s] = 1
	}

	_, ok := rel[val]
	return ok
}
