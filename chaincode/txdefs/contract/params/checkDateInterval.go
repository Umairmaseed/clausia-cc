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
	intervalDays   int = 1
	intervalWeeks  int = 2
	intervalMonths int = 3
	intervalYears  int = 4

	intervalBefore int = -1
	intervalAfter  int = 1

	defaultName string = "dateIntervalCheck"
)

type CheckDateInterval struct {
}

type ParametersCheckDateInterval struct {
	Name             string `json:"name"`
	IntervalType     int    `json:"intervalType"`
	DeadlineInterval int    `json:"deadlineInterval"`
	ReferenceDate    string `json:"referenceDate"` // Optional
}

type InputsCheckDateInterval struct {
	ReferenceDate string `json:"referenceDate"` // Optional
	EvaluatedDate string `json:"evaluatedDate"`
}

func (a *CheckDateInterval) Type() datatypes.ActionType {
	return datatypes.CheckDateInterval
}

func (a *CheckDateInterval) getFeedback(isOnTime bool) string {
	if isOnTime {
		return "Within the deadline."
	} else {
		return "Outside the deadline."
	}
}

func (a *CheckDateInterval) Execute(input interface{}, data map[string]interface{}) (*models.Result, bool, errors.ICCError) {
	// Marshal input to get bytes
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to marshal input")
	}

	// Unmarshal to get CheckDateIntervalArgs
	var args ParametersCheckDateInterval
	err = json.Unmarshal(inputBytes, &args)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal to CheckDateIntervalArgs")
	}

	// Unmarshal to get CheckDateIntervalInput
	var inputData InputsCheckDateInterval
	err = json.Unmarshal(inputBytes, &inputData)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to unmarshal to CheckDateIntervalInput")
	}

	// Default name if not provided
	if args.Name == "" {
		args.Name = defaultName
	}

	// Determine the ReferenceDate
	refDateStr := args.ReferenceDate
	if refDateStr == "" {
		// Fallback to input data's ReferenceDate if not provided in args
		refDateStr = inputData.ReferenceDate
	}

	if refDateStr == "" {
		return &models.Result{
			Success:  false,
			Feedback: "Reference date is not provided",
		}, false, nil
	}

	// Parse the ReferenceDate
	refDate, err := time.Parse(time.RFC3339, refDateStr)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to parse reference date")
	}

	// Get and parse the EvaluatedDate
	evaluatedDateStr := inputData.EvaluatedDate
	if evaluatedDateStr == "" {
		return &models.Result{
			Success:  false,
			Feedback: "Evaluated date is not provided",
		}, false, nil
	}

	evaluatedDate, err := time.Parse(time.RFC3339, evaluatedDateStr)
	if err != nil {
		return nil, false, errors.WrapError(err, "Failed to parse evaluated date")
	}

	// Calculate the deadline in days based on IntervalType
	days := args.DeadlineInterval
	switch args.IntervalType {
	case intervalWeeks:
		days *= 7
	case intervalMonths:
		days *= 30
	case intervalYears:
		days *= 365
	}

	// Calculate the interval date based on IntervalType
	var intervalDate time.Time
	switch args.IntervalType {
	case intervalBefore:
		intervalDate = refDate.AddDate(0, 0, -days)
	case intervalAfter:
		intervalDate = refDate.AddDate(0, 0, days)
	default:
		return nil, false, errors.NewCCError("Invalid interval type, it must be 'before' (-1) or 'after'(1)", http.StatusBadRequest)
	}

	// Determine if the evaluated date is within the deadline
	isWithinDeadline := evaluatedDate.Before(intervalDate) || evaluatedDate.Equal(intervalDate)
	feedback := a.getFeedback(isWithinDeadline)

	// Calculate the number of days from the deadline
	daysFromDeadline := int(intervalDate.Sub(evaluatedDate).Hours() / 24)
	if daysFromDeadline < 0 {
		daysFromDeadline = -daysFromDeadline
	}

	return &models.Result{
		Success:  isWithinDeadline,
		Feedback: feedback,
		Data: map[string]interface{}{
			args.Name: daysFromDeadline,
		},
	}, true, nil
}

func (a *CheckDateInterval) GetParameters() interface{} {
	return ParametersCheckDateInterval{}
}

func (a *CheckDateInterval) GetInputs() interface{} {
	return InputsCheckDateInterval{}
}
