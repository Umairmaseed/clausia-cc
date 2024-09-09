package contractassettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var TemplateClause = assets.AssetType{
	Tag:         "templateClause",
	Label:       "Template Clause",
	Description: "Template Clause",

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
			Tag:      "template",
			Label:    "Template",
			DataType: "->template",
		},
		{
			Required:    true,
			Tag:         "number",
			Label:       "Number",
			DataType:    "number",
			Description: "Order of the clause in the template",
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
			Tag:      "category",
			Label:    "Category",
			DataType: "string",
		},
		{
			Tag:      "dependencies",
			Label:    "Dependencies",
			DataType: "[]->templateClause",
		},
		{
			Required: true,
			Tag:      "actionType",
			Label:    "Action Type",
			DataType: "actionType",
		},
		{
			Tag:         "defaultInputs",
			Label:       "Default Inputs",
			DataType:    "@object",
			Description: "This object is meant to hold default values for inputs that may be used on the creation of a contract with the template",
		},
		{
			Tag:         "defaultParameters",
			Label:       "Default Parameters",
			DataType:    "@object",
			Description: "This object is meant to hold default values for parameters that may be used on the creation of a contract with the template",
		},
		{
			Tag:      "optional",
			Label:    "Optional",
			DataType: "boolean",
		},
	},
}
