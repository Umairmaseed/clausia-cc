package utils

import (
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// Function to remove non-existing dependencies
func RemoveUnexisting(dependencies []interface{}, mapOfCurrClauses map[string]interface{}, stub *sw.StubWrapper) ([]interface{}, error) {
	newDependencies := []interface{}{}
	for _, depInterface := range dependencies {
		dep, ok := depInterface.(assets.Key)
		if !ok {
			continue
		}

		clause, err := dep.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause to add in dependencies")
		}

		key, _ := (*clause)["@key"].(string)

		found := false
		for _, item := range mapOfCurrClauses {
			currClause, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if currKey, ok := currClause["@key"].(string); ok && currKey == key {
				found = true
				break
			}
		}

		if found {
			newDependencies = append(newDependencies, depInterface)
		}
	}
	return newDependencies, nil
}
