package document

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateSigner = tx.Transaction{
	Tag:         "updateSigner",
	Label:       "Update Signer",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:      "signer",
			Label:    "Signer",
			Required: true,
			DataType: "->user",
		},
		{
			Tag:      "updates",
			Label:    "Updates",
			Required: true,
			DataType: "@object",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		signerKey, ok := req["signer"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'signer' must be an asset")
		}

		signerAsset, err := signerKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get signer asset from the ledger")
		}

		signerMap := *signerAsset

		updates, ok := req["updates"].(map[string]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'updates' must be a map")
		}

		for key, value := range updates {
			signerMap[key] = value
		}

		updatedDocument, err := signerKey.Update(stub, signerMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update signer asset in the ledger")
		}

		updatedDocumentJSON, e := json.Marshal(updatedDocument)
		if e != nil {
			return nil, errors.WrapError(err, "Failed to marshal updated document asset")
		}

		return updatedDocumentJSON, nil
	},
}
