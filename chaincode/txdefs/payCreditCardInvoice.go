package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var PayCreditCardInvoice = tx.Transaction{
	Tag:         "payCreditCardInvoice",
	Label:       "Pay credit card invoice",
	Description: "Pay credit card invoice",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{

		{
			Tag:         "creditCard",
			Label:       "Credit card",
			Description: "Credit card",
			DataType:    "->creditCard",
			Required:    true,
		},
		{
			Tag:         "valueToPay",
			Label:       "Value to pay",
			Description: "Value to pay",
			DataType:    "number",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		valueToPay, ok := req["valueToPay"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter valueToPay must be a number")
		}
		creditCardKey, ok := req["creditCard"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter creditCard must be an asset")
		}

		creditCardAsset, err := creditCardKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get creditCard from the ledger")
		}
		creditCardMap := (map[string]interface{})(*creditCardAsset)

		ownerProp, _ := creditCardAsset.GetProp("owner").(map[string]interface{})
		ownerKey := (assets.Key)(ownerProp)

		ownerMap, err := ownerKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}

		if valueToPay == 0 {
			valueToPay = creditCardMap["limitUsed"].(float64)
		}
		if valueToPay > creditCardMap["limitUsed"].(float64) {
			return nil, errors.WrapError(err, "Value to pay can't be higher than the used limit")
		}
		if valueToPay > ownerMap["cash"].(float64) {
			return nil, errors.WrapError(err, "Owner doesn't have enough balance")
		}

		//UPDATING DATA
		ownerMap["cash"] = ownerMap["cash"].(float64) - valueToPay
		creditCardMap["limitUsed"] = creditCardMap["limitUsed"].(float64) - valueToPay
		creditCardMap["owner"] = ownerMap

		creditCardMap, err = creditCardAsset.Update(stub, creditCardMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update creditCard asset")
		}

		ownerMap, err = ownerKey.Update(stub, ownerMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update holder asset")
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
