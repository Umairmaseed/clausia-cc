package params

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Date                time.Time `json:"date"`
	Payment             float64   `json:"payment"`
	ReceiptHash         string    `json:"receiptHash"`
	ReceiptUrl          string    `json:"receiptUrl"`
	FinalPayment        bool      `json:"finalPayment"`
	StripeToken         string    `json:"stripeToken"`
	PayPalTransactionID string    `json:"payPalTransactionID"`
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

	// Calculate total amount including bonuses and fines
	totalAmount, bonusPayment, finePayment, err := a.calculateTotalAmount(params, data)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to calculate total amount")
	}

	// Update payment data
	a.updatePaymentData(data, params, inputs)

	// Create result struct
	result := a.createResult(data, params, inputs, bonusPayment, finePayment)

	// Process the payment based on whether it's partial or full
	success, feedback := a.processPayment(params, inputs, totalAmount, data)
	result.Feedback = feedback
	result.Success = success

	// Return result
	return &result, success, nil
}

func (a *MakePaymentClause) calculateTotalAmount(params MakePaymentParams, data map[string]interface{}) (float64, float64, float64, errors.ICCError) {
	var bonusAmount, fineAmount, bonusPaidAmount, finePaidAmount float64

	if params.AddBonus {
		bonusAmount = a.getAmountFromData(data, "bonus")
		bonusPaidAmount = a.getAmountFromData(data, "bonusPaid")
	}
	if params.AddFine {
		fineAmount = a.getAmountFromData(data, "fine")
		finePaidAmount = a.getAmountFromData(data, "finePaid")
	}

	remainingBonus := bonusAmount - bonusPaidAmount
	remainingFine := fineAmount - finePaidAmount

	currBonusPayment := remainingBonus * (params.PaymentRate / 100)
	currFinePayment := remainingFine * (params.PaymentRate / 100)

	totalAmount := params.Amount + currBonusPayment - currFinePayment

	return totalAmount, currBonusPayment, currFinePayment, nil
}

func (a *MakePaymentClause) getAmountFromData(data map[string]interface{}, key string) float64 {
	if amount, exists := data[key].(float64); exists {
		return amount
	}
	return 0.0
}

func (a *MakePaymentClause) updatePaymentData(data map[string]interface{}, params MakePaymentParams, inputs MakePaymentInputs) {
	paymentData := map[string]interface{}{
		"date":    inputs.Date,
		"payment": inputs.Payment,
	}

	if inputs.ReceiptHash != "" {
		paymentData["receiptHash"] = inputs.ReceiptHash
	}
	if inputs.ReceiptUrl != "" {
		paymentData["receiptUrl"] = inputs.ReceiptUrl
	}

	if inputs.StripeToken != "" {
		paymentData["stripeToken"] = inputs.StripeToken
	}
	if inputs.PayPalTransactionID != "" {
		paymentData["payPalTransactionID"] = inputs.PayPalTransactionID
	}

	data[params.Name] = paymentData
}

func (a *MakePaymentClause) createResult(data map[string]interface{}, params MakePaymentParams, inputs MakePaymentInputs, bonusPayment, finePayment float64) models.Result {
	paymentID := generatePaymentID(params, inputs)

	assetData := map[string]interface{}{
		"@assetType": "payment",
		"hash":       paymentID,
		"name":       params.Name,
		"payment":    inputs.Payment,
	}

	if inputs.ReceiptUrl != "" {
		assetData["receiptUrl"] = inputs.ReceiptUrl
	}

	if inputs.StripeToken != "" {
		assetData["stripeToken"] = inputs.StripeToken
	}

	if inputs.PayPalTransactionID != "" {
		assetData["payPalTransactionID"] = inputs.PayPalTransactionID
	}

	result := models.Result{
		Data: data,
		Meta: map[string]interface{}{
			"payment": inputs.Payment,
			"bonus":   bonusPayment,
			"fine":    finePayment,
		},
		Assets: []map[string]interface{}{
			assetData,
		},
	}

	if params.AddBonus {
		result.Data["bonusPaid"] = a.getAmountFromData(data, "bonusPaid") + bonusPayment
	}
	if params.AddFine {
		result.Data["finePaid"] = a.getAmountFromData(data, "finePaid") + finePayment
	}

	return result
}

func generatePaymentID(params MakePaymentParams, inputs MakePaymentInputs) string {
	if inputs.ReceiptHash != "" {
		return inputs.ReceiptHash
	}

	var uniqueData string
	if inputs.StripeToken != "" {
		uniqueData = inputs.StripeToken
	} else if inputs.PayPalTransactionID != "" {
		uniqueData = inputs.PayPalTransactionID
	} else {
		// Fallback to creating a hash based on the payment name and amount
		uniqueData = fmt.Sprintf("%s-%f", params.Name, inputs.Payment)
	}

	hash := sha256.Sum256([]byte(uniqueData))
	return hex.EncodeToString(hash[:])
}

func (a *MakePaymentClause) processPayment(params MakePaymentParams, inputs MakePaymentInputs, totalAmount float64, data map[string]interface{}) (bool, string) {
	if params.PartialPayment {
		return a.processPartialPayment(params, inputs, totalAmount, data)
	}
	return a.processFullPayment(inputs, totalAmount, data)
}

func (a *MakePaymentClause) processPartialPayment(params MakePaymentParams, inputs MakePaymentInputs, totalAmount float64, data map[string]interface{}) (bool, string) {
	previousPayment := a.getAmountFromData(data, "previousPartialPayment")
	partialPaymentAmount := totalAmount * (params.PaymentRate / 100)

	data["previousPartialPayment"] = previousPayment + inputs.Payment

	if inputs.Payment < partialPaymentAmount {
		return false, "Partial payment is less than expected. Payment incomplete."
	}
	return true, "Partial payment successful."
}

func (a *MakePaymentClause) processFullPayment(inputs MakePaymentInputs, totalAmount float64, data map[string]interface{}) (bool, string) {
	if inputs.Payment < totalAmount {
		data["paidAmount"] = totalAmount
		return false, "Payment is less than the required amount. Payment incomplete."
	}
	data["paidAmount"] = totalAmount
	return true, "Full payment successful."
}
