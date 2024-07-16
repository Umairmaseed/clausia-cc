package main

import (
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/assettypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/assettypes/contractassettypes"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/assettypes/documentassettypes"
)

var assetTypeList = []assets.AssetType{
	assettypes.Secret,
	documentassettypes.Signer,
	documentassettypes.Document,
	contractassettypes.AutoExecutableContract,
	contractassettypes.Clause,
	contractassettypes.Deduction,
	contractassettypes.Credit,
}
