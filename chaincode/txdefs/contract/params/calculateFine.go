package params

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/models"
)

const (
	fineName = "calculateFine"
)

type CalculateFineParameters struct {
	ImposeFine        bool    `json:"imposeFine"`
	FineName          string  `json:"fineName"`
	MaxPercentage     float64 `json:"maxPercentage"`
	MaxReferenceValue float64 `json:"maxReferenceValue"`
}

type CalculateFineInput struct {
	ReferenceValue  float64 `json:"referenceValue"`
	DailyPercentage float64 `json:"dailyPercentage"`
	Days            float64 `json:"days"`
}

type CalculateFine struct{}

func (a *CalculateFine) Type() datatypes.ActionType {
	return datatypes.GetDeduction
}

func (a *CalculateFine) GetParameters() interface{} {
	return CalculateFineParameters{}
}

func (a *CalculateFine) GetInputs() interface{} {
	return CalculateFineInput{}
}

func (a *CalculateFine) Execute(input interface{}, data map[string]interface{}) (*models.Result, bool, errors.ICCError) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, false, errors.WrapError(err, "failed to marshal input")
	}

	var fineInput CalculateFineInput
	err = json.Unmarshal(inputBytes, &fineInput)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal input")
	}

	var fineParams CalculateFineParameters
	err = json.Unmarshal(inputBytes, &fineParams)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal parameters")
	}

	if !fineParams.ImposeFine {
		return &models.Result{
			Success:  true,
			Feedback: "Fine calculation is not imposed.",
		}, false, nil
	}

	if fineInput.ReferenceValue <= 0 || fineInput.DailyPercentage <= 0 || fineInput.Days <= 0 {
		return nil, false, errors.NewCCError("Invalid input values: ReferenceValue, DailyPercentage, and Days should be greater than 0", http.StatusBadRequest)
	}

	// Calculate the fine
	fine := fineInput.ReferenceValue * fineInput.DailyPercentage / 100 * fineInput.Days

	// Apply upper limit if necessary
	var shouldConsiderUpperLimit bool
	if fineParams.MaxPercentage > 0 && fineParams.MaxReferenceValue > 0 {
		shouldConsiderUpperLimit = true
	}

	if shouldConsiderUpperLimit {
		limit := fineParams.MaxPercentage / 100 * fineParams.MaxReferenceValue
		if fine > limit {
			fine = limit
		}
	}

	// Handle empty fine name
	if fineParams.FineName == "" {
		fineParams.FineName = fineName
	}

	updateData := data

	// Add fine to the "fine" field if it exists
	if currentFine, exists := updateData["fine"]; exists {
		if fineValue, ok := currentFine.(float64); ok {
			updateData["fine"] = fineValue + fine
		}
	} else {
		updateData["fine"] = fine
	}

	// Add the current fine to the "listOfFines" field
	newFineEntry := map[string]interface{}{
		"name":     fineParams.FineName,
		"fine":     fine,
		"feedback": "Fine calculated successfully.",
		"success":  true,
	}

	if listOfFines, exists := updateData["listOfFines"]; exists {
		if fines, ok := listOfFines.([]map[string]interface{}); ok {
			updateData["listOfFines"] = append(fines, newFineEntry)
		} else {
			updateData["listOfFines"] = []map[string]interface{}{newFineEntry}
		}
	} else {
		updateData["listOfFines"] = []map[string]interface{}{newFineEntry}
	}

	// Prepare the result with updated data
	result := &models.Result{
		Success:  true,
		Feedback: "Fine calculated successfully.",
		Data:     updateData,
	}

	return result, true, nil
}
