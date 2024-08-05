package utils

// Join two maps prioritizing map2 values over map1
func JoinMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	if map1 == nil {
		map1 = make(map[string]interface{})
	}

	for k, v := range map2 {
		map1[k] = v
	}

	return map1
}
