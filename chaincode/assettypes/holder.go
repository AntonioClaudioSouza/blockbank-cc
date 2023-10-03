package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Holder = assets.AssetType{
	Tag:         "holder",
	Label:       "Holder",
	Description: "Holder",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "document",
			Label:    "ID Document",
			DataType: "string",
			Writers:  []string{"orgMSP"},
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "name",
			Label:    "Name",
			DataType: "string",
			Writers:  []string{"orgMSP"},
		},
		{
			Tag:      "cash",
			Label:    "Cash",
			DataType: "number",
		},
		{
			Tag:      "ccAvailable",
			Label:    "Credit Card Available",
			DataType: "boolean",
		},
		{
			Tag:      "creditCard",
			Label:    "Credit Card",
			DataType: "->creditCard",
		},
	},
}
