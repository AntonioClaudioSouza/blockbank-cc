package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var CreditCardPurchase = assets.AssetType{
	Tag:         "creditCardPurchase",
	Label:       "Credit Card Purchase",
	Description: "Credit Card Purchase",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Label:    "tdId",
			DataType: "string",
			Tag:      "txId",
		},
		{
			Required: true,
			IsKey:    true,
			Tag:      "description",
			Label:    "Purchase description",
			DataType: "string",
			Writers:  []string{"orgMSP"},
		},

		{
			Required: true,
			IsKey:    true,
			Tag:      "creditCard",
			Label:    "Credit Card",
			DataType: "->creditCard",
			Writers:  []string{"orgMSP"},
		},
		{
			Required: true,
			Tag:      "value",
			Label:    "Purchase value",
			DataType: "number",
		},
		{
			Required: true,
			Tag:      "date",
			Label:    "Purchase date",
			DataType: "datetime",
		},
	},
}
