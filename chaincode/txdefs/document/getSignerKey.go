package document

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetSignerKey = tx.Transaction{
	Tag:         "getSignerKey",
	Label:       "Get Signer Key",
	Description: "Retrieves the key of the signer by CPF",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:         "cpf",
			Description: "CPF of the signer",
			DataType:    "cpf",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		cpf := req["cpf"]

		fields := map[string]interface{}{
			"@assetType": "signer",
			"cpf":        cpf,
		}

		key, err := assets.NewKey(fields)
		if err != nil {
			return nil, errors.WrapError(err, "failed to generate key")
		}

		keyJSON, er := json.Marshal(key)
		if er != nil {
			return nil, errors.WrapError(err, "failed to marshal key to JSON")
		}

		return keyJSON, nil
	},
}
