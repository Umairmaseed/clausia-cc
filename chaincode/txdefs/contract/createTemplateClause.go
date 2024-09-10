package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var CreateTemplateClause = tx.Transaction{
	Tag:         "createTemplateClause",
	Label:       "Create Template Clause",
	Description: "Transaction to create a new template clause",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "id",
			Label:    "Clause ID",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "template",
			Label:    "Template",
			DataType: "->template",
		},
		{
			Required: true,
			Tag:      "number",
			Label:    "Clause Number",
			DataType: "number",
		},
		{
			Required: true,
			Tag:      "name",
			Label:    "Clause Name",
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
			Tag:      "dependencies",
			Label:    "Dependencies",
			DataType: "[]->templateClause",
		},
		{
			Required: true,
			Tag:      "actionType",
			Label:    "Action Type",
			DataType: "actionType",
		},
		{
			Tag:      "defaultInputs",
			Label:    "Default Inputs",
			DataType: "@object",
		},
		{
			Tag:      "defaultParameters",
			Label:    "Default Parameters",
			DataType: "@object",
		},
		{
			Tag:      "optional",
			Label:    "Optional",
			DataType: "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		id, ok := req["id"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'id' must be a string")
		}

		template, ok := req["template"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'template' must be an asset key")
		}

		number, ok := req["number"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'number' must be a number")
		}

		name, ok := req["name"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'name' must be a string")
		}

		actionType, ok := req["actionType"].(datatypes.ActionType)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'actionType' must be of type action type")
		}

		templateClause := map[string]interface{}{
			"@assetType": "templateClause",
			"id":         id,
			"template":   template,
			"number":     number,
			"name":       name,
			"actionType": actionType,
		}

		if description, ok := req["description"].(string); ok {
			templateClause["description"] = description
		}

		if category, ok := req["category"].(string); ok {
			templateClause["category"] = category
		}

		if dependencies, ok := req["dependencies"].([]interface{}); ok {
			templateClause["dependencies"] = dependencies
		}

		if defaultInputs, ok := req["defaultInputs"].(map[string]interface{}); ok {
			templateClause["defaultInputs"] = defaultInputs
		}

		if defaultParameters, ok := req["defaultParameters"].(map[string]interface{}); ok {
			templateClause["defaultParameters"] = defaultParameters
		}
		if optional, ok := req["optional"].(bool); ok {
			templateClause["optional"] = optional
		}

		newTemplateClause, err := assets.NewAsset(templateClause)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create template clause asset")
		}

		res, err := newTemplateClause.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to write template clause asset to the ledger")
		}

		resBytes, e := json.Marshal(res)
		if e != nil {
			return nil, errors.WrapError(e, "Failed to marshal response")
		}

		return resBytes, nil
	},
}
