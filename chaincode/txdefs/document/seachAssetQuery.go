package document

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var SearchAssetQuery = tx.Transaction{
	Tag:         "searchAssetQuery",
	Label:       "Search for asset through query",
	Description: "",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:      "query",
			Label:    "query",
			Required: true,
			DataType: "@object",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		query, ok := req["query"].(map[string]interface{})
		if !ok {
			return nil, errors.NewCCError("Invalid query type", http.StatusBadRequest)
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for query", http.StatusInternalServerError)
		}
		responseJSON, errr := json.Marshal(response)
		if errr != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", http.StatusInternalServerError)
		}

		return responseJSON, nil
	},
}
