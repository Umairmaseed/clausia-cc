package document

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/clausia-cc/chaincode/datatypes"
)

var PutSignature = tx.Transaction{
	Tag:         "putSignature",
	Label:       "Put Signature",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:      "document",
			Label:    "Document",
			Required: true,
			DataType: "@asset",
		},
		{
			Tag:      "user",
			Label:    "User",
			Required: true,
			DataType: "@asset",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		documentAsset, ok := req["document"].(assets.Asset)
		if !ok {
			return nil, errors.NewCCError("Failed to get document parameter", 400)
		}

		signerAsset, ok := req["user"].(assets.Asset)
		if !ok {
			return nil, errors.NewCCError("Failed to get user parameter", 400)
		}

		documentKey, err := assets.NewKey(documentAsset)
		if err != nil {
			return nil, errors.WrapError(err, "Invalid document key")
		}

		response := make(map[string]interface{})

		exists, err := documentKey.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to check if document exists in ledger")
		}

		signerAvailable, err := signerAsset.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to check if signer exists in ledger")
		}

		if !signerAvailable {
			return nil, errors.NewCCError("Signer is not registered in blockchain", 400)
		}

		if exists && signerAvailable {
			document, err := documentKey.Get(stub)
			if err != nil {
				return nil, errors.WrapError(err, "Failed to get document from the ledger")
			}

			requiredSignatures, ok := document.GetProp("requiredSignatures").([]interface{})
			if !ok {
				return nil, errors.NewCCError("Failed to get requiredSignatures property", 400)
			}

			successfulSignatures, _ := document.GetProp("successfulSignatures").([]interface{})
			if successfulSignatures == nil {
				successfulSignatures = []interface{}{}
			}

			rejectedSignatures, _ := document.GetProp("rejectedSignatures").([]interface{})
			if rejectedSignatures == nil {
				rejectedSignatures = []interface{}{}
			}

			status := document.GetProp("status").(datatypes.StatusType)
			if status != 0 {
				return nil, errors.NewCCError("Document is not waiting for signatures", 400)
			}

			signerKey := map[string]interface{}{
				"@assetType": signerAsset["@assetType"],
				"@key":       signerAsset["@key"],
			}

			isSignerRequired := false
			for _, sig := range requiredSignatures {
				if sigKey, ok := sig.(map[string]interface{}); ok {
					if sigKey["@key"] == signerKey["@key"] {
						isSignerRequired = true
						break
					}
				}
			}

			if !isSignerRequired {
				return nil, errors.NewCCError("Signer is not required for this document", 400)
			}

			isSignerSuccessful := false
			for _, sig := range successfulSignatures {
				if sigKey, ok := sig.(map[string]interface{}); ok {
					if sigKey["@key"] == signerKey["@key"] {
						isSignerSuccessful = true
						break
					}
				}
			}

			if !isSignerSuccessful {
				successfulSignatures = append(successfulSignatures, signerKey)
			}

			isLastSignature := len(successfulSignatures) == len(requiredSignatures)

			fields := map[string]interface{}{
				"successfulSignatures": successfulSignatures,
			}

			if len(rejectedSignatures) > 0 {
				fields["rejectedSignatures"] = rejectedSignatures
			}

			if isLastSignature {
				fields["status"] = 3
			}

			updatedDocument, err := documentKey.Update(stub, fields)
			if err != nil {
				return nil, errors.WrapError(err, "Failed to update document")
			}

			response = map[string]interface{}{
				"document": updatedDocument,
			}

		} else {
			documentAsset["successfulSignatures"] = []interface{}{signerAsset}
			documentAsset["status"] = 0

			newDocument, err := documentAsset.PutNew(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to write asset to the ledger")
			}

			response["document"] = newDocument
		}

		resBytes, e := json.Marshal(response)
		if e != nil {
			return nil, errors.WrapError(e, "failed to marshal response")
		}

		return resBytes, nil
	},
}
