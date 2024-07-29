package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var User = assets.AssetType{
	Tag:         "user",
	Label:       "User",
	Description: "User",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "cpf",
			Label:    "cpf",
			DataType: "cpf",
		},
		{
			Tag:      "email",
			Label:    "email",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "name",
			Label:    "name",
			DataType: "string",
			Required: true,
		},
		{
			Tag:      "phone",
			Label:    "phone",
			DataType: "string",
			Required: true,
		},
	},
}
