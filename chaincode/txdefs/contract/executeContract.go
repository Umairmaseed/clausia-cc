package contract

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/clausia-cc/chaincode/txdefs/contract/models"
	"github.com/hyperledger-labs/clausia-cc/chaincode/txdefs/contract/params"
	"github.com/hyperledger-labs/clausia-cc/chaincode/utils"
)

var ExecuteAutoExecutableContract = tx.Transaction{
	Tag:         "executeAutoExecutableContract",
	Label:       "Execute Auto Executable Contract",
	Description: "Executes an auto-executable contract",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "contract",
			Label:       "Contract",
			Description: "Auto Executable Contract",
			DataType:    "->autoExecutableContract",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		contractKey, _ := req["contract"].(assets.Key)

		contract, err := models.GetAutoExecutableContract(stub, contractKey)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get auto executable contract")
		}

		for _, clause := range contract.Clauses {
			err = ExecuteClause(stub, contract, clause)
			if err != nil {
				fmt.Printf("Error executing clause %s: %v\n", clause.Key, err)
				continue
			}
		}

		updatedContract, err := contract.Asset.Update(stub, map[string]interface{}{
			"data": contract.Data,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update contract")
		}

		responseJSON, nerr := json.Marshal(updatedContract)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode response to JSON format")
		}

		return responseJSON, nil
	},
}

func ExecuteClause(stub *sw.StubWrapper, contract *models.AutoExecutableContract, clause *models.Clause) errors.ICCError {
	if clause.Finalized || !clause.Executable {
		return nil
	}

	for _, dep := range clause.Dependencies {
		depClause := contract.GetClause(dep.Key())
		if depClause == nil {
			return errors.NewCCError("Clause does not belong to contract", http.StatusBadRequest)
		}

		err := ExecuteClause(stub, contract, depClause)
		if err != nil {
			return err
		}
	}

	inputs := utils.JoinMaps(clause.Input, clause.Parameters, contract.Data)

	action := params.Get(clause.ActionType)
	result, shouldFinalizeClause, err := action.Execute(inputs, contract.Data)
	if err != nil {
		return errors.WrapError(err, "Failed to execute action")
	}

	if len(result.Assets) > 0 {
		err = utils.SaveGeneratedAssets(stub, result.Assets, contract, clause)
		if err != nil {
			return errors.WrapError(err, "Failed to save generated assets")
		}
	}

	contract.Data = mergeData(contract.Data, result.Data)

	return updateClause(stub, clause, shouldFinalizeClause, result.Success, result.Feedback)
}

func updateClause(stub *sw.StubWrapper, clause *models.Clause, shouldFinalize bool, success bool, feedback string) errors.ICCError {
	_, err := clause.Asset.Update(stub, map[string]interface{}{
		"finalized": shouldFinalize,
		"result": map[string]interface{}{
			"success":  success,
			"feedback": feedback,
		},
	})

	clause.Finalized = shouldFinalize
	return err
}

func mergeData(existingData map[string]interface{}, newData map[string]interface{}) map[string]interface{} {
	for k, v := range newData {
		existingData[k] = v
	}
	return existingData
}
