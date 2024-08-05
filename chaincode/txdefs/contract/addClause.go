package contract

import (
	"encoding/json"
	"reflect"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/params"
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
		actionTypeFloat, ok := req["actionType"].(datatypes.ActionType)
		if !ok {
			return nil, errors.WrapError(nil, "Invalid type for actionType")
		}
		actionType := datatypes.ActionType(actionTypeFloat)

		paramHandler := params.Get(actionType)

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

		if actionType == datatypes.NonExecutable {
			clause["executable"] = false
		}

		if input, ok := req["input"].(map[string]interface{}); ok {
			formatOutputNames(input)
			inputType := paramHandler.GetInputs()
			filteredInput := filterFields(input, inputType)
			clause["input"] = filteredInput
		}

		if parameters, ok := req["parameters"].(map[string]interface{}); ok {
			formatOutputNames(parameters)
			paramsType := paramHandler.GetParameters()
			filteredParams := filterFields(parameters, paramsType)
			clause["parameters"] = filteredParams
		}

		if actionType == datatypes.NonExecutable {
			clause["executable"] = false
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

func formatOutputNames(params map[string]interface{}) {
	for key, value := range params {
		if strValue, ok := value.(string); ok {
			params[key] = utils.ValidateAndCleanData(strValue)
		}
	}
}

func filterFields(data map[string]interface{}, structType interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(structType)
	typ := reflect.TypeOf(structType)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Name
		}

		if value, exists := data[fieldName]; exists {
			result[fieldName] = value
		}
	}

	return result
}
