package document

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetUserKey = tx.Transaction{
	Tag:         "getUserKey",
	Label:       "Get User Key",
	Description: "Retrieves the key of the User by CPF",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:         "cpf",
			Description: "CPF of the User",
			DataType:    "cpf",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		cpf := req["cpf"]

		fields := map[string]interface{}{
			"@assetType": "user",
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
