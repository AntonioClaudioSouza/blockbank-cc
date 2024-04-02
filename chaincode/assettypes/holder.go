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
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "name",
			Label:    "Name",
			DataType: "string",
		},
		{
			Tag:          "cash",
			Label:        "Cash",
			DataType:     "number",
			DefaultValue: 0,
		},
		{
			Tag:          "ccAvailable",
			Label:        "Credit Card Available",
			DataType:     "boolean",
			DefaultValue: false,
		},
	},
}
