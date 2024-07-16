package contractassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Credit = assets.AssetType{
	Tag:         "credit",
	Label:       "Credit",
	Description: "Credit",

	Props: []assets.AssetProp{
		{
			Required: true,
			Tag:      "description",
			Label:    "Description",
			DataType: "string",
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "contract",
			Label:    "Contract",
			DataType: "->autoExecutableContract",
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "clause",
			Label:    "Clause",
			DataType: "->clause",
		},
		{
			Required: true,
			Tag:      "value",
			Label:    "Value",
			DataType: "number",
		},
	},
}
