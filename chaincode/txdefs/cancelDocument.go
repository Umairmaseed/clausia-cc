package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var CancelDocument = tx.Transaction{
	Tag:         "cancelDocument",
	Label:       "Cancel Document",
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
			Tag:      "status",
			Label:    "Status",
			Required: true,
			DataType: "statusType",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		documentKey, ok := req["document"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'document' must be an asset")
		}

		status, ok := req["status"].(datatypes.StatusType)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'status' must be of type 'statusType'")
		}

		// Retrieve the document asset from the ledger
		documentAsset, err := documentKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get document asset from the ledger")
		}

		documentMap := *documentAsset

		currentStatus, ok := documentMap["status"].(datatypes.StatusType)
		if ok && ((currentStatus == 1 && status == 1) || (currentStatus == 2 && status == 2)) {
			var statusMessage string
			switch currentStatus {
			case 1:
				statusMessage = "Document status already set to cancelled"
			case 2:
				statusMessage = "Document status already set to expired"
			}
			return nil, errors.NewCCError(statusMessage, 400)
		}

		if status != 1 && status != 2 {
			return nil, errors.NewCCError("Transaction is only used to set document status to cancel or expired", 400)
		}

		// Update the status of the document
		documentMap["status"] = status

		// Update the document asset in the ledger
		updatedDocument, err := documentKey.Update(stub, documentMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update document asset in the ledger")
		}

		// Marshal the updated document asset to JSON format
		updatedDocumentJSON, e := json.Marshal(updatedDocument)
		if e != nil {
			return nil, errors.WrapError(err, "Failed to marshal updated document asset")
		}

		return updatedDocumentJSON, nil
	},
}
