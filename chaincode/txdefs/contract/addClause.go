package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/utils"
)

var AddClause = tx.Transaction{
	Tag:         "addClause",
	Label:       "Add Clause",
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
			Tag:      "id",
			Label:    "Id",
			DataType: "string",
		},
		{
			Tag:      "description",
			Label:    "Description",
			DataType: "string",
		},
		{
			Tag:      "category",
			Label:    "Category",
			DataType: "string",
		},
		{
			Tag:      "parameters",
			Label:    "Parameters",
			DataType: "@object",
		},
		{
			Tag:      "input",
			Label:    "Input",
			DataType: "@object",
		},
		{
			Tag:      "dependencies",
			Label:    "Dependencies",
			DataType: "[]->clause",
		},
		{
			Required: true,
			Tag:      "actionType",
			Label:    "Action Type",
			DataType: "actionType",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		actionType, _ := req["actionType"].(datatypes.ActionType)

		clause := map[string]interface{}{
			"@assetType": "clause",
			"id":         req["id"],
			"actionType": actionType,
			"executable": true,
			"finalized":  false,
		}

		if description, ok := req["description"].(string); ok {
			clause["description"] = description
		}
		if category, ok := req["category"].(string); ok {
			clause["category"] = category
		}
		if input, ok := req["input"].(map[string]interface{}); ok {
			clause["input"] = utils.ValidateAndCleanData(input)
		}

		if actionType == datatypes.NonExecutable {
			clause["executable"] = false
		}

		if parameters, ok := req["parameters"].(map[string]interface{}); ok {
			clause["parameters"] = utils.ValidateAndCleanData(parameters)
		}

		contractKey, ok := req["autoExecutableContract"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'autoExecutableContract' must be an asset key")
		}

		contract, err := contractKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get autoExecutableContract asset from ledger")
		}

		clauses, exists := (*contract)["clauses"].([]interface{})
		if !exists {
			clauses = make([]interface{}, 0)
		}
		mapOfCurrClauses := utils.GenMapOfCurrClauses(clauses)

		if dependencies, ok := req["dependencies"].([]interface{}); ok {
			removedUnexistingClause, err := utils.RemoveUnexisting(dependencies, mapOfCurrClauses, stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to removed unexisting clauses")
			}
			clause["dependencies"] = removedUnexistingClause
		}

		newClause, err := assets.NewAsset(clause)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create clause asset")
		}

		clauseAsset, err := newClause.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to save clause asset on ledger")
		}

		clauses = append(clauses, clauseAsset)

		updatedContract, err := contractKey.Update(stub, map[string]interface{}{
			"clauses": clauses,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update contract asset with new clause")
		}

		responseJSON, nerr := json.Marshal(updatedContract)
		if nerr != nil {
			return nil, errors.WrapError(err, "Failed to marshal response to JSON format")
		}

		return responseJSON, nil
	},
}
