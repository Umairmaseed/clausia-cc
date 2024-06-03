package datatypes

import (
	"fmt"
	"strconv"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

type StatusType float64

const (
	waiting StatusType = iota
	cancelled
	expired
	finalized
	partiallyFinalized
)

func (b StatusType) CheckType() errors.ICCError {
	switch b {
	case waiting:
		return nil
	case cancelled:
		return nil
	case expired:
		return nil
	case finalized:
		return nil
	case partiallyFinalized:
		return nil
	default:
		return errors.NewCCError("invalid type", 400)
	}

}

var statusType = assets.DataType{
	AcceptedFormats: []string{"number"},
	DropDownValues: map[string]interface{}{
		"waiting for signatures": waiting,
		"cancelled":              cancelled,
		"expired":                expired,
		"finalized":              finalized,
		"partially finalized":    partiallyFinalized,
	},
	Description: "Status of the signature",
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

		retVal := (StatusType)(dataVal)
		err := retVal.CheckType()
		return fmt.Sprint(retVal), retVal, err
	},
}
