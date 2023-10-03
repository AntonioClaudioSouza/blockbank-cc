package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var Manager = assets.AssetType{
	Tag:         "manager",
	Label:       "Manager",
	Description: "Manager",

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
	},
}
