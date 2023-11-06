package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var Withdrawal = assets.AssetType{
	Tag:         "withdrawal",
	Label:       "Withdrawal",
	Description: "Withdrawal",

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
			Label:        "Withdrawal value",
			DataType:     "number",
			DefaultValue: 0,
		},
		{
			Tag:      "date",
			Label:    "Withdrawal date",
			DataType: "datetime",
			Required: true,
		},
	},
}
