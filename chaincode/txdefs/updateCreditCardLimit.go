package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateCreditCardLimit = tx.Transaction{
	Tag:         "updateCreditCardLimit",
	Label:       "Update credit card limit",
	Description: "Update a credit card limit value",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "creditCard",
			Label:       "Credit card",
			Description: "Credit card to be updated",
			DataType:    "->creditCard",
			Required:    true,
		},
		{
			Tag:         "value",
			Label:       "Updated limit value",
			Description: "Updated limit value",
			DataType:    "number",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		newLimitValue, _ := req["value"].(float64)
		creditCardKey, _ := req["creditCard"].(assets.Key)

		creditCardMap, err := creditCardKey.Update(stub, map[string]interface{}{
			"limit": newLimitValue,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Error saving new limit asset on blockchain")
		}

		creditCardJSON, nerr := json.Marshal(creditCardMap)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return creditCardJSON, nil
	},
}
