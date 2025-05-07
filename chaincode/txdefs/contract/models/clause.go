package models

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger-labs/clausia-cc/chaincode/datatypes"
)

type Clause struct {
	Key          string                 `json:"@key"`
	Id           string                 `json:"id"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	Parameters   map[string]interface{} `json:"parameters"`
	Input        map[string]interface{} `json:"input"`
	Executable   bool                   `json:"executable"`
	Dependencies []assets.Key           `json:"dependencies"`
	ActionType   datatypes.ActionType   `json:"actionType"`
	Finalized    bool                   `json:"finalized"`
	Result       map[string]interface{} `json:"result"`
	Asset        *assets.Asset
}

func GetClause(stub *sw.StubWrapper, key assets.Key) (*Clause, errors.ICCError) {
	asset, err := key.Get(stub)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get clause asset")
	}

	bytes, jerr := json.Marshal(asset)
	if jerr != nil {
		return nil, errors.WrapError(jerr, "Failed to marshal clause asset")
	}

	var clause Clause
	jerr = json.Unmarshal(bytes, &clause)
	if jerr != nil {
		return nil, errors.WrapError(jerr, "Failed to unmarshal to get clause")
	}

	clause.Asset = asset
	return &clause, nil
}
