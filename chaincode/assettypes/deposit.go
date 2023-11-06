package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var Deposit = assets.AssetType{
	Tag:         "deposit",
	Label:       "Deposit",
	Description: "Deposit",

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
			Tag:      "holder",
			Label:    "Holder",
			DataType: "->holder",
			IsKey:    true,
		},
		{
			Tag:          "value",
			Label:        "Deposit value",
			DataType:     "number",
			DefaultValue: 0,
		},
		{
			Tag:      "date",
			Label:    "Deposit date",
			DataType: "datetime",
			Required: true,
		},
	},
}
