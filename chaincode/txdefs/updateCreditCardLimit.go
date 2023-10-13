package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateCreditCardLimit = tx.Transaction{
	Tag:         "updateCreditCardLimit",
	Label:       "Update credit card limit",
	Description: "Update a credit card limit value",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

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
		newLimitValue := req["value"].(float64)
		creditCardKey, ok := req["creditCard"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter creditCard must be an asset")
		}

		creditCardAsset, err := creditCardKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		creditCardMap := (map[string]interface{})(*creditCardAsset)

		creditCardMap["limit"] = newLimitValue

		creditCardMap, err = creditCardAsset.Update(stub, creditCardMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update sender asset")
		}

		creditCardJSON, nerr := json.Marshal(creditCardAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		// // Marshall message to be logged
		// logMsg, ok := json.Marshal(fmt.Sprintf("New library name: %s", name))
		// if ok != nil {
		// 	return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		// }

		// // Call event to log the message
		// events.CallEvent(stub, "createLibraryLog", logMsg)

		return creditCardJSON, nil
	},
}
