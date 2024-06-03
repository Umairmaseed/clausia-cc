package main

import (
	"github.com/hyperledger-labs/goprocess-cc/chaincode/assettypes"
	"github.com/hyperledger-labs/cc-tools/assets"
)

var assetTypeList = []assets.AssetType{
	assettypes.Secret,
	assettypes.Signer,
	assettypes.Document,
	assettypes.Signature,
}
