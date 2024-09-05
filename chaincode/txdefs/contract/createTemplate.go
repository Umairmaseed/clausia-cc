package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateTemplate = tx.Transaction{
	Tag:         "createTemplate",
	Label:       "Create Template",
	Description: "Transaction to create a new template",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "id",
			Label:    "Template ID",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "name",
			Label:    "Template Name",
			DataType: "string",
		},
		{
			Tag:      "description",
			Label:    "Description",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "creator",
			Label:    "Creator",
			DataType: "->user",
		},
		{
			Required: true,
			Tag:      "public",
			Label:    "Public",
			DataType: "boolean",
		},
		{
			Tag:      "clauses",
			Label:    "Clauses",
			DataType: "[]->templateClause",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		id, ok := req["id"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'id' must be a string")
		}

		name, ok := req["name"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'name' must be a string")
		}

		creator, ok := req["creator"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'creator' must be an asset key")
		}

		public, ok := req["public"].(bool)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'public' must be boolean value")
		}

		template := map[string]interface{}{
			"@assetType": "template",
			"id":         id,
			"name":       name,
			"creator":    creator,
			"public":     public,
		}

		if description, ok := req["description"].(string); ok {
			template["description"] = description
		}

		if clauses, ok := req["clauses"].([]interface{}); ok {
			template["clauses"] = clauses
		}

		newTemplate, err := assets.NewAsset(template)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create template asset")
		}

		res, err := newTemplate.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to write template asset to the ledger")
		}

		resBytes, e := json.Marshal(res)
		if e != nil {
			return nil, errors.WrapError(e, "Failed to marshal response")
		}

		return resBytes, nil
	},
}
