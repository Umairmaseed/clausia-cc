package document

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateSigner = tx.Transaction{
	Tag:         "createSigner",
	Label:       "Create Signer",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:      "cpf",
			Label:    "Cpf",
			Required: true,
			DataType: "cpf",
		},
		{
			Tag:      "email",
			Label:    "Email",
			Required: true,
			DataType: "string",
		},
		{
			Tag:      "name",
			Label:    "Name",
			Required: true,
			DataType: "string",
		},
		{
			Tag:      "phone",
			Label:    "Phone",
			Required: true,
			DataType: "string",
		},
		{
			Tag:      "userName",
			Label:    "UserName",
			Required: true,
			DataType: "string",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		cpf, ok := req["cpf"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get cpf parameter", 400)
		}

		email, ok := req["email"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get email parameter", 400)
		}

		name, ok := req["name"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get name parameter", 400)
		}

		phone, ok := req["phone"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get phone parameter", 400)
		}

		userName, ok := req["userName"].(string)
		if !ok {
			return nil, errors.NewCCError("Failed to get userName parameter", 400)
		}

		signer := map[string]interface{}{
			"@assetType": "user",
			"cpf":        cpf,
			"email":      email,
			"name":       name,
			"phone":      phone,
			"userName":   userName,
		}

		newSigner, err := assets.NewAsset(signer)
		if err != nil {
			return nil, errors.WrapError(err, "failed to create user asset")
		}

		res, err := newSigner.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to write user asset to the ledger")
		}

		resBytes, e := json.Marshal(res)
		if e != nil {
			return nil, errors.WrapError(e, "failed to marshal response")
		}

		return resBytes, nil
	},
}
