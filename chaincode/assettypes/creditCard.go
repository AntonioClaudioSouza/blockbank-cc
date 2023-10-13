package assettypes

import (
	"fmt"
	"strconv"

	"github.com/hyperledger-labs/cc-tools/assets"
)

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
			Validate: func(number interface{}) error {
				numberSrt := number.(string)
				_, err := strconv.Atoi(numberSrt)
				if err != nil {
					return fmt.Errorf("Credit card number must contain only numbers")
				}
				if len(numberSrt) != 16 {
					return fmt.Errorf("Credit card number must have 16 digits")
				}
				return nil
			},
		},
		{
			Tag:      "creditCardName",
			Label:    "Credit Card Name",
			DataType: "string",
		},
		{
			Tag:      "owner",
			Label:    "Credit Card Owner",
			DataType: "->holder",
		},
		{
			Tag:          "limit",
			Label:        "Credit Card Limit",
			DataType:     "number",
			DefaultValue: 100,
		},
		{
			Tag:          "limitUsed",
			Label:        "Credit Card Used Limit",
			DataType:     "number",
			DefaultValue: 0,
		},
	},
}
