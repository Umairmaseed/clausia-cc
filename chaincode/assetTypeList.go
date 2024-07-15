package main

import (
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/goprocess-cc/chaincode/assettypes"
)

var assetTypeList = []assets.AssetType{
	assettypes.Secret,
	assettypes.Signer,
	assettypes.Document,
}
