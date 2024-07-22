package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var UploadDocument = tx.Transaction{
	Tag:         "uploadDocument",
	Label:       "Upload Document",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "originalHash",
			Label:    "originalHash",
			DataType: "sha256",
		},
		{
			Tag:      "finalHash",
			Label:    "finalHash",
			DataType: "sha256",
		},
		{
			Tag:      "status",
			Label:    "status",
			DataType: "statusType",
			Required: true,
		},
		{
			Tag:      "requiredSignatures",
			Label:    "requiredSignatures",
			DataType: "[]->signer",
			Required: true,
		},
		{
			Tag:      "successfulSignatures",
			Label:    "successfulSignatures",
			DataType: "[]->signer",
		},
		{
			Tag:      "rejectedSignatures",
			Label:    "rejectedSignatures",
			DataType: "[]->signer",
		},
		{
			Tag:      "originalDocURL",
			Label:    "originalDocURL",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "finalDocURL",
			Label:    "finalDocURL",
			DataType: "string",
		},
		{
			Tag:      "name",
			Label:    "name",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "owner",
			Label:    "Owner",
			DataType: "->signer",
			Required: true,
		},
		{
			Tag:      "timeout",
			Label:    "Timeout",
			DataType: "datetime",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		originalHash, ok := req["originalHash"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get originalHash parameter", 400)
		}

		status, ok := req["status"].(datatypes.StatusType)
		if !ok {
			return nil, errors.NewCCError("Failed to get status parameter", 400)
		}

		requiredSignatures, ok := req["requiredSignatures"].([]interface{})
		if !ok {
			return nil, errors.NewCCError("Failed to get requiredSignatures parameter", 400)
		}

		successfulSignatures, _ := req["successfulSignatures"].([]interface{})

		originalDocURL, ok := req["originalDocURL"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get originalDocURL parameter", 400)
		}

		name, ok := req["name"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get name parameter", 400)
		}
		owner, ok := req["owner"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("Failed to get owner parameter", 400)
		}
		timeout := req["timeout"]

		doc := map[string]interface{}{
			"@assetType":           "document",
			"originalHash":         originalHash,
			"status":               status,
			"requiredSignatures":   requiredSignatures,
			"successfulSignatures": successfulSignatures,
			"originalDocURL":       originalDocURL,
			"name":                 name,
			"owner":                owner,
			"timeout":              timeout,
		}
		if finalHash, ok := req["finalHash"].(string); ok && finalHash != "" {
			doc["finalHash"] = finalHash
		}

		if signature, ok := req["signature"].(assets.Key); ok {
			doc["signature"] = signature
		}

		if finalDocURL, ok := req["finalDocURL"].(string); ok && finalDocURL != "" {
			doc["finalDocURL"] = finalDocURL
		}

		if rejectedSignatures, ok := req["rejectedSignatures"].([]interface{}); ok {
			doc["rejectedSignatures"] = rejectedSignatures
		}

		document, err := assets.NewAsset(doc)
		if err != nil {
			return nil, errors.WrapError(err, "failed to create asset")
		}

		res, err := document.Put(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to write asset to the ledger")
		}

		resBytes, e := json.Marshal(res)
		if e != nil {
			return nil, errors.WrapError(e, "failed to marshal response")
		}

		return resBytes, nil
	},
}
