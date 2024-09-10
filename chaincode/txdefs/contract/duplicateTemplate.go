package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var DuplicateTemplate = tx.Transaction{
	Tag:         "duplicateTemplate",
	Label:       "Duplicate Template",
	Description: "Create a new template with all data from the original one, optionally with a new owner",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "originalTemplate",
			Label:    "Original Template",
			DataType: "->template",
		},
		{
			Required: true,
			Tag:      "id",
			Label:    "ID",
			DataType: "string",
		},
		{
			Tag:      "newOwner",
			Label:    "New Owner",
			DataType: "->user",
		},
		{
			Tag:      "name",
			Label:    "New Template Name",
			DataType: "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		originalTemplateKey, ok := req["originalTemplate"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'originalTemplate' must be an asset key")
		}

		originalTemplate, err := originalTemplateKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get original template asset from ledger")
		}

		newTemplateData := make(map[string]interface{})

		for key, value := range *originalTemplate {
			newTemplateData[key] = value
		}

		id, ok := req["id"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'id' must be an asset a string")

		}

		newTemplateData["@assetType"] = "template"
		newTemplateData["id"] = id

		if newOwner, ok := req["newOwner"].(assets.Key); ok {
			newTemplateData["owner"] = newOwner
		}

		if newName, ok := req["name"].(string); ok {
			newTemplateData["name"] = newName
		}

		newTemplate, err := assets.NewAsset(newTemplateData)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create new template asset")
		}

		res, err := newTemplate.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to write new template asset to the ledger")
		}

		resBytes, nerr := json.Marshal(res)
		if nerr != nil {
			return nil, errors.WrapError(err, "Failed to marshal response")
		}

		return resBytes, nil
	},
}
