package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var EditTemplate = tx.Transaction{
	Tag:         "editTemplate",
	Label:       "Edit Template",
	Description: "Edit the description, name, and public fields of a Template",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "template",
			Label:    "Template",
			DataType: "->template",
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
			Tag:      "public",
			Label:    "Public",
			DataType: "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		templateKey, ok := req["template"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'template' must be an asset key")
		}

		template, err := templateKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get template asset from ledger")
		}

		updateReq := map[string]interface{}{}

		if name, ok := req["name"].(string); ok {
			updateReq["name"] = name
		}

		if description, ok := req["description"].(string); ok {
			updateReq["description"] = description
		}

		if public, ok := req["public"].(bool); ok {
			updateReq["public"] = public
		}

		updatedTemplate, err := template.Update(stub, updateReq)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update template asset on the ledger")
		}

		responseJSON, nerr := json.Marshal(updatedTemplate)
		if nerr != nil {
			return nil, errors.WrapError(nerr, "Failed to marshal response to JSON format")
		}

		return responseJSON, nil
	},
}
