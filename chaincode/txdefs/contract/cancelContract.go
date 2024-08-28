package contract

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/models"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/params"
)

var CancelContract = tx.Transaction{
	Tag:         "cancelContract",
	Label:       "Cancel Contract",
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
			Tag:      "forceCancellation",
			Label:    "Force Cancellation",
			DataType: "boolean",
		},
		{
			Tag:      "requestedCancellation",
			Label:    "Requested Cancellation",
			DataType: "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		clauseKey, ok := req["clause"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'clause' must be an asset key")
		}

		clause, err := models.GetClause(stub, clauseKey)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause")
		}

		if clause.ActionType != datatypes.FinishContract {
			return nil, errors.NewCCError("action type must be of Finish Contract to cancel contract", 400)
		}

		if req["forceCancellation"] == nil && req["requestedCancellation"] == nil {
			return nil, errors.NewCCError("please provide condition to cancel the contract", 400)
		}

		if req["forceCancellation"] != nil && req["requestedCancellation"] != nil {
			return nil, errors.NewCCError("please provide one condition to cancel the contract", 400)
		}

		var parameters params.FinalizeContractParams
		if clause.Parameters != nil {
			bytes, jerr := json.Marshal(clause.Parameters)
			if jerr != nil {
				return nil, errors.WrapError(jerr, "Failed to marshal parameters")
			}

			jerr = json.Unmarshal(bytes, &parameters)
			if jerr != nil {
				return nil, errors.WrapError(jerr, "Failed to unmarshal to get parameters")
			}
		}

		if req["forceCancellation"] != nil {
			parameters.ForceCancellation = req["forceCancellation"].(bool)
		}

		if req["requestedCancellation"] != nil {
			parameters.RequestedCancellation = req["requestedCancellation"].(bool)
		}

		parametersJSON, nerr := json.Marshal(parameters)
		if nerr != nil {
			return nil, errors.WrapError(err, "Failed to marshal parameters to JSON")
		}

		clauseUpdate := map[string]interface{}{
			"parameters": parametersJSON,
		}

		updatedClauseAsset, err := clauseKey.Update(stub, clauseUpdate)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update clause asset in ledger")
		}

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "autoExecutableContract",
				"clauses": map[string]interface{}{
					"$elemMatch": map[string]interface{}{
						"@key": clause.Key,
					},
				},
			},
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for contracts with executable clauses", http.StatusInternalServerError)
		}

		if len(response.Result) > 0 {

			updateClauseAsset, err := models.GetClause(stub, clauseKey)
			if err != nil {
				return nil, errors.WrapError(err, "Failed to get clause")
			}

			contractKey := assets.Key{
				"@assetType": "autoExecutableContract",
				"@key":       response.Result[0]["@key"],
			}

			contractAsset, err := models.GetAutoExecutableContract(stub, contractKey)
			if err != nil {
				return nil, errors.WrapError(err, "Failed to get contract")
			}

			err = ExecuteClause(stub, contractAsset, updateClauseAsset)
			if err != nil {
				return nil, errors.WrapError(err, "Failed to execute clause")
			}
		} else {
			return nil, errors.NewCCError("clause is not associated with any contract", 400)
		}

		updatedClauseJSON, nerr := json.Marshal(updatedClauseAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "Failed to marshal updated clause asset")
		}

		return updatedClauseJSON, nil

	},
}
