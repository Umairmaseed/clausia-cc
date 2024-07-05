package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

type DocumentHistoryRecord struct {
	TxID      string          `json:"txId"`
	Timestamp string          `json:"timestamp"`
	Value     json.RawMessage `json:"value"`
	IsDeleted bool            `json:"isDeleted"`
}

var GetDocHistory = tx.Transaction{
	Tag:         "getDocHistory",
	Label:       "Get Document History",
	Description: "Retrieves document history by key",
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
		mapKey, ok := req["key"].(map[string]interface{})
		if !ok {
			return nil, errors.NewCCError("key is required and must be an object", http.StatusBadRequest)
		}

		key, ok := mapKey["@key"].(string)
		if !ok {
			return nil, errors.NewCCError("invalid type of key", http.StatusBadRequest)
		}

		historyIterator, err := stub.GetHistoryForKey(key)
		if err != nil {
			return nil, errors.WrapError(err, "failed to read document history from blockchain")
		}
		defer historyIterator.Close()

		var history []DocumentHistoryRecord
		for historyIterator.HasNext() {
			queryResponse, err := historyIterator.Next()
			if err != nil {
				return nil, errors.WrapError(err, "error iterating response")
			}

			var record DocumentHistoryRecord
			record.TxID = queryResponse.TxId
			record.Timestamp = queryResponse.Timestamp.String()
			record.Value = queryResponse.Value
			record.IsDeleted = queryResponse.IsDelete

			history = append(history, record)
		}

		responseJSON, nerr := json.Marshal(history)
		if nerr != nil {
			return nil, errors.WrapError(err, "error marshaling history response")
		}
		return responseJSON, nil
	},
}
