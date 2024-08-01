package contract

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var AddReferenceDateCDI = tx.Transaction{
	Tag:         "addReferenceDateCDI",
	Label:       "Add Reference Date to CDI Clause",
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
			Required: true,
			Tag:      "referenceDate",
			Label:    "Reference Data",
			DataType: "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		clauseKey, ok := req["clause"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("Invalid clause format", 400)
		}

		referenceDate, ok := req["referenceDate"].(string)
		if !ok {
			return nil, errors.NewCCError("Invalid reference date format", 400)
		}

		_, err := time.Parse(time.RFC3339, referenceDate)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to parse reference date")
		}

		clauseAsset, err := clauseKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause asset from ledger")
		}

		actionType, ok := (*clauseAsset)["actionType"].(datatypes.ActionType)
		if !ok {
			return nil, errors.NewCCError("Invalid action type format", 400)
		}

		if actionType == datatypes.CheckDateInterval {
			input, ok := (*clauseAsset)["input"].(map[string]interface{})
			if !ok {
				input = make(map[string]interface{})
				(*clauseAsset)["input"] = input
			}

			input["referenceDate"] = referenceDate

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

		return nil, errors.NewCCError("Action type is not CheckDateInterval", 400)
	},
}
