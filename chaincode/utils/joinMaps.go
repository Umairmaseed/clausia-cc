package utils

// Join three maps prioritizing map3 over map2, and map2 over map1.
func JoinMaps(map1, map2, map3 map[string]interface{}) map[string]interface{} {
	if map1 == nil {
		map1 = make(map[string]interface{})
	}

	// Merge map2 into map1, prioritizing map2 values
	for k, v := range map2 {
		map1[k] = v
	}

	// Merge map3 into the result, prioritizing map3 values
	for k, v := range map3 {
		map1[k] = v
	}

	return map1
}
