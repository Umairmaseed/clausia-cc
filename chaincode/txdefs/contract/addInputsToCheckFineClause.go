package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var AddInputToCheckFineClause = tx.Transaction{
	Tag:         "addInputToCheckFineClause",
	Label:       "Add Input To Check Fine Clause",
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
			Tag:      "referenceValue",
			Label:    "Reference Value",
			DataType: "number",
		},
		{
			Tag:      "dailyPercentage",
			Label:    "Daily Percentage",
			DataType: "number",
		},
		{
			Tag:      "days",
			Label:    "days",
			DataType: "number",
		},
		{
			Tag:      "referenceClauseDays",
			Label:    "Reference Clause Days",
			DataType: "boolean",
		},
		{
			Tag:      "referenceClauseName",
			Label:    "Reference Clause Name",
			DataType: "string",
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

		if actionType == datatypes.GetDeduction {
			input, ok := (*clauseAsset)["input"].(map[string]interface{})
			if !ok {
				input = make(map[string]interface{})
				(*clauseAsset)["input"] = input
			}

			if referenceValue, ok := req["referenceValue"].(float64); ok {
				input["referenceValue"] = referenceValue
			}
			if dailyPercentage, ok := req["dailyPercentage"].(float64); ok {
				input["dailyPercentage"] = dailyPercentage
			}
			if days, ok := req["days"].(float64); ok {
				input["days"] = days
			}
			if referenceClauseDays, ok := req["referenceClauseDays"].(bool); ok {
				input["referenceClauseDays"] = referenceClauseDays
			}
			if referenceClauseName, ok := req["referenceClauseName"].(string); ok {
				input["referenceClauseName"] = referenceClauseName
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

		return nil, errors.NewCCError("Action type is not GEtDeduction", 400)
	},
}
