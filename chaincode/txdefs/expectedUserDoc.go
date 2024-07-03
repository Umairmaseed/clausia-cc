package txdefs

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/datatypes"
)

var ExpectedUserDoc = tx.Transaction{
	Tag:         "expectedUserDoc",
	Label:       "Expected User Document",
	Description: "",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:      "signer",
			Label:    "Signer",
			Required: true,
			DataType: "->signer",
		},
		{
			Tag:      "status",
			Label:    "Status",
			DataType: "statusType",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		signerObj, signerOk := req["signer"].(assets.Key)
		if !signerOk {
			return nil, errors.NewCCError("Invalid document key", http.StatusBadRequest)
		}

		signerKey, keyOk := signerObj["@key"].(string)
		if !keyOk {
			return nil, errors.NewCCError("Invalid signer key", http.StatusBadRequest)
		}

		// Check if status is provided and convert it to string
		statusStr := ""
		if status, statusOk := req["status"].(datatypes.StatusType); statusOk {
			statusStr = strconv.Itoa(int(status))
		}

		// Construct the query string based on whether status is available
		queryString := `{"selector":{"@assetType":"document","requiredSignatures":{"$elemMatch":{"$eq":{"@assetType":"signer","@key":"` + signerKey + `"}}}`
		if statusStr != "" {
			queryString += `,"status":` + statusStr
		}
		queryString += `}}`

		resultsIterator, err := stub.GetQueryResult(queryString)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to get query result", http.StatusInternalServerError)
		}
		defer resultsIterator.Close()

		// Collect the results
		var searchResult []map[string]interface{}
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "error iterating query response", http.StatusInternalServerError)
			}

			var data map[string]interface{}
			err = json.Unmarshal(queryResponse.Value, &data)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal query response value", http.StatusInternalServerError)
			}

			searchResult = append(searchResult, data)
		}

		// Prepare the response
		response := struct {
			Result []map[string]interface{} `json:"result"`
		}{
			Result: searchResult,
		}

		responseJSON, errr := json.Marshal(response)
		if errr != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", http.StatusInternalServerError)
		}

		return responseJSON, nil
	},
}
