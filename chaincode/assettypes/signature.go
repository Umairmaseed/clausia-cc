package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Signature = assets.AssetType{
	Tag:         "signature",
	Label:       "Signature",
	Description: "Signature",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "firstHash",
			Label:    "firstHash",
			DataType: "sha256",
		},
		{
			Tag:      "hash",
			Label:    "hash",
			DataType: "sha256",
			Required: true,
		},
		{
			Tag:      "prevHash",
			Label:    "prevHash",
			DataType: "sha256",
		},
		{
			Tag:      "pubKey",
			Label:    "pubKey",
			DataType: "pemPubKey",
		},
		{
			Tag:      "signer",
			Label:    "signer",
			DataType: "->signer",
			Required: true,
		},
		{
			Tag:      "url",
			Label:    "url",
			DataType: "string",
			Required: true,
		},
	},
}
