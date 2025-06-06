package contract

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/clausia-cc/chaincode/datatypes"
)

var AddInputsToMakePaymentClause = tx.Transaction{
	Tag:         "addInputsToMakePaymentClause",
	Label:       "Add Input To make payment Clause",
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
			Tag:      "date",
			Label:    "Date",
			DataType: "datetime",
			Required: true,
		},
		{
			Tag:      "payment",
			Label:    "Payment",
			DataType: "number",
			Required: true,
		},
		{
			Tag:      "receiptHash",
			Label:    "Receipt hash",
			DataType: "string",
		},
		{
			Tag:      "finalPayment",
			Label:    "FinalPayment",
			DataType: "boolean",
			Required: true,
		},
		{
			Tag:      "receiptUrl",
			Label:    "receipt url",
			DataType: "string",
		},
		{
			Tag:      "stripeToken",
			Label:    "Stripe Token",
			DataType: "string",
		},
		{
			Tag:      "payPalTransactionID",
			Label:    "PayPal Transaction ID",
			DataType: "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		clauseKey, ok := req["clause"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("Invalid clause format", 400)
		}

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "autoExecutableContract",
				"clauses": map[string]interface{}{
					"$elemMatch": clauseKey,
				},
			},
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for query", http.StatusInternalServerError)
		}

		contractKey := assets.Key{}
		if len(response.Result) != 0 {
			contractKey = assets.Key{
				"@assetType": "autoExecutableContract",
				"@key":       response.Result[0]["@key"],
			}
		}

		contractAsset, err := contractKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get autoExecutableContract asset from ledger")
		}

		clauseAsset, err := clauseKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause asset from ledger")
		}

		actionType, ok := (*clauseAsset)["actionType"].(datatypes.ActionType)
		if !ok {
			return nil, errors.NewCCError("Invalid action type format", 400)
		}

		paymentName, ok := (*clauseAsset)["parameters"].(map[string]interface{})["name"].(string)
		if !ok {
			return nil, errors.NewCCError("Invalid payment name format", 400)
		}

		if actionType == datatypes.Payment {
			input, ok := (*clauseAsset)["input"].(map[string]interface{})
			if !ok {
				input = make(map[string]interface{})
				(*clauseAsset)["input"] = input
			}

			input["date"] = req["date"]
			input["payment"] = req["payment"]
			input["finalPayment"] = req["finalPayment"]

			if req["receiptHash"] != nil {
				hash, _, err := datatypes.Sha256.Parse(req["receiptHash"])
				if err != nil {
					return nil, errors.WrapError(err, "Failed to update clause")
				}

				input["receiptHash"] = hash
			}

			if req["receiptUrl"] != nil {
				input["receiptUrl"] = req["receiptUrl"]
			}

			if req["stripeToken"] != nil {
				input["stripeToken"] = req["stripeToken"]
			}
			if req["payPalTransactionID"] != nil {
				input["payPalTransactionID"] = req["payPalTransactionID"]
			}

			clauseUpdated, err := clauseAsset.Update(stub, map[string]interface{}{
				"input": input,
			})
			if err != nil {
				return nil, errors.WrapError(err, "Failed to update clause")
			}

			contractDates, ok := (*contractAsset)["dates"].(map[string]interface{})
			if !ok {
				contractDates = make(map[string]interface{})
			}

			contractDates[paymentName] = req["date"]

			_, err = contractAsset.Update(stub, map[string]interface{}{
				"dates": contractDates,
			})

			if err != nil {
				return nil, errors.WrapError(err, "Failed to update contract asset")
			}

			response, nerr := json.Marshal(clauseUpdated)
			if nerr != nil {
				return nil, errors.WrapError(err, "Failed to marshal updated clause")
			}

			return response, nil
		}

		return nil, errors.NewCCError("Action type is not payment", 400)
	},
}
