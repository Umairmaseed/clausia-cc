package contractassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var AutoExecutableContract = assets.AssetType{
	Tag:         "autoExecutableContract",
	Label:       "Auto Executable Contract",
	Description: "Auto Executable Contract",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "name",
			Label:    "Name",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "signatureDate",
			Label:    "Signature Date",
			DataType: "datetime",
		},
		{
			Tag:      "clauses",
			Label:    "Clauses",
			DataType: "[]->clause",
		},
		{
			Tag:      "data",
			Label:    "Data",
			DataType: "@object",
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "owner",
			Label:    "Owner",
			DataType: "->user",
		},
		{
			Tag:      "participants",
			Label:    "Participants",
			DataType: "[]->user",
		},
		{
			Tag:         "dates",
			Label:       "Dates",
			DataType:    "@object",
			Description: "This includes all the dates that are linked to the contract.",
		},
	},
}
