package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var RemoveTemplate = tx.Transaction{
	Tag:         "removeTemplate",
	Label:       "Remove Template",
	Description: "Remove or delete a template",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "template",
			Label:    "Template",
			DataType: "->template",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		templateKey, ok := req["template"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'template' must be an asset key")
		}

		templateAsset, err := templateKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get template asset from ledger")
		}

		_, err = templateKey.Delete(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to delete template asset from ledger")
		}

		response := map[string]interface{}{
			"message": "Template successfully deleted",
			"id":      (*templateAsset)["id"],
			"name":    (*templateAsset)["name"],
		}

		responseJSON, nerr := json.Marshal(response)
		if nerr != nil {
			return nil, errors.WrapError(nil, "Failed to marshal response")
		}

		return responseJSON, nil
	},
}
