package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/clausia-cc/chaincode/txdefs/contract/params"
)

var AddReviewToContract = tx.Transaction{
	Tag:         "addReviewToContract",
	Label:       "Add review to contract",
	Description: "",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "autoExecutableContract",
			Label:    "Auto Executable Contract",
			DataType: "->autoExecutableContract",
		},
		{
			Tag:      "review",
			Label:    "Review",
			DataType: "@object",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		contractKey, ok := req["autoExecutableContract"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'autoExecutableContract' must be an asset key")
		}

		reviewMap, ok := req["review"].(map[string]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'review' must be an object")
		}

		reviewBytes, err := json.Marshal(reviewMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to marshal review map")
		}

		var review params.Review
		err = json.Unmarshal(reviewBytes, &review)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to unmarshal review data into Review struct")
		}

		contract, err := contractKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get autoExecutableContract asset from ledger")
		}

		data, ok := (*contract)["data"].(map[string]interface{})
		if !ok {
			data = make(map[string]interface{})
		}

		if _, exists := data["review"].(map[string]interface{}); exists {
			return nil, errors.WrapError(nil, "Contract already contains a review; cannot add a new review")
		}

		data["review"] = review

		updateReq := map[string]interface{}{
			"data": data,
		}

		updatedContract, err := contract.Update(stub, updateReq)

		responseJSON, nerr := json.Marshal(updatedContract)
		if nerr != nil {
			return nil, errors.WrapError(err, "Failed to marshal response to JSON format")
		}

		return responseJSON, nil
	},
}
