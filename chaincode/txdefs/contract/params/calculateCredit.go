package params

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/models"
)

const (
	creditName = "calculateCredit"
)

type CalculateCreditParam struct {
	ImposeCredit    bool    `json:"imposeCredit"`
	CreditName      string  `json:"creditName"`
	Percentage      float64 `json:"percentage"`
	PredefinedValue float64 `json:"predefinedValue"`
	ConditionName   string  `json:"conditionName"`
}

type CalculateCreditInput struct {
	ConditionValue interface{} `json:"conditionValue"`
	StoredValue    float64     `json:"storedValue"`
}

type CalculateCredit struct{}

func (a *CalculateCredit) Type() datatypes.ActionType {
	return datatypes.GetCredit
}

func (a *CalculateCredit) GetParameters() interface{} {
	return CalculateCreditParam{}
}

func (a *CalculateCredit) GetInputs() interface{} {
	return CalculateCreditInput{}
}

func (a *CalculateCredit) Execute(input interface{}) (*models.Result, bool, errors.ICCError) {
	// Unmarshal input
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, false, errors.WrapError(err, "failed to marshal input")
	}

	var creditInput CalculateCreditInput
	err = json.Unmarshal(inputBytes, &creditInput)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal input")
	}

	// Separate unmarshal parameters
	var parameters CalculateCreditParam
	err = json.Unmarshal(inputBytes, &parameters)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal input")
	}

	if !parameters.ImposeCredit {
		return &models.Result{
			Success:  false,
			Feedback: "Credit is not imposed.",
		}, false, nil
	}

	// Check the condition
	conditionMet := false
	switch v := creditInput.ConditionValue.(type) {
	case bool:
		if v {
			conditionMet = true
		}
	default:
		conditionMet = creditInput.ConditionValue != nil
	}

	if !conditionMet {
		return &models.Result{
			Success:  false,
			Feedback: "Conditions for credit are not met.",
		}, false, nil
	}

	var creditAmount float64
	if parameters.Percentage > 0 && creditInput.StoredValue > 0 {
		creditAmount = (parameters.Percentage / 100) * creditInput.StoredValue
	} else {
		creditAmount = parameters.PredefinedValue
	}

	if creditAmount <= 0 {
		return nil, false, errors.NewCCError("Invalid credit amount calculated", http.StatusBadRequest)
	}

	// Handle empty credit name
	if parameters.CreditName == "" {
		parameters.CreditName = creditName
	}

	// Prepare result
	result := &models.Result{
		Success:  true,
		Feedback: "Credit calculated successfully.",
		Data: map[string]interface{}{
			parameters.CreditName: creditAmount,
		},
	}

	return result, true, nil
}
