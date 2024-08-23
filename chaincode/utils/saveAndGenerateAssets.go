package utils

import (
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/txdefs/contract/models"
)

func SaveGeneratedAssets(stub *sw.StubWrapper, genAssets []map[string]interface{}, contract *models.AutoExecutableContract, c *models.Clause) errors.ICCError {
	for _, a := range genAssets {
		// Add contract and clause infos
		a["autoExecutableContract"] = map[string]interface{}{
			"@assetType": "autoExecutableContract",
			"@key":       contract.Key,
		}
		a["clause"] = map[string]interface{}{
			"@assetType": "clause",
			"@key":       c.Key,
		}

		// Save on ledger
		asset, err := assets.NewAsset(a)
		if err != nil {
			return errors.WrapError(err, "Failed to create asset")
		}

		_, err = asset.PutNew(stub)
		if err != nil {
			return errors.WrapError(err, "Failed to save asset on ledger")
		}
	}

	return nil
}
