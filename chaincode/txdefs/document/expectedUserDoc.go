package document

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var ExpectedUserDoc = tx.Transaction{
	Tag:         "expectedUserDoc",
	Label:       "Expected User Document",
	Description: "",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:      "signer",
			Label:    "Signer",
			Required: true,
			DataType: "->signer",
		},
		{
			Tag:      "status",
			Label:    "Status",
			DataType: "statusType",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		signerObj, signerOk := req["signer"].(assets.Key)
		if !signerOk {
			return nil, errors.NewCCError("Invalid document key", http.StatusBadRequest)
		}

		signerKey, keyOk := signerObj["@key"].(string)
		if !keyOk {
			return nil, errors.NewCCError("Invalid signer key", http.StatusBadRequest)
		}

		status, statusOk := req["status"].(datatypes.StatusType)

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "document",
				"requiredSignatures": map[string]interface{}{
					"$elemMatch": map[string]interface{}{
						"$eq": map[string]interface{}{
							"@assetType": "signer",
							"@key":       signerKey,
						},
					},
				},
			},
		}

		if statusOk {
			query["selector"].(map[string]interface{})["status"] = status
		}

		// Prepare the response
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for document", http.StatusInternalServerError)
		}
		responseJSON, errr := json.Marshal(response)
		if errr != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", http.StatusInternalServerError)
		}

		return responseJSON, nil
	},
}
