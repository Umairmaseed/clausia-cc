package contract

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var RemoveTemplateClause = tx.Transaction{
	Tag:         "removeTemplateClause",
	Label:       "Remove Template Clause",
	Description: "Remove or delete a clause from a template",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "template",
			Label:    "Template",
			DataType: "->template",
		},
		{
			Required: true,
			Tag:      "templateClause",
			Label:    "Template Clause",
			DataType: "->templateClause",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		templateKey, ok := req["template"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'template' must be an asset key")
		}

		templateAsset, err := templateKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get template asset from ledger")
		}

		templateClauseKey, ok := req["templateClause"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter 'templateClause' must be an asset key")
		}

		templateClauseAsset, err := templateClauseKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause template asset from ledger")
		}

		clauses, ok := (*templateAsset)["clauses"].([]interface{})
		if !ok {
			return nil, errors.WrapError(nil, "Clauses field is not an array")
		}

		var updatedClauses []interface{}
		var clauseFound bool
		for _, c := range clauses {
			clauseMap, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			if clauseMap["@key"] != (*templateClauseAsset)["@key"] {
				updatedClauses = append(updatedClauses, c)
			} else {
				clauseFound = true
			}
		}

		if !clauseFound {
			return nil, errors.NewCCError("Template Clause does not belong to the template", http.StatusBadRequest)
		}

		contractUpdates := map[string]interface{}{
			"clauses": updatedClauses,
		}

		updatedTemplateAsset, err := templateKey.Update(stub, contractUpdates)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update template asset in ledger")
		}

		_, err = templateClauseKey.Delete(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to delete template clause")
		}

		updatedTemplateJSON, nerr := json.Marshal(updatedTemplateAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "Failed to marshal updated template asset")
		}

		return updatedTemplateJSON, nil
	},
}
