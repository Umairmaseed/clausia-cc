package datatypes

import (
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

type ArgDt string

const (
	ArgDtString string = "string"
	ArgDtInt    string = "int"
	ArgDtFloat  string = "float64"
	ArgDtBool   string = "boolean"
	ArgDtTime   string = "datetime"
	ArgDtObject string = "object"
	ArgDtEnum   string = "enum"
)

func (b ArgDt) CheckType() errors.ICCError {
	switch string(b) {
	case ArgDtString, ArgDtInt, ArgDtFloat, ArgDtBool, ArgDtTime, ArgDtObject, ArgDtEnum:
		return nil
	default:
		return errors.NewCCError("invalid type", http.StatusBadRequest)
	}

}

var argDt = assets.DataType{
	DropDownValues: map[string]interface{}{
		string(ArgDtString): ArgDtString,
		string(ArgDtInt):    ArgDtInt,
		string(ArgDtFloat):  ArgDtFloat,
		string(ArgDtBool):   ArgDtBool,
		string(ArgDtTime):   ArgDtTime,
		string(ArgDtObject): ArgDtObject,
		string(ArgDtEnum):   ArgDtEnum,
	},

	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		dataValue, ok := data.(string)
		if !ok {
			return "", nil, errors.NewCCError("Property must be a string", http.StatusBadRequest)
		}

		retVal := ArgDt(dataValue)
		err := retVal.CheckType()

		return dataValue, dataValue, err
	},
}
