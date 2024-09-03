package params

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/models"
)

type FinalizeContract struct{}

type DataType string

const (
	IntType    DataType = "int"
	FloatType  DataType = "float"
	StringType DataType = "string"
	BoolType   DataType = "bool"
	DateType   DataType = "date"
)

type ConditionalCheck string

const (
	Equal    ConditionalCheck = "equal"
	NotEqual ConditionalCheck = "notEqual"
	Greater  ConditionalCheck = "greater"
	Smaller  ConditionalCheck = "smaller"
)

type checkValue struct {
	Tag              string           `json:"tag"`
	ReferenceValue   interface{}      `json:"referenceValue"`
	DataType         DataType         `json:"dataType"`
	ConditionalCheck ConditionalCheck `json:"conditionalCheck"`
}

type FinalizeContractParams struct {
	AutoFinalizationValue  checkValue `json:"autoFinalizationValue"`
	CancellationCheckValue checkValue `json:"cancellationCheckValue"`
	ForceCancellation      bool       `json:"forceCancellation"`
	RequestedCancellation  bool       `json:"requestedCancellation"`
}

func (a *FinalizeContract) Type() datatypes.ActionType {
	return datatypes.FinishContract
}

func (a *FinalizeContract) GetParameters() interface{} {
	return FinalizeContractParams{}
}

func (a *FinalizeContract) GetInputs() interface{} {
	return nil
}

func (a *FinalizeContract) Execute(input interface{}, data map[string]interface{}) (*models.Result, bool, errors.ICCError) {

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to marshal input")
	}

	var params FinalizeContractParams
	err = json.Unmarshal(inputBytes, &params)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal to FinalizeContractParams")
	}

	if params.ForceCancellation {
		return &models.Result{
			Success:  true,
			Feedback: "Contract cancelled upon force cancellation",
		}, true, nil
	}

	if (params.CancellationCheckValue != checkValue{}) {
		if params.RequestedCancellation {
			if err := a.evaluateCheckValue(params.CancellationCheckValue, data); err == nil {
				return &models.Result{
					Success:  true,
					Feedback: "Contract cancelled upon request based on defined conditions",
				}, true, nil
			} else {
				return nil, false, errors.WrapError(err, "Contract cancellation condition not met")
			}
		}
	}

	if (params.AutoFinalizationValue != checkValue{}) {
		if err := a.evaluateCheckValue(params.AutoFinalizationValue, data); err == nil {
			return &models.Result{
				Success:  true,
				Feedback: "Contract automatically finalized based on defined conditions.",
			}, true, nil
		} else {
			return nil, false, errors.WrapError(err, "Contract finalization condition not met")
		}
	}

	return &models.Result{
		Success:  false,
		Feedback: "Contract remains active; no conditions for finalization were met.",
	}, false, nil
}

func (a *FinalizeContract) evaluateCheckValue(cv checkValue, data map[string]interface{}) errors.ICCError {
	value, exists := data[cv.Tag]
	if !exists {
		return errors.NewCCError("Value not found in contract data", 400)
	}

	switch cv.DataType {
	case IntType:
		return a.evaluateInt(cv, value)
	case FloatType:
		return a.evaluateFloat(cv, value)
	case BoolType:
		return a.evaluateBool(cv, value)
	case DateType:
		return a.evaluateDate(cv, value)
	case StringType:
		return a.evaluateString(cv, value)
	default:
		return errors.NewCCError("Unsupported data type", 400)
	}
}

func (a *FinalizeContract) evaluateInt(cv checkValue, value interface{}) errors.ICCError {
	intVal, ok := value.(int)
	if !ok {
		return errors.NewCCError("Value must be of type int", 400)
	}

	val, ok := cv.ReferenceValue.(int)
	if !ok {
		return errors.NewCCError("Reference value must be of int type", 400)
	}

	switch cv.ConditionalCheck {
	case Equal:
		if intVal == val {
			return nil
		}
	case NotEqual:
		if intVal != val {
			return nil
		}
	case Greater:
		if intVal > val {
			return nil
		}
	case Smaller:
		if intVal < val {
			return nil
		}
	}

	return errors.NewCCError("Condition not met for int evaluation", 400)
}

func (a *FinalizeContract) evaluateFloat(cv checkValue, value interface{}) errors.ICCError {
	floatVal, ok := value.(float64)
	if !ok {
		return errors.NewCCError("Invalid float value", 400)
	}

	val, ok := cv.ReferenceValue.(float64)
	if !ok {
		return errors.NewCCError("Reference value must be of type float", 400)
	}

	switch cv.ConditionalCheck {
	case Equal:
		if floatVal == val {
			return nil
		}
	case NotEqual:
		if floatVal != val {
			return nil
		}
	case Greater:
		if floatVal > val {
			return nil
		}
	case Smaller:
		if floatVal < val {
			return nil
		}
	}

	return errors.NewCCError("Condition not met for float evaluation", 400)
}

func (a *FinalizeContract) evaluateBool(cv checkValue, value interface{}) errors.ICCError {
	boolVal, ok := value.(bool)
	if !ok {
		return errors.NewCCError("Invalid bool value", 400)
	}

	expectedVal, ok := cv.ReferenceValue.(bool)
	if !ok {
		return errors.NewCCError("Reference value must be of type boolean", 400)
	}

	if cv.ConditionalCheck == Equal && boolVal == expectedVal {
		return nil
	}

	if cv.ConditionalCheck == NotEqual && boolVal != expectedVal {
		return nil
	}

	return errors.NewCCError("Condition not met for bool evaluation", 400)
}

func (a *FinalizeContract) evaluateDate(cv checkValue, value interface{}) errors.ICCError {
	dateVal, err := time.Parse(time.RFC3339, value.(string))
	if err != nil {
		return errors.WrapError(err, "Invalid date value")
	}

	refDateStr, ok := cv.ReferenceValue.(string)
	if !ok {
		return errors.NewCCError("Reference value must be of type string", 400)
	}

	expectedDate, err := time.Parse(time.RFC3339, refDateStr)
	if err != nil {
		return errors.WrapError(err, "Invalid condition for date")
	}

	switch cv.ConditionalCheck {
	case Equal:
		if dateVal.Equal(expectedDate) {
			return nil
		}
	case NotEqual:
		if !dateVal.Equal(expectedDate) {
			return nil
		}
	case Greater:
		if dateVal.After(expectedDate) {
			return nil
		}
	case Smaller:
		if dateVal.Before(expectedDate) {
			return nil
		}
	}

	return errors.NewCCError("Condition not met for date evaluation", 400)
}

func (a *FinalizeContract) evaluateString(cv checkValue, value interface{}) errors.ICCError {
	strVal, ok := value.(string)
	if !ok {
		return errors.NewCCError("Value must be of type string", 400)
	}

	refVal, ok := cv.ReferenceValue.(string)
	if !ok {
		return errors.NewCCError("Reference value must be of type string", 400)
	}

	switch cv.ConditionalCheck {
	case Equal:
		if strVal == refVal {
			return nil
		}
	case NotEqual:
		if strVal != refVal {
			return nil
		}
	}

	return errors.NewCCError("Condition not met for string evaluation", 400)
}
