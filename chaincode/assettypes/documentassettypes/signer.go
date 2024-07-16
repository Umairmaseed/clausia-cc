package documentassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Signer = assets.AssetType{
	Tag:         "signer",
	Label:       "Signer",
	Description: "Signer",

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
