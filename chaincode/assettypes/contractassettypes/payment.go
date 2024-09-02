package contractassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Payment = assets.AssetType{
	Tag:         "payment",
	Label:       "Payment",
	Description: "Payment",

	Props: []assets.AssetProp{
		{
			Required: true,
			Tag:      "name",
			Label:    "name",
			DataType: "string",
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "hash",
			Label:    "hash",
			DataType: "sha256",
		},
		{
			Required: true,
			Tag:      "receiptUrl",
			Label:    "receiptUrl",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "payment",
			Label:    "payment",
			DataType: "number",
		},
		{
			Required: true,
			Tag:      "autoExecutableContract",
			Label:    "autoExecutableContract",
			DataType: "->autoExecutableContract",
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "clause",
			Label:    "clause",
			DataType: "->clause",
		},
	},
}
