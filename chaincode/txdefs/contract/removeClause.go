package contract

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var RemoveClause = tx.Transaction{
	Tag:         "removeClause",
	Label:       "Remove Clause",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "autoExecutableContract",
			Label:    "Auto Executable Contract",
			DataType: "->autoExecutableContract",
		},
		{
			Required: true,
			Tag:      "clause",
			Label:    "Clause",
			DataType: "->clause",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		contractKey, ok := req["autoExecutableContract"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'contract' must be an asset key")
		}

		contractAsset, err := contractKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get autoExecutableContract asset from ledger")
		}

		clauseKey, ok := req["clause"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'clause' must be an asset key")
		}

		clauseAsset, err := clauseKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause asset from ledger")
		}

		clauses, ok := (*contractAsset)["clauses"].([]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Clauses field is not an array")
		}

		var updatedClauses []interface{}
		var clauseFound bool
		for _, c := range clauses {
			clauseMap, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			if clauseMap["@key"] != (*clauseAsset)["@key"] {
				updatedClauses = append(updatedClauses, c)
			} else {
				clauseFound = true
			}
		}

		if !clauseFound {
			return nil, errors.NewCCError("Clause does not belong to contract", http.StatusBadRequest)
		}

		if len(updatedClauses) == 0 {
			updatedClauses = []interface{}{}
		}

		contractUpdates := map[string]interface{}{
			"clauses": updatedClauses,
		}

		updatedContractAsset, err := contractKey.Update(stub, contractUpdates)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update autoExecutableContract asset in ledger")
		}

		_, err = clauseKey.Delete(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to delete clause")
		}

		updatedContractJSON, nerr := json.Marshal(updatedContractAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "Failed to marshal updated contract asset")
		}

		return updatedContractJSON, nil

	},
}
