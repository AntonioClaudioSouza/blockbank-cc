package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var Purchase = assets.AssetType{
	Tag:         "purchase",
	Label:       "Purchase",
	Description: "Purchase",

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
			Tag:      "buyer",
			Label:    "buyer",
			DataType: "->holder",
			Writers:  []string{"orgMSP"},
		},
		{
			Required: true,
			Tag:      "value",
			Label:    "Purchase value",
			DataType: "number",
		},
		{
			Tag:      "date",
			Label:    "Purchase date",
			DataType: "datetime",
			Required: true,
		},
	},
}
