package datatypes

import (
	"fmt"
	"strconv"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

type ActionType float64

const (
	checkDateInterval ActionType = iota
	getDeduction
	getCredit
	payment
	finishContract
)

func (b ActionType) CheckType() errors.ICCError {
	switch b {
	case checkDateInterval:
		return nil
	case getDeduction:
		return nil
	case getCredit:
		return nil
	case payment:
		return nil
	case finishContract:
		return nil
	default:
		return errors.NewCCError("invalid type", 400)
	}

}

var actionType = assets.DataType{
	AcceptedFormats: []string{"number"},
	DropDownValues: map[string]interface{}{
		"check date interval": checkDateInterval,
		"get Deduction":       getDeduction,
		"get Credit":          getCredit,
		"payment":             payment,
		"finish contract":     finishContract,
	},
	Description: "action type for clause",
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		var dataVal float64
		switch v := data.(type) {
		case float64:
			dataVal = v
		case int:
			dataVal = (float64)(v)
		case StatusType:
			dataVal = (float64)(v)
		case string:
			var err error
			dataVal, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return "", nil, errors.WrapErrorWithStatus(err, "asset property must be an integer, is %t", 400)
			}
		default:
			return "", nil, errors.NewCCError("asset property must be an integer, is %t", 400)
		}

		retVal := (ActionType)(dataVal)
		err := retVal.CheckType()
		return fmt.Sprint(retVal), retVal, err
	},
}
