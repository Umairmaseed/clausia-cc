package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// GetDoc retrieves a document by its key
var GetDoc = tx.Transaction{
	Tag:         "getDoc",
	Label:       "Get Document",
	Description: "Retrieves a document by key",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:         "key",
			Label:       "Document Key",
			Description: "The key of the document to retrieve",
			DataType:    "@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		key, ok := req["key"].(map[string]interface{})
		if !ok {
			return nil, errors.NewCCError("key is required and must be an object", http.StatusBadRequest)
		}

		query := map[string]interface{}{
			"selector": key,
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for document", http.StatusInternalServerError)
		}

		responseJSON, er := json.Marshal(response)
		if er != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", http.StatusInternalServerError)
		}

		return responseJSON, nil
	},
}
