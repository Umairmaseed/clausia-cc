package models

import (
	"time"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

type AutoExecutableContract struct {
	Key           string                 `json:"key"`
	Name          string                 `json:"name"`
	SignatureDate string                 `json:"signatureDate"`
	Clauses       []*Clause              `json:"clauses"`
	Data          map[string]interface{} `json:"data"`
	Owner         assets.Key             `json:"owner"`
	Participants  []assets.Key           `json:"participants"`
	Asset         *assets.Asset
}

func GetAutoExecutableContract(stub *sw.StubWrapper, key assets.Key) (*AutoExecutableContract, errors.ICCError) {
	asset, err := key.Get(stub)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get auto executable contract asset")
	}

	var contract AutoExecutableContract

	prop := asset.GetProp("clauses")
	clauseKeys, _ := prop.([]interface{})
	for _, c := range clauseKeys {
		keyMap, _ := c.(map[string]interface{})
		key, err := assets.NewKey(keyMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to make clause key")
		}

		clause, err := GetClause(stub, key)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get clause")
		}

		contract.Clauses = append(contract.Clauses, clause)
	}

	data, ok := asset.GetProp("data").(map[string]interface{})
	if !ok {
		data = map[string]interface{}{}
	}

	contract.Key, _ = asset.GetProp("@key").(string)
	contract.Name, _ = asset.GetProp("name").(string)

	signatureDate, _ := asset.GetProp("signatureDate").(time.Time)
	contract.SignatureDate = signatureDate.Format(time.RFC3339)

	contract.Data = data

	ownerKey, _ := asset.GetProp("owner").(map[string]interface{})
	contract.Owner, _ = assets.NewKey(ownerKey)

	participantKeys, _ := asset.GetProp("participants").([]interface{})
	for _, p := range participantKeys {
		participantKeyMap, _ := p.(map[string]interface{})
		participantKey, err := assets.NewKey(participantKeyMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to make participant key")
		}
		contract.Participants = append(contract.Participants, participantKey)
	}

	contract.Asset = asset

	return &contract, nil
}

func (c *AutoExecutableContract) GetClause(key string) *Clause {
	for i := range c.Clauses {
		if c.Clauses[i].Key == key {
			return c.Clauses[i]
		}
	}
	return nil
}
