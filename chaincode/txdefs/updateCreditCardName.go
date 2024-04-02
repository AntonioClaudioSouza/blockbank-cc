package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateCreditCardName = tx.Transaction{
	Tag:         "updateCreditCardName",
	Label:       "Update credit card name",
	Description: "Update a credit card name",
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
			Tag:         "name",
			Label:       "Updated name",
			Description: "Updated name",
			DataType:    "string",
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		newName, _ := req["name"].(string)
		creditCardKey, _ := req["creditCard"].(assets.Key)

		creditCardMap, err := creditCardKey.Update(stub, map[string]interface{}{
			"creditCardName": newName,
		})
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update sender asset")
		}

		creditCardJSON, nerr := json.Marshal(creditCardMap)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return creditCardJSON, nil
	},
}
