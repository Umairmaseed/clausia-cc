package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var EditTemplateClause = tx.Transaction{
	Tag:         "editTemplateClause",
	Label:       "Edit Template Clause",
	Description: "Edit the description, name, category, dependencies, actionType, defaultInputs, and defaultParameters fields of a TemplateClause",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "templateClause",
			Label:    "Template Clause",
			DataType: "->templateClause",
		},
		{
			Tag:      "name",
			Label:    "Name",
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
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		templateClauseKey, ok := req["templateClause"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'templateClause' must be an asset key")
		}

		templateClause, err := templateClauseKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get templateClause asset from ledger")
		}

		updateReq := map[string]interface{}{}

		if name, ok := req["name"].(string); ok {
			updateReq["name"] = name
		}

		if description, ok := req["description"].(string); ok {
			updateReq["description"] = description
		}

		if category, ok := req["category"].(string); ok {
			updateReq["category"] = category
		}

		if dependencies, ok := req["dependencies"].([]interface{}); ok {
			updateReq["dependencies"] = dependencies
		}

		if actionType, ok := req["actionType"].(string); ok {
			updateReq["actionType"] = actionType
		}

		if defaultInputs, ok := req["defaultInputs"].(map[string]interface{}); ok {
			updateReq["defaultInputs"] = defaultInputs
		}

		if defaultParameters, ok := req["defaultParameters"].(map[string]interface{}); ok {
			updateReq["defaultParameters"] = defaultParameters
		}

		updatedTemplateClause, err := templateClause.Update(stub, updateReq)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update templateClause asset on the ledger")
		}

		responseJSON, nerr := json.Marshal(updatedTemplateClause)
		if nerr != nil {
			return nil, errors.WrapError(nerr, "Failed to marshal response to JSON format")
		}

		return responseJSON, nil
	},
}
