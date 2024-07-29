package utils

// Function to generate a map of current clauses for quick lookup
func GenMapOfCurrClauses(currClauses []interface{}) map[string]interface{} {
	currCl := map[string]interface{}{}
	for i := range currClauses {
		clause, ok := currClauses[i].(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := clause["id"].(string)
		currCl[key] = currClauses[i]
	}
	return currCl
}
