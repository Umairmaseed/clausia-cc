package documentassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Document = assets.AssetType{
	Tag:         "document",
	Label:       "Document",
	Description: "Document",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "originalHash",
			Label:    "originalHash",
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
			DataType: "[]->user",
			Required: true,
		},
		{
			Tag:      "successfulSignatures",
			Label:    "successfulSignatures",
			DataType: "[]->user",
		},
		{
			Tag:      "rejectedSignatures",
			Label:    "rejectedSignatures",
			DataType: "[]->user",
		},
		{
			Tag:      "originalDocURL",
			Label:    "originalDocURL",
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
			Tag:      "owner",
			Label:    "Owner",
			DataType: "->user",
			Required: true,
		},
		{
			Tag:      "timeout",
			Label:    "Timeout",
			DataType: "datetime",
		},
	},
}
