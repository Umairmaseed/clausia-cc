package datatypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

var pemPubKey = assets.DataType{
	AcceptedFormats: []string{"string"},
	Description:     "Pem encoded public key",
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		pubKey, ok := data.(string)
		if !ok {
			return "", nil, errors.NewCCError("property must be a string", 400)
		}

		// Validates public key format
		// block, _ := pem.Decode([]byte(pubKey))
		// if block == nil || block.Type != "PUBLIC KEY" {
		// 	return "", nil, errors.NewCCError("The key format is not valid", 400)
		// }

		return pubKey, pubKey, nil
	},
}
