package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools-demo/chaincode/utils"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"

	// "github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateNewCreditCardPurchase = tx.Transaction{
	Tag:         "createNewCreditCardPurchase",
	Label:       "Create new credit card purchase",
	Description: "Create a new credit card purchase",
	Method:      "POST",
	Callers:     []string{"$orgMSP"},

	Args: []tx.Argument{

		{
			Tag:         "creditCard",
			Label:       "Credit Card",
			Description: "Credit Card",
			DataType:    "->creditCard",
			Required:    true,
		},
		{
			Tag:         "description",
			Label:       "Purchase description",
			Description: "Purchase description",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "value",
			Label:       "Purchase value",
			Description: "Purchase value",
			DataType:    "number",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		purchaseDescription := req["description"].(string)
		purchaseValue := req["value"].(float64)
		timeStamp, err := stub.Stub.GetTxTimestamp()

		creditCardKey, ok := req["creditCard"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter creditCard must be an asset")
		}

		creditCardAsset, err := creditCardKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to get owner from the ledger")
		}
		creditCardMap := (map[string]interface{})(*creditCardAsset)

		availableLimit := creditCardMap["limit"].(float64) - creditCardMap["limitUsed"].(float64)

		if availableLimit < purchaseValue {
			return nil, errors.WrapError(err, "Credit card doesn't have enough limit to make this purchase")
		}

		//REMOVE BALANCE FROM BUYER
		creditCardMap["limitUsed"] = creditCardMap["limitUsed"].(float64) + purchaseValue

		ccPurchaseMap := make(map[string]interface{})
		ccPurchaseMap["@assetType"] = "creditCardPurchase"
		ccPurchaseMap["txId"] = stub.Stub.GetTxID()
		ccPurchaseMap["description"] = purchaseDescription
		ccPurchaseMap["creditCard"] = creditCardMap
		ccPurchaseMap["value"] = purchaseValue
		ccPurchaseMap["date"] = utils.ReturnDate(timeStamp)

		creditCardMap, err = creditCardAsset.Update(stub, creditCardMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to update creditCard asset")
		}

		ccPurchaseAsset, err := assets.NewAsset(ccPurchaseMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create purchase asset")
		}

		_, err = ccPurchaseAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		ccPurchaseJSON, nerr := json.Marshal(ccPurchaseAsset)
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

		return ccPurchaseJSON, nil
	},
}
