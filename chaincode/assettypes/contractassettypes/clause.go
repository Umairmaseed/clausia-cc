package contractassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Clause = assets.AssetType{
	Tag:         "clause",
	Label:       "Clause",
	Description: "Clause",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "Id",
			DataType: "string",
		},
		{
			Tag:      "description",
			Label:    "Description",
			DataType: "string",
		},
		{
			Tag:      "category",
			Label:    "Category",
			DataType: "string",
		},
		{
			Tag:      "parameters",
			Label:    "Parameters",
			DataType: "@object",
		},
		{
			Tag:      "input",
			Label:    "Input",
			DataType: "@object",
		},
		{
			Required: true,
			Tag:      "executable",
			Label:    "Executable",
			DataType: "boolean",
		},
		{
			Tag:      "dependencies",
			Label:    "Dependencies",
			DataType: "[]->clauses",
		},
		{
			Required: true,
			Tag:      "actionType",
			Label:    "Action Type",
			DataType: "actionType",
		},
		{
			Tag:      "finalized",
			Label:    "Finalized",
			DataType: "boolean",
		},
		{
			Tag:      "result",
			Label:    "Result",
			DataType: "@object",
		},
	},
}
