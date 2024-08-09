package utils

import (
	"regexp"
	"strings"
)

// ValidateAndCleanData validates and cleans data by removing special characters and spaces
func ValidateAndCleanData(data map[string]interface{}) map[string]interface{} {
	// regex to remove special characters (allowing only letters, numbers, and underscores)
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)

	// Clean the data
	cleanedData := make(map[string]interface{})
	for key, value := range data {
		// Remove spaces and special characters from the key
		cleanedKey := re.ReplaceAllString(strings.ReplaceAll(key, " ", "_"), "")

		// Handle the value if it's a nested map (recursive cleaning)
		if nestedMap, ok := value.(map[string]interface{}); ok {
			cleanedData[cleanedKey] = ValidateAndCleanData(nestedMap)
		} else {
			cleanedData[cleanedKey] = value
		}
	}

	return cleanedData
}
