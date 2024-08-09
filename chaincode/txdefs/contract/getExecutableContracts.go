package contract

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ContractsWithExecutableClauses = tx.Transaction{
	Tag:         "contractsWithExecutableClauses",
	Label:       "Contracts with Executable Clauses",
	Description: "Retrieves all contracts containing clauses that are executable but not finalized",
	Method:      "GET",

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "autoExecutableContract",
			},
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for contracts with executable clauses", http.StatusInternalServerError)
		}

		var filteredContracts []map[string]interface{}

		for _, item := range response.Result {

			clauses, ok := item["clauses"].([]interface{})
			if !ok {
				continue
			}

			for _, clause := range clauses {
				clauseMap, ok := clause.(map[string]interface{})
				if !ok {
					continue
				}

				finalized, _ := clauseMap["finalized"].(bool)
				executable, _ := clauseMap["executable"].(bool)

				if !finalized && executable {
					filteredContracts = append(filteredContracts, item)
					break
				}
			}
		}

		responseJSON, nerr := json.Marshal(filteredContracts)
		if nerr != nil {
			return nil, errors.WrapErrorWithStatus(nerr, "error marshaling response", http.StatusInternalServerError)
		}

		return responseJSON, nil
	},
}
