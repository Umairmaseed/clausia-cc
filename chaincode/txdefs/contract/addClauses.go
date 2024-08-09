package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var AddClauses = tx.Transaction{
	Tag:         "addClauses",
	Label:       "Add Multiple Clauses",
	Description: "Adds multiple clauses to a contract by calling the AddClause transaction for each clause.",
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
			Tag:      "clauses",
			Label:    "Clauses",
			DataType: "[]@object",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		contractKey, ok := req["autoExecutableContract"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'autoExecutableContract' must be an asset key")
		}

		clauseList, ok := req["clauses"].([]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'clauses' must be a list of clauses")
		}

		for _, item := range clauseList {
			clauseMap, ok := item.(map[string]interface{})
			if !ok {
				return nil, errors.WrapError(nil, "Each clause must be an object")
			}

			args := map[string]interface{}{
				"autoExecutableContract": contractKey,
				"id":                     clauseMap["id"],
				"actionType":             clauseMap["actionType"],
			}

			if description, ok := clauseMap["description"].(string); ok {
				args["description"] = description
			}
			if category, ok := clauseMap["category"].(string); ok {
				args["category"] = category
			}
			if parameters, ok := clauseMap["parameters"].(map[string]interface{}); ok {
				args["parameters"] = parameters
			}
			if input, ok := clauseMap["input"].(map[string]interface{}); ok {
				args["input"] = input
			}
			if dependencies, ok := clauseMap["dependencies"].([]interface{}); ok {
				args["dependencies"] = dependencies
			}

			_, err := AddClause.Routine(stub, args)
			if err != nil {
				return nil, errors.WrapError(err, "Failed to add clause")
			}
		}

		// Get the final state of the contract
		contract, err := contractKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get updated contract asset from ledger")
		}

		// Return the end state of the contract
		responseJSON, nerr := json.Marshal(contract)
		if nerr != nil {
			return nil, errors.WrapError(err, "Failed to marshal response to JSON format")
		}

		return responseJSON, nil
	},
}
