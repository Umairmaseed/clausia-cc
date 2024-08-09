package document

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateDocument = tx.Transaction{
	Tag:         "updateDocument",
	Label:       "Update Document",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:      "document",
			Label:    "Document",
			Required: true,
			DataType: "->document",
		},
		{
			Tag:      "updates",
			Label:    "Updates",
			Required: true,
			DataType: "@object",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		documentKey, ok := req["document"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'document' must be an asset")
		}

		documentAsset, err := documentKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get document asset from the ledger")
		}

		documentMap := *documentAsset

		updates, ok := req["updates"].(map[string]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'updates' must be a map")
		}

		for key, value := range updates {
			documentMap[key] = value
		}

		updatedDocument, err := documentKey.Update(stub, documentMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update document asset in the ledger")
		}

		updatedDocumentJSON, e := json.Marshal(updatedDocument)
		if e != nil {
			return nil, errors.WrapError(err, "Failed to marshal updated document asset")
		}

		return updatedDocumentJSON, nil
	},
}
