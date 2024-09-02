package datatypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

// CustomDataTypes contain the user-defined primary data types
var CustomDataTypes = map[string]assets.DataType{
	"sha256":     Sha256,
	"statusType": statusType,
	"pemPubKey":  pemPubKey,
	"cpf":        cpf,
	"actionType": actionType,
	"argDt":      argDt,
}
