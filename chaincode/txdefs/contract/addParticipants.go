package contract

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var AddParticipants = tx.Transaction{
	Tag:         "addParticipants",
	Label:       "Add Participants",
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
			Tag:      "participants",
			Label:    "Participants",
			DataType: "[]->user",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		contractKey, ok := req["autoExecutableContract"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'autoExecutableContract' must be an asset key")
		}

		participants, ok := req["participants"].([]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'participants' must be an slice")
		}

		contract, err := contractKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get autoExecutableContract asset from ledger")
		}

		updateReq := map[string]interface{}{
			"participants": participants,
		}

		updatedContract, err := contract.Update(stub, updateReq)

		responseJSON, nerr := json.Marshal(updatedContract)
		if nerr != nil {
			return nil, errors.WrapError(err, "Failed to marshal response to JSON format")
		}

		return responseJSON, nil
	},
}
