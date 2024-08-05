package utils

// Function to generate a map of current clauses for quick lookup
func GenMapOfCurrClauses(currClauses []interface{}) map[string]interface{} {
	currCl := make(map[string]interface{})
	for i := range currClauses {
		clause, ok := currClauses[i].(map[string]interface{})
		if !ok {
			continue
		}
		key, ok := clause["@key"].(string)
		if !ok {
			continue
		}
		currCl[key] = currClauses[i]
	}
	return currCl
}
