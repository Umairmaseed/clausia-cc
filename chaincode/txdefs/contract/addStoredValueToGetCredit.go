package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var AddStoredValueToGetCredit = tx.Transaction{
	Tag:         "addStoredValueToGetCredit",
	Label:       "Add Stored Value to Get Credit",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "clause",
			Label:    "Clause",
			DataType: "->clause",
		},
		{
			Tag:      "storedValue",
			Label:    "Stored Value",
			DataType: "number",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		clauseKey, ok := req["clause"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("Invalid clause format", 400)
		}

		clauseAsset, err := clauseKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause asset from ledger")
		}

		actionType, ok := (*clauseAsset)["actionType"].(datatypes.ActionType)
		if !ok {
			return nil, errors.NewCCError("Invalid action type format", 400)
		}

		if actionType == datatypes.GetCredit {
			input, ok := (*clauseAsset)["input"].(map[string]interface{})
			if !ok {
				input = make(map[string]interface{})
				(*clauseAsset)["input"] = input
			}

			if storedValue, ok := req["storedValue"].(float64); ok {
				input["storedValue"] = storedValue
			}

			clauseUpdated, err := clauseAsset.Update(stub, map[string]interface{}{
				"input": input,
			})
			if err != nil {
				return nil, errors.WrapError(err, "Failed to update clause")
			}

			response, nerr := json.Marshal(clauseUpdated)
			if nerr != nil {
				return nil, errors.WrapError(err, "Failed to marshal updated clause")
			}

			return response, nil
		}

		return nil, errors.NewCCError("Action type is not GetCredit", 400)
	},
}
