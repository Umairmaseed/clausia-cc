package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Document = assets.AssetType{
	Tag:         "document",
	Label:       "Document",
	Description: "Document",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "orignalHash",
			Label:    "orignalHash",
			DataType: "sha256",
		},
		{
			Tag:      "finalHash",
			Label:    "finalHash",
			DataType: "sha256",
		},
		{
			Tag:      "status",
			Label:    "status",
			DataType: "statusType",
			Required: true,
		},
		{
			Tag:      "requiredSignatures",
			Label:    "requiredSignatures",
			DataType: "[]->signer",
			Required: true,
		},
		{
			Tag:      "successfulSignatures",
			Label:    "successfulSignatures",
			DataType: "[]->signer",
		},
		{
			Tag:      "rejectedSignatures",
			Label:    "rejectedSignatures",
			DataType: "[]->signer",
		},
		{
			Tag:      "orignalDocURL",
			Label:    "orignalDocURL",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "finalDocURL",
			Label:    "finalDocURL",
			DataType: "string",
		},
		{
			Tag:      "name",
			Label:    "name",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "signature",
			Label:    "signature",
			DataType: "->signature",
		},
	},
}
