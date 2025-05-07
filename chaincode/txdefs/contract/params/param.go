package params

import (
	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/clausia-cc/chaincode/datatypes"
	"github.com/hyperledger-labs/clausia-cc/chaincode/txdefs/contract/models"
)

type param interface {
	Type() datatypes.ActionType

	Execute(args interface{}, data map[string]interface{}) (*models.Result, bool, errors.ICCError)

	GetParameters() interface{}

	GetInputs() interface{}
}

type Output struct {
	Name  string
	Type  string
	Label string
}

func Get(actionType datatypes.ActionType) param {
	switch actionType {
	case datatypes.CheckDateInterval:
		return &CheckDateInterval{}
	case datatypes.GetDeduction:
		return &CalculateFine{}
	case datatypes.GetCredit:
		return &CalculateCredit{}
	case datatypes.Payment:
		return &MakePaymentClause{}
	case datatypes.FinishContract:
		return &FinalizeContract{}
	default:
		return nil // to be changed according to non executable action type
	}
}
