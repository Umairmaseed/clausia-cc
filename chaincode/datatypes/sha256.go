package datatypes

import (
	"regexp"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

var sha256 = assets.DataType{
	AcceptedFormats: []string{"string"},
	Description:     "Sha256 hash digest string (hexadecimal)",
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		hash, ok := data.(string)
		if !ok {
			return "", nil, errors.NewCCError("property must be a string", 400)
		}

		// Check if the string is in hexadecimal
		isHex, err := regexp.MatchString(`^[a-fA-F0-9]+$`, hash)
		if err != nil {
			return "", nil, errors.WrapError(err, "failed to match hexadecimal regular expression")
		}
		if !isHex {
			return "", nil, errors.NewCCError("The hash digest must be a hexadecimal string", 400)
		}

		// If the hash has 64 hex chars, than it has 256 bits
		has256Bits := len(hash) == 64
		if !has256Bits {
			return "", nil, errors.NewCCError("The hash digest must have 64 characters", 400)
		}

		return hash, hash, nil
	},
}
