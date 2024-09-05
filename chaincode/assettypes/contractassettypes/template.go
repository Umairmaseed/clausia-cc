package contractassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Template = assets.AssetType{
	Tag:         "template",
	Label:       "Template",
	Description: "Template",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "Id",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "name",
			Label:    "Name",
			DataType: "string",
		},
		{
			Tag:      "description",
			Label:    "Description",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "creator",
			Label:    "Creator",
			DataType: "->user",
		},
		{
			Required: true,
			Tag:      "public",
			Label:    "Public",
			DataType: "boolean",
		},
		{
			Tag:      "clauses",
			Label:    "Clauses",
			DataType: "[]->templateClause",
		},
	},
}
