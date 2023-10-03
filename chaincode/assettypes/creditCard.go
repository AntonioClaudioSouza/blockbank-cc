package assettypes

import "github.com/hyperledger-labs/cc-tools/assets"

var CreditCard = assets.AssetType{
	Tag:         "creditCard",
	Label:       "Credit Card",
	Description: "Credit Card",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "number",
			Label:    "Credit Card Number",
			DataType: "string",
			Writers:  []string{"orgMSP"},
		},
		{
			Tag:      "owner",
			Label:    "Credit Card Owner",
			DataType: "->holder",
		},
		{
			Tag:      "limit",
			Label:    "Credit Card Limit",
			DataType: "number",
		},
		{
			Tag:      "limitUsed",
			Label:    "Credit Card Used Limit",
			DataType: "number",
		},
	},
}
