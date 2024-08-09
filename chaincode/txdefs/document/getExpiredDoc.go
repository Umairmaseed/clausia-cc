package document

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var GetExpiredDoc = tx.Transaction{
	Tag:         "getExpiredDocuments",
	Label:       "Get Expired Documents",
	Description: "Return Expired documents",
	Method:      "GET",

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": "document",
				"status":     float64(0),
			},
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for documents", http.StatusInternalServerError)
		}

		if response == nil || len(response.Result) == 0 {
			return []byte("[]"), nil
		}

		currentTime := time.Now().UTC()

		var filteredDocs []map[string]interface{}
		for _, doc := range response.Result {
			timeoutStr, ok := doc["timeout"].(string)
			if !ok {
				continue
			}

			timeout, err := time.Parse(time.RFC3339, timeoutStr)
			if err != nil {
				continue
			}

			if timeout.Before(currentTime) {
				filteredDocs = append(filteredDocs, doc)
			}
		}

		filteredJSON, er := json.Marshal(filteredDocs)
		if er != nil {
			return nil, errors.WrapErrorWithStatus(er, "error marshaling filtered response", http.StatusInternalServerError)
		}

		return filteredJSON, nil
	},
}
