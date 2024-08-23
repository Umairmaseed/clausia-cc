package params

import (
	"encoding/json"
	"net/http"
	"time"

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
	ReviewCondition bool    `json:"reviewCondition"`
}

type CalculateCreditInput struct {
	Review      Review  `json:"review"`
	StoredValue float64 `json:"storedValue"`
}

type Review struct {
	User       map[string]interface{} `json:"user"`
	Rating     int                    `json:"rating"`
	Comments   string                 `json:"comments"`
	Date       time.Time              `json:"date"`
	ContractID map[string]interface{} `json:"contract_id"`
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

func (a *CalculateCredit) Execute(input interface{}, data map[string]interface{}) (*models.Result, bool, errors.ICCError) {
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

	if parameters.ImposeCredit || parameters.ReviewCondition {
		var creditAmount float64

		// If ImposeCredit is true, calculate credit based on the provided parameters
		if parameters.ImposeCredit {
			if parameters.Percentage > 0 && creditInput.StoredValue > 0 {
				creditAmount = (parameters.Percentage / 100) * creditInput.StoredValue
			} else {
				creditAmount = parameters.PredefinedValue
			}

			if creditAmount <= 0 {
				return nil, false, errors.NewCCError("Invalid credit amount calculated", http.StatusBadRequest)
			}

			if parameters.CreditName == "" {
				parameters.CreditName = creditName
			}

			feedback := "Credit calculated successfully."

			updateData := updateBonusData(data, creditAmount, parameters.CreditName, feedback)

			return &models.Result{
				Success:  true,
				Feedback: feedback,
				Data:     updateData,
			}, true, nil
		}

		// If ReviewCondition is true, check the Review rating
		if parameters.ReviewCondition && creditInput.Review.Rating >= 3 {
			if parameters.Percentage > 0 && creditInput.StoredValue > 0 {
				creditAmount = (parameters.Percentage / 100) * creditInput.StoredValue
			} else {
				creditAmount = parameters.PredefinedValue
			}

			if creditAmount <= 0 {
				return nil, false, errors.NewCCError("Invalid credit amount calculated", http.StatusBadRequest)
			}

			if parameters.CreditName == "" {
				parameters.CreditName = creditName
			}

			feedback := "Credit calculated based on review rating."

			updateData := updateBonusData(data, creditAmount, parameters.CreditName, feedback)

			return &models.Result{
				Success:  true,
				Feedback: feedback,
				Data:     updateData,
			}, true, nil
		}
	}

	// If neither condition is met, return no credit calculation
	return &models.Result{
		Success:  false,
		Feedback: "Conditions for credit are not met.",
	}, false, nil
}

func updateBonusData(data map[string]interface{}, creditAmount float64, creditName, feedback string) map[string]interface{} {
	// Update the "bonus" field if it exists
	if currentBonus, exists := data["bonus"]; exists {
		if bonusValue, ok := currentBonus.(float64); ok {
			data["bonus"] = bonusValue + creditAmount
		}
	} else {
		data["bonus"] = creditAmount
	}

	// Create a new entry for the bonus
	newBonusEntry := map[string]interface{}{
		"name":     creditName,
		"bonus":    creditAmount,
		"feedback": feedback,
		"success":  true,
	}

	// Update the "listOfBonus" field
	if listOfBonus, exists := data["listOfBonus"]; exists {
		if bonuses, ok := listOfBonus.([]map[string]interface{}); ok {
			data["listOfBonus"] = append(bonuses, newBonusEntry)
		} else {
			data["listOfBonus"] = []map[string]interface{}{newBonusEntry}
		}
	} else {
		data["listOfBonus"] = []map[string]interface{}{newBonusEntry}
	}

	return data
}
