package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateAutoExecutableContract = tx.Transaction{
	Tag:         "createContract",
	Label:       "Create Auto executable contract",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "name",
			Label:    "Name",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "signatureDate",
			Label:    "Signature Date",
			DataType: "datetime",
		},
		{
			Tag:      "clauses",
			Label:    "Clauses",
			DataType: "[]->clause",
		},
		{
			Tag:      "data",
			Label:    "Data",
			DataType: "@object",
		},
		{
			Required: true,
			Tag:      "owner",
			Label:    "Owner",
			DataType: "->user",
		},
		{
			Tag:      "participants",
			Label:    "Participants",
			DataType: "[]->user",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		name, ok := req["name"].(string)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'name' must be a string")
		}
		signatureDate := req["signatureDate"]

		owner, ok := req["owner"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'owner' must be an asset key")
		}

		contract := map[string]interface{}{
			"@assetType":    "autoExecutableContract",
			"name":          name,
			"signatureDate": signatureDate,
			"owner":         owner,
			"dates":         map[string]interface{}{"signature": signatureDate},
		}

		if clauses, ok := req["clauses"].([]interface{}); ok {
			contract["clauses"] = clauses
		}
		if participants, ok := req["participants"].([]interface{}); ok {
			contract["participants"] = participants
		}
		if data, ok := req["data"].(map[string]interface{}); ok {
			contract["data"] = data
		}

		newContract, err := assets.NewAsset(contract)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create contract asset")
		}

		res, err := newContract.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to write contract asset to the ledger")
		}

		resBytes, e := json.Marshal(res)
		if e != nil {
			return nil, errors.WrapError(e, "Failed to marshal response")
		}

		return resBytes, nil
	},
}
