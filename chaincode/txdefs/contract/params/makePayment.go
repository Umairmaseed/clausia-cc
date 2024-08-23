package params

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/models"
)

const (
	paymentName string = "paymentMade"
)

type MakePaymentClause struct{}

type MakePaymentParams struct {
	Name           string  `json:"name"`
	Amount         float64 `json:"amount"`
	PaymentRate    float64 `json:"paymentRate"`
	PartialPayment bool    `json:"partialPayment"`
	AddBonus       bool    `json:"addBonus"`
	AddFine        bool    `json:"addFine"`
}

type MakePaymentInputs struct {
	Date         time.Time `json:"date"`
	Payment      float64   `json:"payment"`
	ReceiptHash  string    `json:"receiptHash"`
	ReceiptUrl   string    `json:"receiptUrl"`
	FinalPayment bool      `json:"finalPayment"`
}

func (a *MakePaymentClause) Type() datatypes.ActionType {
	return datatypes.Payment
}

func (a *MakePaymentClause) GetParameters() interface{} {
	return MakePaymentParams{}
}

func (a *MakePaymentClause) GetInputs() interface{} {
	return MakePaymentInputs{}
}

func (a *MakePaymentClause) Execute(input interface{}, data map[string]interface{}) (*models.Result, bool, errors.ICCError) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to marshal input")
	}

	var inputs MakePaymentInputs
	err = json.Unmarshal(inputBytes, &inputs)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal MakePaymentInputs")
	}

	var params MakePaymentParams
	err = json.Unmarshal(inputBytes, &params)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal MakePaymentParams")
	}

	// Retrieve bonus and fine from data
	var bonusAmount, fineAmount, bonusPaidAmount, finePaidAmount float64
	if params.AddBonus {
		if bonus, exists := data["bonus"].(float64); exists {
			bonusAmount = bonus
		}
		if bonusPaid, exists := data["bonusPaid"].(float64); exists {
			bonusPaidAmount = bonusPaid
		}
	}
	if params.AddFine {
		if fine, exists := data["fine"].(float64); exists {
			fineAmount = fine
		}
		if finePaid, exists := data["finePaid"].(float64); exists {
			finePaidAmount = finePaid
		}
	}

	if params.Name == "" {
		params.Name = paymentName
	}

	remainingBonus := bonusAmount - bonusPaidAmount

	remainingFine := fineAmount - finePaidAmount

	currBonusPayment := remainingBonus * (params.PaymentRate / 100)

	currFinePayment := remainingFine * (params.PaymentRate / 100)

	// Adjust the amount based on bonus and fine
	totalAmount := params.Amount + currBonusPayment - currFinePayment

	updateData := data

	updateData[params.Name] = map[string]interface{}{
		"receiptHash": inputs.ReceiptHash,
		"receiptUrl":  inputs.ReceiptUrl,
		"date":        inputs.Date,
		"payment":     inputs.Payment,
	}

	// Result to be returned
	result := models.Result{
		Data: updateData,
		Meta: map[string]interface{}{
			"payment": inputs.Payment,
			"bonus":   currBonusPayment,
			"fine":    currFinePayment,
		},
		Assets: []map[string]interface{}{
			{
				"@assetType":  "payment",
				"name":        params.Name,
				"receiptHash": inputs.ReceiptHash,
				"receiptUrl":  inputs.ReceiptUrl,
				"payment":     inputs.Payment,
			},
		},
	}

	// Update remaining bonus and fine
	if params.AddBonus {
		result.Data["bonusPaid"] = bonusPaidAmount + currBonusPayment
	}
	if params.AddFine {
		result.Data["finePaid"] = finePaidAmount + currFinePayment
	}

	// Handle partial payment
	if params.PartialPayment {
		// Check for previous partial payments
		previousPayment := 0.0
		if prev, exists := data["previousPartialPayment"]; exists {
			previousPayment, _ = prev.(float64)
		}
		partialPaymentAmount := totalAmount * (params.PaymentRate / 100)

		// Update result with the current partial payment
		result.Data["previousPartialPayment"] = previousPayment + inputs.Payment

		if inputs.Payment < partialPaymentAmount {
			result.Feedback = "Partial payment is less than expected. Payment incomplete."
			result.Success = false
			return &result, false, nil
		} else {
			result.Feedback = "Partial payment successful."
			result.Success = true
			return &result, true, nil
		}
	} else {
		// Handle full payment
		if inputs.Payment < totalAmount {
			result.Feedback = "Payment is less than the required amount. Payment incomplete."
			result.Data["paidAmount"] = totalAmount
			result.Success = false
			return &result, false, nil
		} else {
			result.Feedback = "Full payment successful."
			result.Data["paidAmount"] = totalAmount
			result.Success = true
			return &result, true, nil
		}
	}
}
