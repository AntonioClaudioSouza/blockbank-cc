package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var InvoicePayment = assets.AssetType{
	Tag:         "invoicePayment",
	Label:       "Credit Card Invoice Payment",
	Description: "Credit Card Invoice Payment",

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
			Tag:      "creditCard",
			Label:    "Credit Card",
			DataType: "->creditCard",
			IsKey:    true,
		},
		{
			Required: true,
			Tag:      "owner",
			Label:    "Credit Card Owner",
			DataType: "->holder",
			IsKey:    true,
		},
		{
			Tag:          "value",
			Label:        "Payment value",
			DataType:     "number",
			DefaultValue: 0,
		},
		{
			Tag:      "date",
			Label:    "Payment date",
			DataType: "datetime",
			Required: true,
		},
	},
}
